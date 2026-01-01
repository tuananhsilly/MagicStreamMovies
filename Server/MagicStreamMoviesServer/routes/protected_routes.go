package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
)

func SetupProtectedRoutes(router *gin.Engine) {
	// Create a route group and apply middleware ONLY to this group
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleWare()) // This applies middleware only to routes in this group

	//HTTP endpoints that are protected, cant be reached without a valid access token
	protected.GET("/movie/:imdb_id", controller.GetMovie())
	protected.POST("/addmovie", controller.AddMovie())
}
