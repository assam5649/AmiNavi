package router

import (
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/internal/handler"
)

func WorkRouter(router *gin.RouterGroup, db *gorm.DB, firebaseAuthClient *firebaseAuth.Client) {
	workHandler := handler.NewWorkHandler(db, firebaseAuthClient)

	router.GET("/works", workHandler.GetAll)
	router.POST("/works", workHandler.Create)

	router.GET("/works/:id", workHandler.GetByID)
	router.PUT("/works/:id", workHandler.PutByID)
	router.PATCH("/works/:id", workHandler.PatchByID)
	router.DELETE("/works/:id", workHandler.DeleteByID)
}
