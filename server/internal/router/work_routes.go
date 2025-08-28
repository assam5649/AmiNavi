package router

import (
	"github.com/gin-gonic/gin"
	"server/internal/handler"
	"server/internal/services"
)

func WorkRouter(router *gin.RouterGroup, service *services.WorkServices) {
	workHandler := handler.NewWorkHandler(service)

	router.GET("/works", workHandler.GetAll)
	router.POST("/works", workHandler.Create)

	router.GET("/works/:id", workHandler.GetByID)
	router.PUT("/works/:id", workHandler.PutByID)
	router.PATCH("/works/:id", workHandler.PatchByID)
	router.DELETE("/works/:id", workHandler.DeleteByID)
}
