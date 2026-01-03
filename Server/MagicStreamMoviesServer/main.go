package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	routes "github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"context"
	"log"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMoviesServer!")
	})

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//Facilitation of separating our HTTP endpoints into protected(requires authentication) and unprotected(does not require authentication) endpoints

	var client *mongo.Client = database.Connect()
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to connect to the server: %v", err)
	}

	defer func(){
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	} ()


	routes.SetupProtectedRoutes(router, client)
	routes.SetupUnprotectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
