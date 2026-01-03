package database

//define the API connection to the database

import(
	"fmt"
	"log"
	"os" // use for environment variables configured in the .env file
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"github.com/joho/godotenv"
)

//connect the Gin framework to the database

func Connect() *mongo.Client {
	//load the environment variables from the .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file") //now just log the error and continue
	}

	MongoDB := os.Getenv("MONGODB_URI")

	if MongoDB == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	fmt.Println("MongoDB URI: ", MongoDB)

	//connect client options to the database

	clientOptions := options.Client().ApplyURI(MongoDB)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Println("Error connecting to MongoDB: ", err)
	}

	return client
}


//want to create a function that will be used to open the actual connection to the database

func OpenCollection (collectionName string, client *mongo.Client) *mongo.Collection {
	err := godotenv.Load(".env")

	if err != nil{
		log.Println("Error loading .env file")
	}

	//read the database name from the .env file
	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("Database name: ", databaseName)

	collection := client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		log.Println("Collection not found")
	}

	return collection
}