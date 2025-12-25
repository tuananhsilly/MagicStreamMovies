package controllers

import(
	"github.com/gin-gonic/gin"
	"context"
	"time"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

//return a collection of movies queried by the user to the client side 
//want the func to be exportable 

var movieCollection *mongo.Collection = database.OpenCollection("movies")

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context){     //how we are hooking to the gin framework
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
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second) //ctx is the context object carry the timeoutlled context.Background() is the parent context, 100*time.Second is the timeout duration
		defer cancel() //defer delay the execution of the function until the context is cancelled or the timeout is reached


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

