package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID              bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          string        `json:"user_id" bson:"user_id"`
	FirstName       string        `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName        string        `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Email           string        `json:"email" bson:"email" validate:"required,email"`
	Password        string        `json:"password" bson:"password" validate:"required,min=6"`
	Role            string        `json:"role" bson:"role" validate:"oneof=ADMIN USER"`
	CreatedAt       time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" bson:"updated_at"`
	Token           string        `json:"token" bson:"token"`
	RefreshToken    string        `json:"refresh_token" bson:"refresh_token"`
	FavouriteGenres []Genre       `json:"favourite_genres" bson:"favourite_genres" validate:"required,dive"`
}
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	// we dont need the bson tag here because we are not storing the password in the database in plain text
}
type UserResponse struct { // this is the response that will be sent to the client side 
	// this is a DTO (Data Transfer Object) that is used to transfer data between the client and the server through network
	// without exposing the internal database structure to the client side 
	//in this case, we are not exposing the password to the client side 
	UserId          string  `json:"user_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Role            string  `json:"role"`
	Token           string  `json:"token"`
	RefreshToken    string  `json:"refresh_token"`
	FavouriteGenres []Genre `json:"favourite_genres"`
}
