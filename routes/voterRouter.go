package routes

import (
    "github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/controllers"
)

func VoterRoutes(incominRoutes *gin.Engine){
	incominRoutes.GET("/voter", controllers.GetVoters())
	incominRoutes.POST("/voter/signup", controllers.SignUpVoter())
	incominRoutes.POST("/voter/signin", controllers.SignInVoter())
	incominRoutes.GET("/voter/:id", controllers.GetVoter())
}
