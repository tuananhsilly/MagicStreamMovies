package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10" //validate the movie data from the request body
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//return a collection of movies queried by the user to the client side
//want the func to be exportable



var validate = validator.New() //validate is a pointer to the validator package

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) { //how we are hooking to the gin framework
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		//cursor is a pointer to the database collection, with the purpose of iterating through the database collection

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
			return
		}
		defer cursor.Close(ctx) //for memory management, close the cursor after the function is done

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies"})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie(client *mongo.Client) gin.HandlerFunc { //easy to map a HTTP endpoint route to the relative HTTP endpoint handler function
	// also easy to create HTTP responses from within the relevant handler function
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second) //ctx is the context object carry the timeoutlled context.Background() is the parent context, 100*time.Second is the timeout duration
		defer cancel()                                                            //defer delay the execution of the function until the context is cancelled or the timeout is reached

		//use c to read a parameter from the HTTP request
		// map the parameter to the movie struct

		movieID := c.Param("imdb_id") // := is used to declare and assign a variable in one go

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movie models.Movie

        var movieCollection *mongo.Collection = database.OpenCollection("movies", client)


		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch movie"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movie models.Movie // store the movie data from the request body passing from client

        var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		// bind the movie data from the request body to the movie struct
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
			return // stop the function from executing further
		}

		//insert the movie data into the database
		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert movie"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Movie added successfully", "data": result})
	}
}

func AdminReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user have to be an admin to update the admin review
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Role not found in context"})
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update the admin review"})
			return
		}

		movieId := c.Param("imdb_id")
		//movieId will be passed in within the URL of the HTTP request to the handler function

		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}

		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		//bind the request body to the request struct
		//ShouldBindJSON is a method that binds the request body to the request struct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the admin review from the request body"})
			return
		}

		//use AI to extract the sentiment from the admin review that bounded to the request struct
		//OpenAI API using langchaingo

		//sort out the prompting instructions for the AI to extract the sentiment from the admin review
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview, client, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract sentiment from admin review"})
			return
		}

		filter := bson.M{"imdb_id": movieId}

		update := bson.M{"$set": bson.M{"admin_review": req.AdminReview, "ranking": bson.M{"ranking_name": sentiment, "ranking_value": rankVal}}}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie ranking"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)

	}
}

// function to extract the sentiment from the admin review using the OpenAI API
func GetReviewRanking(admin_review string, client *mongo.Client, c *gin.Context) (string, int, error) {
	//querry the ranking name and value from the database based on the admin review
	//add the admin_review to the prompt
	rankings, err := GetRankings(client, c)
	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited += ranking.RankingName + ","
		}
	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	//load the base prompt into memory
	err = godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: Failed to load .env file")
	}

	//now need to read the environment variable for the OpenAI API key

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")

	if OpenAiApiKey == "" {
		return "", 0, errors.New("OPENAI_API_KEY is not set")
	}

	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}

	//read the base prompt template from the .env file
	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	//we've read in all the rankings name and the sentiment delimited string, now we need to create the base prompt for the AI to extract the sentiment from the admin review

	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(c, base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}

	rankVal := 0
	//check rankings sent back by the AI is indeed in the ranking collection
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil

}

func GetRankings(client *mongo.Client, c *gin.Context) ([]models.Ranking, error) {
	var rankings []models.Ranking

	var rankingCollection *mongo.Collection = database.OpenCollection("rankings", client)

	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	//pull all the rankings from the database into the rankings slice
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

func GetReccomendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		//user have to login to get the recommended movies
		userId, err := utils.GetUserIdFromContext(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}

		favourite_genres, err := GetUsersFavouriteGenres(userId, client, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env file not found")
		}
		var recommendedMovieLimitVal int64 = 5

		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")

		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}

		findOptions := options.Find()

		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}}) //sort the movies by the ranking value in ascending order
		// need to set the limit to the number of movies to be recommended
		findOptions.SetLimit(recommendedMovieLimitVal)

		filter := bson.D{
			{Key: "genre.genre_name", Value: bson.D{
				{Key: "$in", Value: favourite_genres},
			}},
		}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie

		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUsersFavouriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {

	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: userId}}

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)
	var result bson.M

	var userCollection *mongo.Collection = database.OpenCollection("users", client)
	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}

	favGenresArray, ok := result["favourite_genres"].(bson.A)

	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}

	var genreNames []string

	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil

}

func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var genreCollection *mongo.Collection = database.OpenCollection("genres", client)

		cursor, err := genreCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie genres"})
			return
		}
		defer cursor.Close(ctx)

		var genres []models.Genre
		if err := cursor.All(ctx, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, genres)

	}
}
