package router

import (
	"server/internal/auth"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(authService *services.AuthServices, workService *services.WorkServices, mediaService *services.MediaServices) *gin.Engine {
	r := gin.Default()
	router := r.Group("/v1")

	router.Use(auth.AuthMiddleware(authService.FirebaseAuth))
	{
		AuthRouter(router, authService)
		WorkRouter(router, workService)
		MediaRouter(router, mediaService)
	}

	return r
}
