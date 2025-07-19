package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/internal/auth"
	"server/internal/config"
	"server/internal/db"
	"server/internal/router"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	_, err := db.Connect()
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("ERROR: Error closing database connection: %v", err)
		}
	}()
	firebaseAuthClient := auth.InitFirebaseAuthClient()

	r := router.NewRouter(db.DB, firebaseAuthClient)

	serverAddr := ":" + cfg.Server.Port
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		log.Printf("INFO: Server starting on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("INFO: Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("FATAL: Server forced to shutdown: %v", err)
	}

	log.Println("INFO: Server exited.")
}
