package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/app/controllers"
	"github.com/mdozairq/e-voting-be/app/middlewares"
)

func VoterRoutes(superRoute *gin.RouterGroup) {
	votersRouter := superRoute.Group("/voters")
	{
		votersRouter.GET("/", controllers.GetVoters())
		votersRouter.POST("/adhaar", controllers.AddAdhaaarCard())
		votersRouter.POST("/signin", controllers.SignInVoter())
		votersRouter.POST("/verify", controllers.VerifyOTP())
		votersRouter.GET("/:id", middlewares.VoterTokenAuthMiddleware(), controllers.GetVoterByID())
		votersRouter.GET("/election", controllers.GetElectionByAadhaarLocation())
		votersRouter.POST("/election/vote", middlewares.VoterTokenAuthMiddleware(), controllers.CreateBallot())
	}
}
