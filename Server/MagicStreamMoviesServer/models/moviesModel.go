package models
// create movies struct

import(
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID int 
	GenreName string
}

type Ranking struct {
	RankingValue int //each sentiment excelent, good, ok, bad, terrible has an associtated value 1, 2, 3, 4, 5

	// excelent: 1, good: 2, ok: 3, bad: 4, terrible: 5

	RankingName string

}

type Movie struct {
	ID bson.ObjectID
	ImdbID string 	
	Title string
	PosterPath string 
	YoutubeID string 
	Genre []Genre    // array of struct genre 
	AdminReview string
	Ranking Ranking
}