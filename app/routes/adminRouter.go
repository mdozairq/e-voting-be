package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/app/controllers"
	"github.com/mdozairq/e-voting-be/app/middlewares"
)

func AdminRoutes(superRoute *gin.RouterGroup) {
	adminRouters := superRoute.Group("/admin")
	{
		adminRouters.POST("/signin", controllers.SignInAdmin())
		adminRouters.POST("/election/initialize", middlewares.AdminTokenAuthMiddleware(), controllers.InitializeElection())
		adminRouters.GET("/election/all", middlewares.AdminTokenAuthMiddleware(), controllers.GetAllElections())
		adminRouters.GET("/election/:id", middlewares.AdminTokenAuthMiddleware(), controllers.GetElectionByID())
		adminRouters.GET("/constituency/all", middlewares.AdminTokenAuthMiddleware(), controllers.GetAllConstituencies())
		adminRouters.POST("/party/create", middlewares.AdminTokenAuthMiddleware(), controllers.CreateParty())
	}

}
