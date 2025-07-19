package handler

import (
	firebase "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"server/internal/auth"
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

// Register は新しいユーザーを登録します (POST /register)。
// クライアントから送られてきたFirebase IDトークンを検証し、その後にユーザー情報をデータベースに保存します。
// このエンドポイント自体に認証ミドルウェアは適用されません。
func (h *AuthHandler) Register(c *gin.Context) {
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
