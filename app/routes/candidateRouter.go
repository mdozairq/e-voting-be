package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/app/controllers"
	"github.com/mdozairq/e-voting-be/app/middlewares"
)

func CandidateRoutes(superRoute *gin.RouterGroup) {
	candidateRouters := superRoute.Group("/candidate")
	{
		candidateRouters.GET("/", middlewares.CandidateTokenAuthMiddleware(), controllers.GetCandidates())
		candidateRouters.POST("/signup", controllers.SignUpCandidate())
		candidateRouters.POST("/signin", controllers.SignInCandidate())
		candidateRouters.GET("/:id", controllers.GetCandidate())
		candidateRouters.PATCH("/:id", controllers.UpdateCandidate())
		candidateRouters.GET("/election",  middlewares.CandidateTokenAuthMiddleware(), controllers.GetRegistrationElections())
		candidateRouters.GET("party/all", middlewares.CandidateTokenAuthMiddleware(), controllers.GetAllParties())
	}
}
