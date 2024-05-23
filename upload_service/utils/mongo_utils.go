package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var RepoCollection *mongo.Collection
var databaseName string

func init() {

	// Replace the placeholder with your connection string

	godotenv.Load(".env")
	var uri, x = os.LookupEnv("MONGO_URI")
	if !x {
		panic("MONGO_URI not found in .env")
	}
	log.Default().Println(uri, x)
	dbName, c := os.LookupEnv("MONGO_DB_NAME")
	if !c {
		panic("MONGO_DB_NAME not found in .env")
	}
	databaseName = dbName

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println(" You successfully connected to MongoDB!")
	mongoClient = client
	RepoCollection = GetCollection(databaseName, "projects")

}

// GetCollection returns a collection from the database
func GetCollection(dbName string, collectionName string) *mongo.Collection {
	collection := mongoClient.Database(dbName).Collection(collectionName)
	return collection
}

type TakaGoProject struct {
	ID      string `json:"id" bson:"id"`
	RepoURL string `json:"repo_url" bson:"repo_url"`
	Branch  string `json:"branch" bson:"branch" default:"main"`
	Status  string `json:"status" bson:"status" default:"initiated"`
}

type TakaGoProjectID struct {
	ID string `json:"id" bson:"id"`
}
type TakaGoProjectID_Status struct {
	ID     string `json:"id" bson:"id"`
	Status string `json:"status" bson:"status" default:"initiated"`
}

// InsertOne inserts a single document into the collection
func InsertOne(collection *mongo.Collection, document interface{}) {
	result, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		panic(err)
	}
	log.Default().Println(result)
}

func CheckIfExist(collection *mongo.Collection, filter interface{}) bool {
	result := FindOne(collection, filter)
	if result.Err() != nil {
		return false
	}
	return true

}

// FindOne finds a single document in the collection
func FindOne(collection *mongo.Collection, filter interface{}) *mongo.SingleResult {
	return collection.FindOne(context.Background(), filter)
}

// Find finds multiple documents in the collection
func Find(collection *mongo.Collection, filter interface{}) (*mongo.Cursor, error) {
	return collection.Find(context.Background(), filter)
}

// UpdateOne updates a single document in the collection
func UpdateOne(collection *mongo.Collection, filter interface{}, update TakaGoProject) (*mongo.UpdateResult, error) {
	return collection.UpdateOne(context.Background(), filter, bson.D{{"$set", update}})
	// return collection.UpdateOne(context.Background(), filter, update)
}

// DeleteOne deletes a single document from the collection
func DeleteOne(collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	return collection.DeleteOne(context.Background(), filter)
}

// DeleteMany deletes multiple documents from the collection
func DeleteMany(collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	return collection.DeleteMany(context.Background(), filter)
}
