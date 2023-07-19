package routes

import (
    "github.com/gin-gonic/gin"
	controllers "github.com/mdozairq/e-voting-be/controllers"
)

func AdminRoutes(incominRoutes *gin.Engine){
	incominRoutes.POST("/admin/signin", controllers.SignInAdmin())
}
