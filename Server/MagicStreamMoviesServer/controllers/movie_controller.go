package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10" //validate the movie data from the request body
	"github.com/joho/godotenv"
	"github.com/langchain-ai/langchaingo/langchain/chat_models/openai"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

//return a collection of movies queried by the user to the client side
//want the func to be exportable

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")
var validate = validator.New() //validate is a pointer to the validator package

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) { //how we are hooking to the gin framework
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

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

func GetMovie() gin.HandlerFunc { //easy to map a HTTP endpoint route to the relative HTTP endpoint handler function
	// also easy to create HTTP responses from within the relevant handler function
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second) //ctx is the context object carry the timeoutlled context.Background() is the parent context, 100*time.Second is the timeout duration
		defer cancel()                                                            //defer delay the execution of the function until the context is cancelled or the timeout is reached

		//use c to read a parameter from the HTTP request
		// map the parameter to the movie struct

		movieID := c.Param("imdb_id") // := is used to declare and assign a variable in one go

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch movie"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie // store the movie data from the request body passing from client

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

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract sentiment from admin review"})
			return
		}

		filter := bson.M{"imdb_id": movieId}

		update := bson.M{"$set": bson.M{"admin_review": req.AdminReview, "ranking": bson.M{"ranking_name": sentiment, "ranking_value": rankVal}}}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

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

		c.JSON(http.StatusOK, gin.H{"message": "Movie ranking updated successfully"})

	}
}

// function to extract the sentiment from the admin review using the OpenAI API
func GetReviewRanking(admin_review string) (string, int, error) {
	//querry the ranking name and value from the database based on the admin review
	//add the admin_review to the prompt
	rankings, err := GetRankings()
	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited += sentimentDelimited + ranking.RankingName + ","
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

	response, err := llm.Call(context.Background(), base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}

	rankVal := 0
	//check rankings sent back by the AI is indeed in the ranking collection
	for _, ranking := range rankings {
		if ranking.RankingName == response.Content {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response.Content, rankVal, nil

}

func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
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
