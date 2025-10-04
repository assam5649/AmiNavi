package router

import (
	"server/internal/handler"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

func MediaRouter(router *gin.RouterGroup, service *services.MediaServices) {
	mediaHandler := handler.NewMediaHandler(service)

	router.POST("/csv-conversions", mediaHandler.HandlerConversion)
	router.POST("/ocr", mediaHandler.HandlerOCR)
}
