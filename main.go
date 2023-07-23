package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/config"
	"github.com/mdozairq/e-voting-be/helpers"
	"github.com/mdozairq/e-voting-be/routes"
	"github.com/mdozairq/e-voting-be/utils"
	// "github.com/twilio/twilio-go"
)



func main() {
	config.LoadEnv()
	utils.LogInfo("env loaded")
	serverConfig := config.NewServerConfig()
	port := serverConfig.Port
	if port == "" {
		port = "8080"
	}

	// accountSid := config.NewTwilioConfig().Sid
	// authToken := config.NewTwilioConfig().AuthToken

	// Initialize Twilio client
	// var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: accountSid,
	// 	Password: authToken,
	// })
	
	router := gin.New()
	router.Use(helpers.ResponseInterceptor())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	// router.Use(func(c *gin.Context) {
	// 	c.Set("twilioClient", client)
	// 	c.Next()
	// })
	allRoute := router.Group(serverConfig.BasePath)
	routes.AddRoutes(allRoute)

	router.Run(":" + port)

}
