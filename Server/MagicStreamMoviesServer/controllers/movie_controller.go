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
		}
		defer cursor.Close(ctx) //for memory management, close the cursor after the function is done

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies"})
		}
		
		c.JSON(http.StatusOK, movies)
	}
}