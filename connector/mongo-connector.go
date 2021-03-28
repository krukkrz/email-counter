package connector

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	uri            = "mongodb://localhost:27017"
	dbName         = "emailslistsdb"
	collectionName = "list"
)

func GetMongoDbCollection() *mongo.Collection {
	client := getMongoDbClient()

	collection := client.Database(dbName).Collection(collectionName)

	return collection
}

func getMongoDbClient() *mongo.Client {

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
