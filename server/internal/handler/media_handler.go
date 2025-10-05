package handler

import "C"
import (
	"io"
	"net/http"
	"server/internal/dto/post"
	"server/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	MediaService *services.MediaServices
}

func NewMediaHandler(mediaService *services.MediaServices) *MediaHandler {
	return &MediaHandler{MediaService: mediaService}
}

func (h *MediaHandler) HandlerFixCSV(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body."})
		return
	}

	csvData, fileName, err := h.MediaService.FixCsv.FixCsv(
		c.Request.Context(),
		data,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fix csv."})
		return
	}

	response := post.FixResponse{
		Csv:      csvData,
		FileName: fileName,
	}

	c.JSON(http.StatusAccepted, response)
}

func (h *MediaHandler) HandlerOCR(c *gin.Context) {
	var response post.OCRResponse

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	imageBytes, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	csv, err := h.MediaService.OCR.OCR(
		c.Request.Context(),
		imageBytes,
	)
	if err != nil {
		errorMessage := err.Error()
		statusCode := http.StatusInternalServerError

		if strings.Contains(errorMessage, "status 400") {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{"error": "OCR processing failed.", "details": errorMessage})
		return
	}

	response.Pattern = string(csv)

	c.JSON(http.StatusOK, response)
}

func (h *MediaHandler) HandlerConversion(c *gin.Context) {
	var response post.MediaResponse

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	imageBytes, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	csv, fileName, err := h.MediaService.CsvService.ConvertAndUpload(
		c.Request.Context(),
		imageBytes,
		h.MediaService.Request,
		h.MediaService.Upload,
	)
	if err != nil {
		errorMessage := err.Error()
		statusCode := http.StatusInternalServerError

		if strings.Contains(errorMessage, "status 400") {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{"error": "ML processing failed.", "details": errorMessage})
		return
	}

	response.Csv = csv
	response.FileName = fileName

	c.JSON(http.StatusOK, response)
}
