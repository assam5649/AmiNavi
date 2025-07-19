package router

import (
	"net/http"
	"server/internal/auth"
	"server/internal/handler"

	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, firebaseAuthClient *firebaseAuth.Client) *gin.Engine {
	r := gin.Default()

	authHandler := handler.NewAuthHandler(db, firebaseAuthClient)

	api := r.Group("/v1")

	api.Use(auth.AuthMiddleware(firebaseAuthClient))
	{
		// 認証済みアクセス確認用テストエンドポイント
		// 有効なトークンがあればアクセス可能。GinのContextからUIDが取得できるかを確認。
		api.GET("/ping-auth", func(c *gin.Context) {
			uid, exists := auth.GetUIDFromContext(c)
			if !exists {
				// ここに到達した場合はミドルウェアかContext設定に問題
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase UID not found in context after authentication"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message":     "pong (authenticated access granted)",
				"firebaseUID": uid,
			})
		})

		r.POST("/register", authHandler.Register)
	}

	return r
}
