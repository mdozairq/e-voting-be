package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/mdozairq/e-voting-be/config"
	"github.com/mdozairq/e-voting-be/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstances() *mongo.Client {
	config.LoadEnv()
	utils.LogInfo("env loaded")
	dbConfig := config.NewDBConfig()
	utils.LogInfo(config.NewDBConfig().MongoUri)
	MongoUri := dbConfig.MongoUri
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
	dbConfig := config.NewDBConfig()
	var collection *mongo.Collection = client.Database(dbConfig.Dbname).Collection(collectionName)
	return collection
}