package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	routes "github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMoviesServer!")
	})

	//Facilitation of separating our HTTP endpoints into protected(requires authentication) and unprotected(does not require authentication) endpoints

	routes.SetupProtectedRoutes(router)
	routes.SetupUnprotectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
