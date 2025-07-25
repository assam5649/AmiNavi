package handler

import (
	firebase "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"server/internal/auth"
	"server/internal/dto/patch"
	"server/internal/dto/post"
	"server/internal/models"
	"server/internal/services"
)

type AuthHandler struct {
	DB           *gorm.DB
	FirebaseAuth *firebase.Client
}

func NewAuthHandler(db *gorm.DB, firebaseAuthClient *firebase.Client) *AuthHandler {
	return &AuthHandler{DB: db, FirebaseAuth: firebaseAuthClient}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var register post.RegisterResponse
	var user *models.User

	uid, ok := auth.GetUIDFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication."})
		return
	}

	if c.Request.ContentLength > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is not allowed."})
		return
	}

	user, exists, err := services.RegisterIfNotExists(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User register failed."})
		return
	}

	register.ID = user.ID
	register.DisplayName = user.DisplayName
	register.CreatedAt = user.CreatedAt

	if exists {
		c.JSON(http.StatusOK, register)
		return
	}

	c.JSON(http.StatusCreated, register)

	return
}

func (h *AuthHandler) Update(c *gin.Context) {
	var user *models.User
	var request patch.UpdateRequest
	var response patch.UpdateResponse

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication."})
		return
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if reflect.DeepEqual(request, patch.UpdateRequest{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err = services.Update(uid, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User update failed."})
		return
	}

	response.ID = user.ID
	response.DisplayName = user.DisplayName
	response.UpdatedAt = user.UpdatedAt

	c.JSON(http.StatusOK, response)

	return
}
