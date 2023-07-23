package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/controllers"
)

func CandidateRoutes(superRoute *gin.RouterGroup) {
	cnadidateRouters := superRoute.Group("/candidate")
	{
		cnadidateRouters.GET("/candidates", controllers.Getcandidates())
		cnadidateRouters.POST("/candidate/signup", controllers.SignUpCandidate())
		cnadidateRouters.POST("/candidate/signin", controllers.SignInCandidate())
		cnadidateRouters.GET("/candidate/:id", controllers.GetCandidate())
		cnadidateRouters.PATCH("/candidate/:id", controllers.UpdateCandidate())
	}
}
