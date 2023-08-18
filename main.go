package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/config"

	// "github.com/mdozairq/e-voting-be/app/helpers"
	"github.com/mdozairq/e-voting-be/app/helpers"
	"github.com/mdozairq/e-voting-be/app/routes"
	"github.com/mdozairq/e-voting-be/eth"
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

	ethClient := eth.GetEthClient()
	Conn := helpers.GetConnection(ethClient, "038eef0c11980a31af4bc3871945a59db253545646191337026ccd6b893aa7ae")

	reply, err := Conn.Balance(&bind.CallOpts{})
	if err != nil {
		// handle error
		log.Println(err)
	}
	log.Println(reply.Int64())
	fmt.Println(Conn)

	router := gin.New()
	// router.Use(helpers.ResponseInterceptor())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // Maximum cache age for CORS preflight requests (12 hours)
	}))
	allRoute := router.Group(serverConfig.BasePath)
	routes.AddRoutes(allRoute)

	router.Run(":" + port)

}
