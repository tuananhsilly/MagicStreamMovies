package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	// Create a route group and apply middleware ONLY to this group
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleWare()) // This applies middleware only to routes in this group

	//HTTP endpoints that are protected, cant be reached without a valid access token
	protected.GET("/movie/:imdb_id", controller.GetMovie(client))
	protected.POST("/addmovie", controller.AddMovie(client))
	protected.GET("/recommended_movies", controller.GetReccomendedMovies(client))
	protected.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))


}
