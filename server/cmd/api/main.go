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
	"server/internal/services"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func main() {
	cfg := config.Load()
	credentialsFile := "/app/gcs.json"

	DB, err := db.Connect()
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("ERROR: Error closing database connection: %v", err)
		}
	}()
	firebaseAuthClient := auth.InitFirebaseAuthClient()

	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("client initialize failed: %v", err)
	}

	register := &services.RegisterServiceImpl{DB: DB}
	update := &services.UpdateServiceImpl{DB: DB}

	authService := &services.AuthServices{
		DB:           DB,
		FirebaseAuth: firebaseAuthClient,
		Register:     register,
		Update:       update,
	}

	getAll := &services.GetAllServiceImpl{DB: DB}
	getCompleted := &services.GetCompletedServiceImpl{DB: DB}
	createWork := &services.CreateWorkServiceImpl{DB: DB}
	get := &services.GetServiceImpl{DB: DB, Storage: gcsClient}
	put := &services.PutServiceImpl{DB: DB}
	patch := &services.PatchServiceImpl{DB: DB}
	deleteWork := &services.DeleteServiceImpl{DB: DB}

	workService := &services.WorkServices{
		DB:           DB,
		FirebaseAuth: firebaseAuthClient,
		Storage:      gcsClient,
		GetAll:       getAll,
		GetCompleted: getCompleted,
		CreateWork:   createWork,
		Get:          get,
		Put:          put,
		Patch:        patch,
		Delete:       deleteWork,
	}

	upload := &services.UploadServiceImpl{DB: DB, Storage: gcsClient}
	request := &services.RequestConvertServiceImpl{DB: DB}
	requestOCR := &services.RequestOCRServiceImpl{DB: DB}
	csvService := &services.CsvConversionServiceImpl{
		RequestService: request,
		UploadService:  upload,
	}
	ocrService := &services.OCRServiceImpl{
		RequestService: requestOCR,
	}
	fixCsv := &services.FixCsvServiceImpl{DB: DB, Storage: gcsClient}

	mediaService := &services.MediaServices{
		DB:           DB,
		FirebaseAuth: firebaseAuthClient,
		Storage:      gcsClient,
		CsvService:   csvService,
		Upload:       upload,
		Request:      request,
		OCR:          ocrService,
		RequestOCR:   requestOCR,
		FixCsv:       fixCsv,
	}

	r := router.NewRouter(authService, workService, mediaService)

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
