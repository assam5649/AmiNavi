package handler

import "C"
import (
	firebase "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"server/internal/auth"
	"server/internal/dto/get"
	"server/internal/dto/post"
	"server/internal/models"
	"server/internal/services"
	"strconv"
)

type WorkHandler struct {
	DB           *gorm.DB
	FirebaseAuth *firebase.Client
}

func NewWorkHandler(db *gorm.DB, firebaseAuthClient *firebase.Client) *WorkHandler {
	return &WorkHandler{DB: db, FirebaseAuth: firebaseAuthClient}
}

func (h *WorkHandler) GetAll(c *gin.Context) {
	var responses []get.WorkResponse
	var works []models.Work

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication"})
		return
	}

	q := c.Query("completed")
	if q != "" {
		completed, err := strconv.ParseBool(q)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if completed {
			works, err = services.GetCompleted(uid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}

			for _, w := range works {
				response := get.WorkResponse{
					Title:       w.Title,
					IsCompleted: w.IsCompleted,
					Description: w.Description,
					CompletedAt: w.CompletedAt,
				}
				responses = append(responses, response)
			}
			c.JSON(http.StatusOK, responses)
			return
		}
	}

	works, err := services.GetAllByUID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	for _, w := range works {
		response := get.WorkResponse{
			Title:       w.Title,
			IsCompleted: w.IsCompleted,
			Description: w.Description,
			CompletedAt: w.CompletedAt,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, responses)

	return
}

func (h *WorkHandler) Create(c *gin.Context) {
	var request post.WorksRequest
	var response post.WorksResponse
	var work models.Work

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	}

	if err := services.CreateWork(uid, &request); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	response.Message = "The work has been successfully created."
	response.Id = work.ID
	response.CreatedAt = work.CreatedAt

	c.JSON(http.StatusCreated, response)

	return
}

func (h *WorkHandler) GetByID(c *gin.Context) {

}

func (h *WorkHandler) PutByID(c *gin.Context) {

}

func (h *WorkHandler) PatchByID(c *gin.Context) {

}

func (h *WorkHandler) DeleteByID(c *gin.Context) {

}
