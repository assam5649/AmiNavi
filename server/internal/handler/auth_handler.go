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

func (h *AuthHandler) UserRegister(c *gin.Context) {
	var register post.RegisterResponse
	var user models.User

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication"})
		return
	}

	user, exists, err := services.RegisterIfNotExists(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	register.ID = user.ID
	register.DisplayName = user.DisplayName
	register.CreatedAt = user.CreatedAt

	c.JSON(http.StatusCreated, register)

	return
}
func (h *AuthHandler) UserUpdate(c *gin.Context) {
	var update patch.UpdateRequest
	var response patch.UpdateResponse

	uid, exists := auth.GetUIDFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication"})
		return
	}

	err := c.ShouldBindJSON(&update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if reflect.DeepEqual(update, patch.UpdateRequest{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err, date := services.Update(uid, &update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	response.Message = "Patch Successfully"
	response.ID = id
	response.UpdatedAt = date

	c.JSON(http.StatusOK, response)

	return
}
