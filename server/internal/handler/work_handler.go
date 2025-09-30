package handler

import "C"
import (
	"errors"
	"net/http"
	"server/internal/auth"
	"server/internal/dto/get"
	"server/internal/dto/patch"
	"server/internal/dto/post"
	"server/internal/dto/put"
	"server/internal/models"
	"server/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkHandler struct {
	WorkService *services.WorkServices
}

func NewWorkHandler(workService *services.WorkServices) *WorkHandler {
	return &WorkHandler{WorkService: workService}
}

func (h *WorkHandler) GetAll(c *gin.Context) {
	var responses []get.WorkResponse
	var works []models.Work

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication."})
		return
	}

	q := c.Query("completed")
	if q != "" {
		completed, err := strconv.ParseBool(q)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'completed' parameter. Expected true or false."})
			return
		}

		if completed {
			works, err = h.WorkService.GetCompleted.GetCompleted(uid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve completed works."})
				return
			}

			for _, w := range works {
				response := get.WorkResponse{
					ID:          w.ID,
					Title:       w.Title,
					IsCompleted: w.IsCompleted,
					Description: w.Description,
					CompletedAt: w.CompletedAt,
					UpdatedAt:   w.UpdatedAt,
					CreatedAt:   w.CreatedAt,
				}
				responses = append(responses, response)
			}

			c.JSON(http.StatusOK, responses)
			return
		}
	}

	works, err := h.WorkService.GetAll.GetAllByUID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve works."})
		return
	}

	for _, w := range works {
		response := get.WorkResponse{
			ID:          w.ID,
			Title:       w.Title,
			IsCompleted: w.IsCompleted,
			Description: w.Description,
			CompletedAt: w.CompletedAt,
			UpdatedAt:   w.UpdatedAt,
			CreatedAt:   w.CreatedAt,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, responses)

	return
}

func (h *WorkHandler) Create(c *gin.Context) {
	var work *models.Work
	var request post.WorksRequest
	var response post.WorksResponse

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication."})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
	}

	work, err := h.WorkService.CreateWork.CreateWork(uid, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create work."})
		return
	}

	response.ID = work.ID
	response.Title = work.Title
	response.WorkUrl = work.WorkURL
	response.RawIndex = work.RawIndex
	response.StitchIndex = work.StitchIndex
	response.IsCompleted = work.IsCompleted
	response.Description = work.Description
	response.CompletedAt = work.CompletedAt
	response.UpdatedAt = work.UpdatedAt
	response.CreatedAt = work.CreatedAt

	c.JSON(http.StatusCreated, response)

	return
}

func (h *WorkHandler) GetByID(c *gin.Context) {
	var work *models.Work
	var response get.WorksIDResponse

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication."})
		return
	}

	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format."})
		return
	}

	work, err = h.WorkService.Get.GetByID(uid, id)
	if err != nil {
		if errors.Is(err, services.ErrWorkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work not found."})
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to access this work."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get work by id."})
		return
	}

	response.ID = work.ID
	response.Title = work.Title
	response.WorkUrl = work.WorkURL
	response.RawIndex = work.RawIndex
	response.StitchIndex = work.RawIndex
	response.IsCompleted = work.IsCompleted
	response.Description = work.Description
	response.CompletedAt = work.CompletedAt

	c.JSON(http.StatusOK, response)

	return
}

func (h *WorkHandler) PutByID(c *gin.Context) {
	var work *models.Work
	var request put.WorksIDRequest
	var response put.WorksIDResponse

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication"})
		return
	}

	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format."})
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	work, err = h.WorkService.Put.PutByID(uid, id, &request)
	if err != nil {
		if errors.Is(err, services.ErrWorkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work not found."})
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to access this work."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get work by id."})
		return
	}

	response.ID = work.ID
	response.Title = work.Title
	response.WorkUrl = work.WorkURL
	response.RawIndex = work.RawIndex
	response.StitchIndex = work.StitchIndex
	response.IsCompleted = work.IsCompleted
	response.Description = work.Description
	response.CompletedAt = work.CompletedAt
	response.UpdatedAt = work.UpdatedAt
	response.CompletedAt = work.CompletedAt

	c.JSON(http.StatusOK, response)

	return
}

func (h *WorkHandler) PatchByID(c *gin.Context) {
	var work *models.Work
	var request patch.WorksIDRequest
	var response patch.WorksIDResponse

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication"})
		return
	}

	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format."})
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	work, err = h.WorkService.Patch.PatchByID(uid, id, &request)
	if err != nil {
		if errors.Is(err, services.ErrWorkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work not found."})
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to access this work."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get work by id."})
		return
	}

	response.ID = work.ID
	response.Title = work.Title
	response.WorkUrl = work.WorkURL
	response.RawIndex = work.RawIndex
	response.StitchIndex = work.StitchIndex
	response.IsCompleted = work.IsCompleted
	response.Description = work.Description
	response.CompletedAt = work.CompletedAt
	response.UpdatedAt = work.UpdatedAt
	response.CompletedAt = work.CompletedAt

	c.JSON(http.StatusOK, response)

	return
}

func (h *WorkHandler) DeleteByID(c *gin.Context) {
	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication."})
		return
	}

	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format."})
		return
	}

	if err = h.WorkService.Delete.DeleteByID(uid, id); err != nil {
		if errors.Is(err, services.ErrWorkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work not found."})
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to access this work."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete work by id."})
		return
	}

	c.JSON(http.StatusNoContent, nil)

	return
}
