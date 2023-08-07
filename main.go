package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/config"
	// "github.com/mdozairq/e-voting-be/app/helpers"
	"github.com/mdozairq/e-voting-be/app/routes"
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
	// router.Use(helpers.ResponseInterceptor())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // Maximum cache age for CORS preflight requests (12 hours)
	}))
	allRoute := router.Group(serverConfig.BasePath)
	routes.AddRoutes(allRoute)

	router.Run(":" + port)

}
