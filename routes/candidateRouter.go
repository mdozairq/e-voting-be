package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/controllers"
)

func CandidateRoutes(superRoute *gin.RouterGroup) {
	cnadidateRouters := superRoute.Group("/candidate")
	{
		cnadidateRouters.GET("/", controllers.GetCandidates())
		cnadidateRouters.POST("/signup", controllers.SignUpCandidate())
		cnadidateRouters.POST("/signin", controllers.SignInCandidate())
		cnadidateRouters.GET("/:id", controllers.GetCandidate())
		cnadidateRouters.PATCH("/:id", controllers.UpdateCandidate())
	}
}
