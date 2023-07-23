package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/helpers"
	"github.com/mdozairq/e-voting-be/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	// router.Use(ResponseInterceptor())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(helpers.ResponseInterceptor())
	// router.Use(middlewares.Auth())
	allRoute := router.Group("/e-voting")
	routes.AddRoutes(allRoute)

	router.Run(":" + port)

}
