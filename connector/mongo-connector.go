package connector

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const uri = "mongodb://localhost:27017"

func getMongoDbConnection() *mongo.Client {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

func GetMongoDbCollection(DbName string, CollectionName string) *mongo.Collection {
	client := getMongoDbConnection()

	collection := client.Database(DbName).Collection(CollectionName)

	return collection
}
