package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
)

type Genre struct {
	GenreID   int    `json:"genre_id" bson:"genre_id"`
	GenreName string `json:"genre_name" bson:"genre_name"`
}

type Ranking struct {
	RankingValue int    `json:"ranking_value" bson:"ranking_value"`
	RankingName  string `json:"ranking_name" bson:"ranking_name"`
}

type Movie struct {
	ImdbID      string   `json:"imdb_id" bson:"imdb_id"`
	Title       string   `json:"title" bson:"title"`
	PosterPath  string   `json:"poster_path" bson:"poster_path"`
	YoutubeID   string   `json:"youtube_id" bson:"youtube_id"`
	Genre       []Genre  `json:"genre" bson:"genre"`
	AdminReview string   `json:"admin_review" bson:"admin_review"`
	Ranking     Ranking  `json:"ranking" bson:"ranking"`
}

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Get MongoDB URI
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/"
	}

	// Get database name
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "magic-stream-movies"
	}

	fmt.Println("Connecting to MongoDB...")
	fmt.Println("URI:", mongoURI)
	fmt.Println("Database:", dbName)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	fmt.Println("✅ Connected to MongoDB successfully!")

	// Get movies collection
	collection := client.Database(dbName).Collection("movies")

	// Read sample movies JSON
	fmt.Println("\nReading sample_movies.json...")
	data, err := ioutil.ReadFile("sample_movies.json")
	if err != nil {
		log.Fatal("Failed to read sample_movies.json:", err)
	}

	var movies []Movie
	err = json.Unmarshal(data, &movies)
	if err != nil {
		log.Fatal("Failed to parse JSON:", err)
	}

	fmt.Printf("Found %d movies in sample data\n", len(movies))

	// Clear existing movies (optional)
	fmt.Println("\nClearing existing movies...")
	deleteResult, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Println("Warning: Failed to clear collection:", err)
	} else {
		fmt.Printf("Deleted %d existing movies\n", deleteResult.DeletedCount)
	}

	// Insert sample movies
	fmt.Println("\nInserting sample movies...")
	var docs []interface{}
	for _, movie := range movies {
		docs = append(docs, movie)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal("Failed to insert movies:", err)
	}

	fmt.Printf("✅ Successfully inserted %d movies!\n", len(result.InsertedIDs))
	fmt.Println("\n🎬 Sample movies added to database!")
	fmt.Println("You can now start the backend server and browse movies in the frontend.")
}
