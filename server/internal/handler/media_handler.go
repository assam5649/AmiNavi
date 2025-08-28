package handler

import "C"
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/dto/post"
	"server/internal/services"
)

type MediaHandler struct {
	MediaService *services.MediaServices
}

func NewMediaHandler(mediaService *services.MediaServices) *MediaHandler {
	return &MediaHandler{MediaService: mediaService}
}

func (h *MediaHandler) HandlerConversion(c *gin.Context) {
	var request post.MediaRequest
	var response post.MediaResponse

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	csv, gcsPath, err := h.MediaService.CsvService.ConvertAndUpload(
		c.Request.Context(),
		request.Image,
		h.MediaService.Request,
		h.MediaService.Upload,
	)
	if err != nil {
		return
	}

	response.Csv = csv
	response.CsvUrl = gcsPath

	c.JSON(http.StatusOK, response)

	return
}
