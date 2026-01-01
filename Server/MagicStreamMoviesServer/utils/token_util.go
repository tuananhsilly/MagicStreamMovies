package utils

import(
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	jwt "github.com/golang-jwt/jwt/v5"
	"os"
	"time"
	"context"
	"errors"
)

type SignedDetails struct { //this is the struct that is used to store the signed details of the JWT
	Email string 
	FirstName string 
	LastName string 	
	Role string
	UserId string
	jwt.RegisteredClaims //another struct that is used to store the registered claims of the JWT
	//contains the standard claims of the JWT like the issuer means who is issuing the token, subject means who is the subject of the token, audience means who is the audience of the token, etc.
	//prevents replay attacks and ensures the integrity of the token
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var REFRESH_SECRET_KEY string = os.Getenv("REFRESH_SECRET_KEY")
var userCollection *mongo.Collection = database.OpenCollection("users")

func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error){ 
	//return the access token and the refresh token, and the error if any
	claims := &SignedDetails{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Role: role,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "auth.magicstreammovies.com",
			Subject: userId,
			Audience: jwt.ClaimStrings{"magicstreammovies.com", "magicstreammovies.app"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //expires in 24 hours
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetails{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Role: role,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "auth.magicstreammovies.com",
			Subject: userId,
			Audience: jwt.ClaimStrings{"magicstreammovies.com", "magicstreammovies.app"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)), //expires in 30 days
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(REFRESH_SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func UpdateAllTokens(userId, token, refreshToken string)(err error) {

	//create the usual resource clearing code (housekeeping code) when the timeout is reached 

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"$set": bson.M{
			"token": token,
			"refresh_token": refreshToken,
			"updated_at": updateAt,
		},
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	//update the user document in the database with the new token and refresh token filtered by the user ID

	if err != nil {
		return err 
	}
	return nil
}


func GetAcccessToken (c *gin.Context) (string, error){
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}
	// Check if header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}
	// Extract the token string after "Bearer "
	tokenString := authHeader[len(bearerPrefix):]

	if tokenString == "" {
		return "", errors.New("Bearer token is required")
	}


	return tokenString, nil
} 

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{} // the struct that embeds or extends the RegisteredClaims struct
	//where the token claims will be decoded into 
	 

	// This use jwt.ParseWithClaims to parse the token and decode the claims into the SignedDetails struct
	//use a callback function to verify the token signature
	//return as the byte slice of the secret key
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}


	//check that the signed algorithm is the same as the algorithm used to sign the token
	//A critical security check to prevent token tampering
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	//check that the token has not expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil

}