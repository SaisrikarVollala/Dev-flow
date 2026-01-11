package routes

import (
	controller "github.com/SaisrikarVollala/Dev-flow/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoute(api *gin.RouterGroup) {
	auth:=api.Group("/auth");
	auth.POST("/login",controller.Login)
	auth.POST("/register",controller.Register)
}