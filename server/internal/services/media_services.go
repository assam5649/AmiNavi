package services

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"net/http"
)

type CsvConversionService interface {
	ConvertAndUpload(ctx context.Context, image []byte, requestService RequestServices, uploadService UploadServices) ([]byte, string, error)
}

type UploadServices interface {
	Upload(ctx context.Context, image []byte, fileId string) (string, error)
}

type RequestServices interface {
	RequestConvert(image []byte) ([]byte, error)
}

type CsvConversionServiceImpl struct{}
type UploadServiceImpl struct {
	DB      *gorm.DB
	Storage *storage.Client
}
type RequestConvertServiceImpl struct {
	DB *gorm.DB
}

func (s *CsvConversionServiceImpl) ConvertAndUpload(ctx context.Context,
	image []byte,
	requestService RequestServices,
	uploadServices UploadServices,
) ([]byte, string, error) {
	fileId := uuid.NewString()
	csv, err := requestService.RequestConvert(image)
	if err != nil {
		return nil, "", err
	}

	gcsPath, err := uploadServices.Upload(ctx, image, fileId)
	if err != nil {
		return nil, "", err
	}

	slog.Info("ConvertAndUpload Completed.")
	return csv, gcsPath, err
}

type MediaServices struct {
	DB           *gorm.DB
	FirebaseAuth *firebase.Client
	Storage      *storage.Client
	CsvService   CsvConversionService
	Upload       UploadServices
	Request      RequestServices
}

func NewMediaService(db *gorm.DB, firebaseAuthClient *firebase.Client, storage *storage.Client, csvService CsvConversionService, upload UploadServices, request RequestServices) *MediaServices {
	return &MediaServices{DB: db, FirebaseAuth: firebaseAuthClient, Storage: storage, CsvService: csvService, Upload: upload, Request: request}
}

func (s *UploadServiceImpl) Upload(ctx context.Context, csvData []byte, fileId string) (string, error) {
	bucketName := "aminavi"
	gcsFileName := "csv/" + fileId

	bucket := s.Storage.Bucket(bucketName)

	obj := bucket.Object(gcsFileName)

	wc := obj.NewWriter(ctx)
	reader := bytes.NewReader(csvData)

	if _, err := io.Copy(wc, reader); err != nil {
		slog.Error("Failed to IO Copy", "error", err, "bucket", bucketName, "object", gcsFileName)
	}

	if err := wc.Close(); err != nil {
		slog.Error("Failed to Close wc", "error", err, "bucket", bucketName, "object", gcsFileName)
	}

	gcsPath := "gs://" + bucketName + "/" + gcsFileName

	slog.Info("Upload Completed.")
	return gcsPath, nil
}

func (s *RequestConvertServiceImpl) RequestConvert(image []byte) ([]byte, error) {
	url := "http://ml-service:8501/convert"

	resp, err := http.Post(url, "application/octet-stream", bytes.NewReader(image))
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				return
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python server returned status %d", resp.StatusCode)
	}

	csvData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	slog.Info("RequestConvert Completed.")
	return csvData, nil
}
