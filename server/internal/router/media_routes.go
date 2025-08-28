package router

import (
	"github.com/gin-gonic/gin"
	"server/internal/handler"
	"server/internal/services"
)

func MediaRouter(router *gin.RouterGroup, service *services.MediaServices) {
	mediaHandler := handler.NewMediaHandler(service)

	router.POST("/csv-conversions", mediaHandler.HandlerConversion)
}
