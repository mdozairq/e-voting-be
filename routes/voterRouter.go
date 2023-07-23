package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/controllers"
)

func VoterRoutes(superRoute *gin.RouterGroup) {
	votersRouter := superRoute.Group("/voters")
	{
		votersRouter.GET("/", controllers.GetVoters())
		votersRouter.POST("/adhaar", controllers.AddAdhaaarCard() )
		votersRouter.POST("/signup", controllers.SignUpVoter())
		votersRouter.POST("/signin", controllers.SignInVoter())
		votersRouter.GET("/:id", controllers.GetVoter())
	}
}
