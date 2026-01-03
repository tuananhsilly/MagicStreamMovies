package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	routes "github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
)

func main() {
	router := gin.Default()

	// Enable CORS for frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
