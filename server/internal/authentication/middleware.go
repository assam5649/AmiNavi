package auth

import (
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware はFirebase IDトークンを検証するGinミドルウェアです。
// authClient: Firebase Admin SDKの認証クライアント
// 戻り値: Ginのミドルウェア関数 (gin.HandlerFunc)
func AuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		idToken := strings.Replace(authHeader, "Bearer ", "", 1)
		if idToken == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in 'Bearer <token>' format"})
			return
		}

		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			log.Printf("Firebase ID token verification failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired authentication token"})
			return
		}

		c.Set("firebaseUID", token.UID)

		c.Next()
	}
}

// GetUIDFromContext はGinのContextからFirebase UIDを取得するヘルパー関数です。
// AuthMiddlewareが適用されているエンドポイントで利用することを想定しています。
func GetUIDFromContext(c *gin.Context) (string, bool) {
	uid, ok := c.Get("firebaseUID")
	if !ok {
		return "", false
	}
	strUID, ok := uid.(string)
	if !ok {
		return "", false
	}
	return strUID, true
}
