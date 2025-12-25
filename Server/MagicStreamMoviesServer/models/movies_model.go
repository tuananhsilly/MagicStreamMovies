package models
// create movies struct

import(
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID int 		`bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string 	`bson:"genre_name" json:"genre_name" validate:"required, min = 2, max = 100"`
}

type Ranking struct {
	RankingValue int `bson:"ranking_value" json:"ranking_value" validate:"required, min = 1, max = 5"`//each sentiment excelent, good, ok, bad, terrible has an associtated value 1, 2, 3, 4, 5

	// excelent: 1, good: 2, ok: 3, bad: 4, terrible: 5

	RankingName string `bson:"ranking_name" json:"ranking_name" validate:"required, min = 2, max = 100"`

}

type Movie struct {
	//BSON format refers to how the fields are stored in the database
	//JSON format refers to how the fields are stored in the API response
	ID bson.ObjectID 	`bson:"_id" json:"_id"`
	ImdbID string 		`bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title string 		`bson:"title" json:"title" validate:"required, min = 2, max = 500"`
	PosterPath string 	`bson:"poster_path" json:"poster_path" validate:"required, url"`
	YoutubeID string 	`bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre []Genre   	`bson:"genre" json:"genre" validate:"required, dive"` // dive is used to validate the nested array of structs
	AdminReview string 	`bson:"admin_review" json:"admin_review" validate:"required"`
	Ranking Ranking 	`bson:"ranking" json:"ranking" validate:"required"`
}

//establish rules for each of fields so that only valid data can be stored in the database