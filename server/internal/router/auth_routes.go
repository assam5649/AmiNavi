package router

import (
	"server/internal/handler"

	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRouter(router *gin.RouterGroup, db *gorm.DB, firebaseAuthClient *firebaseAuth.Client) {
	authHandler := handler.NewAuthHandler(db, firebaseAuthClient)

	router.POST("/users", authHandler.Register)
	router.PATCH("/users", authHandler.Update)
}
