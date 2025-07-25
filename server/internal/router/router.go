package router

import (
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/internal/auth"
)

func NewRouter(db *gorm.DB, firebaseAuthClient *firebaseAuth.Client) *gin.Engine {
	r := gin.Default()
	router := r.Group("/v1")

	router.Use(auth.AuthMiddleware(firebaseAuthClient))
	{
		AuthRouter(router, db, firebaseAuthClient)
		WorkRouter(router, db, firebaseAuthClient)
	}

	return r
}
