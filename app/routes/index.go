package routes

import "github.com/gin-gonic/gin"


func AddRoutes(superRoute *gin.RouterGroup) {
	VoterRoutes(superRoute)
	AdminRoutes(superRoute)
	CandidateRoutes(superRoute)
}