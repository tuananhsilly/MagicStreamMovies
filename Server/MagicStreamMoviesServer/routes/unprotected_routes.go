package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
)

func SetupUnprotectedRoutes(router *gin.Engine) {
	router.GET("/movies", controller.GetMovies())
	// router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
}
