package handler

import (
	"log"
	"net/http"

	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server/internal/models"
)

type AuthHandler struct {
	DB           *gorm.DB
	FirebaseAuth *auth.Client
}

func NewAuthHandler(db *gorm.DB, firebaseAuthClient *auth.Client) *AuthHandler {
	return &AuthHandler{DB: db, FirebaseAuth: firebaseAuthClient}
}

type RegisterRequest struct {
	IDToken         string `json:"id_token" binding:"required"`
	LoginID         string `json:"login_id" binding:"required,alphanum"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
	EMail           string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: Invalid request body for /register: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.FirebaseAuth.VerifyIDToken(c.Request.Context(), req.IDToken)
	if err != nil {
		log.Printf("WARN: Firebase ID token verification failed for /register: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired Firebase ID token"})
		return
	}
	firebaseUID := token.UID

	var existingUser models.User

	if h.DB.Where("firebase_uid = ?", firebaseUID).First(&existingUser).Error == nil {
		log.Printf("WARN: Attempted to register existing Firebase UID: %s", firebaseUID)
		c.JSON(http.StatusConflict, gin.H{"error": "User with this Firebase UID already exists"})
		return
	}

	if h.DB.Where("login_id = ?", req.LoginID).First(&existingUser).Error == nil {
		log.Printf("WARN: Attempted to register existing LoginID: %s", req.LoginID)
		c.JSON(http.StatusConflict, gin.H{"error": "LoginID already taken"})
		return
	}

	newUser := models.User{
		FirebaseUID:     firebaseUID,
		LoginID:         req.LoginID,
		DisplayName:     req.DisplayName,
		ProfileImageURL: req.ProfileImageURL,
		EMail:           req.EMail,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if result := h.DB.Create(&newUser); result.Error != nil {
		log.Printf("ERROR: Failed to create user in DB: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user due to database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "User registered successfully",
		"user_id":      newUser.ID,
		"firebase_uid": newUser.FirebaseUID,
	})
}
