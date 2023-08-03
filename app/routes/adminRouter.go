package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/app/controllers"
)

func AdminRoutes(superRoute *gin.RouterGroup) {
	adminRouters := superRoute.Group("/admin")
	{
		adminRouters.POST("/signin", controllers.SignInAdmin())
		adminRouters.POST("/election/initialize", controllers.InitializeElection())
	}

}