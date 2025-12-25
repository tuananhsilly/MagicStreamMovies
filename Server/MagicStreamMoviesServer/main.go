package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"

)

func main(){
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context){
		c.String(200, "Hello, MagicStreamMoviesServer!")  
	})

	router.GET("/movies", controllers.GetMovies())
	
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}