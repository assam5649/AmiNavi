package router

import (
	"server/internal/handler"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.RouterGroup, service *services.AuthServices) {
	authHandler := handler.NewAuthHandler(service)

	router.POST("/users", authHandler.Register)
	router.PATCH("/users", authHandler.Update)
}
