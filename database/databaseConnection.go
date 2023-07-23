package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstances() *mongo.Client {
	MongoUri := "mongodb://localhost:27017"
	fmt.Println(MongoUri)
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoUri))

	if err!= nil {
			log.Fatal(err)
		}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel();

	err = client.Connect(ctx)

	if err!= nil {
			log.Fatal(err)
		}
	
		fmt.Println("Connected to MongoDB")
		return client
}

var Client *mongo.Client = DBinstances()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("e-voting").Collection(collectionName)
	return collection
}