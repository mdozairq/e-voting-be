package routes

import (
    "github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/controllers"
)

func CandidateRoutes(incominRoutes *gin.Engine){
	incominRoutes.GET("/candidate", controllers.Getcandidates())
	incominRoutes.POST("/candidate/signup", controllers.SignUpCandidate())
	incominRoutes.POST("/candidate/signin", controllers.SignInCandidate())
	incominRoutes.GET("/candidate/:id", controllers.GetCandidate())
	incominRoutes.PATCH("/candidate/:id", controllers.UpdateCandidate())


}
