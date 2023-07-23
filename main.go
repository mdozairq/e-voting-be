package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/config"
	"github.com/mdozairq/e-voting-be/helpers"
	"github.com/mdozairq/e-voting-be/routes"
	"github.com/mdozairq/e-voting-be/utils"
)

func main() {
	config.LoadEnv()
	utils.LogInfo("env loaded")
	serverConfig := config.NewServerConfig()
	port := serverConfig.Port
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(helpers.ResponseInterceptor())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	allRoute := router.Group(serverConfig.BasePath)
	routes.AddRoutes(allRoute)

	router.Run(":" + port)

}
