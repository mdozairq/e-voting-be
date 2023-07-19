
package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/database"
	"github.com/mdozairq/e-voting-be/routes"
	"github.com/mdozairq/e-voting-be/middlewares"
	"go.mongodb.org/mongo-driver/mongo"
)

var voterCollection *mongo.Collection = database.OpenCollection(database.Client, "voter")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
			port = "8080"
		}

	router := gin.new()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.Auth())
	router.voterRoutes(router)

	router.Run(":"+port)
}
