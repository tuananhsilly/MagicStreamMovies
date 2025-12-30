package utils

import(
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/tuananhsilly/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	jwt "github.com/golang-jwt/jwt/v5"
	"os"
	"time"
	"context"

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
