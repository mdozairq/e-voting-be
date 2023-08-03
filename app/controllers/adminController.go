package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mdozairq/e-voting-be/config"
	"github.com/mdozairq/e-voting-be/app/helpers"
)

type AdminSignInDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}


func SignInAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var adminSignInDto AdminSignInDto
		if err := c.ShouldBindJSON(&adminSignInDto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := validator.New()
		if err := validate.Struct(adminSignInDto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		adminConfig := config.NewAdminConfig()
		if adminSignInDto.Email != adminConfig.AdminEmail || adminSignInDto.Password != adminConfig.AdminPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		adminToken, err := helpers.GenerateAdminToken(adminConfig.AdminEmail, adminConfig.AdminAuthToken, "admin")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate admin token"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, gin.H{"token": adminToken})
	}
}
