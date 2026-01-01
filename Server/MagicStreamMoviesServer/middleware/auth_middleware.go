package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"net/http"
)

func AuthMiddleWare() gin.HandlerFunc {
	//used to validate the access token and grant access or prohibit access to the protected resources

	return func(c *gin.Context) {
		//first need to extract the access token from the request header in token_util.go
		token, err := utils.GetAcccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		
		claims, err := utils.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		//retrieve the claims from the token
		//once the token is validated, means the user is authenticated, so we can set the claims to the context
		c.Set("role", claims.Role)
		c.Set("user_id", claims.UserId)

		c.Next() // continue to execute the targeted HTTP endpoint handler function
	}
}