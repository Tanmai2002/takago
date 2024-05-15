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

func init() {

	// Replace the placeholder with your connection string

	godotenv.Load(".env")
	var uri, x = os.LookupEnv("MONGO_URI")
	log.Default().Println(uri, x)

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println(" You successfully connected to MongoDB!")
	mongoClient = client

}

// GetCollection returns a collection from the database
func GetCollection(dbName string, collectionName string) *mongo.Collection {
	return mongoClient.Database(dbName).Collection(collectionName)
}

// InsertOne inserts a single document into the collection
func InsertOne(collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(context.Background(), document)
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
func UpdateOne(collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return collection.UpdateOne(context.Background(), filter, update)
}

// DeleteOne deletes a single document from the collection
func DeleteOne(collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	return collection.DeleteOne(context.Background(), filter)
}

// DeleteMany deletes multiple documents from the collection
func DeleteMany(collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	return collection.DeleteMany(context.Background(), filter)
}
