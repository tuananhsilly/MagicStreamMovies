package controllers

import(
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"github.com/gin-gonic/gin"
	//wanna valiadate the user data from the request body to server 
	"github.com/go-playground/validator/v10"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
)



//basic structure for creating the user endpoint handler function

func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	 return string(hashedPassword), nil
}


func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context){
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {         // reference the user data from the request body to the user struct
			 c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			 return
		} 

		validate := validator.New()

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
			return
		}

		//before inserting the user data into the database, we dont want to save the password in the database in plain text
		//let's hash the password using bcrypt

		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		//object helps clear the resources allocated to the function and extract the data from the user struct

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		//each user have a validate email address
		//count how many users have the same email address

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing user"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}

		// if the user is new, generate a unique user ID for the user
		user.UserID = bson.NewObjectID().Hex() //generate a unique user ID for the user
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": result})
	}

}

func LoginUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context){
		var userLogin models.UserLogin
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)	
		defer cancel()

		var foundUser models.User
		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return 
		}

		//compare the password from the request body with the password from the database
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		//user is authenticated, generate an access token and a refresh token for the user
		//an access token is a short-lived token that is used to access the protected resources
		//a refresh token is a long-lived token that is used to refresh the access token
		//this is wonderful when we have an app with disparated loosely coupled components like a frontend and a backend, where the frontend needs to access the protected resources of the backend
		// client dont need to provide the credentials to the backend everytime they want to access the protected resources, they can just use the access token to access the protected resources


		//now generate the tokens for the user
		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		//update the user document in the database with the new token and refresh token
		err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		//once the user logged in, we want to return the user data to the client side
		//return the token and refresh token to the client side
		c.JSON(http.StatusOK, models.UserResponse{
			UserId: foundUser.UserID,
			FirstName: foundUser.FirstName,
			LastName: foundUser.LastName,
			Email: foundUser.Email,
			Role: foundUser.Role,
			FavouriteGenres: foundUser.FavouriteGenres,
			Token: token,
			RefreshToken: refreshToken,
		})
	}
}