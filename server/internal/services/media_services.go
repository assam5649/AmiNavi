package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CsvConversionService interface {
	ConvertAndUpload(ctx context.Context, image []byte, requestService RequestServices, uploadService UploadServices) ([]byte, string, error)
}

type UploadServices interface {
	Upload(ctx context.Context, dataToUpload []byte, fileId string) (string, error)
}

type RequestServices interface {
	RequestConvert(image []byte) ([]byte, error)
}

type RequestOCRServices interface {
	RequestOCR(image []byte) ([]byte, error)
}

type OCRServices interface {
	OCR(ctx context.Context, image []byte) ([]byte, error)
}

type FixCsvService interface {
	FixCsv(ctx context.Context, csv []byte) ([]byte, string, error)
}

type CsvConversionServiceImpl struct {
	RequestService RequestServices
	UploadService  UploadServices
}
type UploadServiceImpl struct {
	DB      *gorm.DB
	Storage *storage.Client
}
type RequestConvertServiceImpl struct {
	DB *gorm.DB
}
type RequestOCRServiceImpl struct {
	DB *gorm.DB
}
type OCRServiceImpl struct {
	RequestService RequestOCRServices
	UploadService  UploadServices
}
type FixCsvServiceImpl struct {
	DB      *gorm.DB
	Storage *storage.Client
}

func (s *FixCsvServiceImpl) FixCsv(ctx context.Context,
	csv []byte,
) ([]byte, string, error) {
	slog.Info("FixCsv Started.")
	fileId := uuid.NewString()
	bucketName := "aminavi"
	gcsFileName := "csv/" + fileId

	bucket := s.Storage.Bucket(bucketName)

	obj := bucket.Object(gcsFileName)

	wc := obj.NewWriter(ctx)
	reader := bytes.NewReader(csv)

	if _, err := io.Copy(wc, reader); err != nil {
		slog.Error("Failed to IO Copy", "error", err, "bucket", bucketName, "object", gcsFileName)
	}

	if err := wc.Close(); err != nil {
		slog.Error("Failed to Close wc", "error", err, "bucket", bucketName, "object", gcsFileName)
	}

	slog.Info("Fix Csv Upload Completed.")
	return csv, fileId, nil
}

func (s *OCRServiceImpl) OCR(ctx context.Context,
	image []byte,
) ([]byte, error) {
	csv, err := s.RequestService.RequestOCR(image)
	if err != nil {
		return nil, err
	}

	slog.Info("OCR completed.")
	slog.Info("Response Data", "pattern", csv)
	return csv, err
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

	filename, err := uploadServices.Upload(ctx, csv, fileId)
	if err != nil {
		return nil, "", err
	}

	slog.Info("convertandupload completed.")
	return csv, filename, err
}

type MediaServices struct {
	DB           *gorm.DB
	FirebaseAuth *firebase.Client
	Storage      *storage.Client
	CsvService   CsvConversionService
	Upload       UploadServices
	Request      RequestServices
	OCR          OCRServices
	RequestOCR   RequestOCRServices
	FixCsv       FixCsvService
}

func NewMediaService(db *gorm.DB, firebaseAuthClient *firebase.Client, storage *storage.Client, csvService CsvConversionService, upload UploadServices, request RequestServices, OCR OCRServices, requestOCR RequestOCRServices, fixCsv FixCsvService) *MediaServices {
	return &MediaServices{DB: db, FirebaseAuth: firebaseAuthClient, Storage: storage, CsvService: csvService, Upload: upload, Request: request, OCR: OCR, RequestOCR: requestOCR, FixCsv: fixCsv}
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

	slog.Info("Upload Completed.")
	return gcsFileName, nil
}

func (s *RequestOCRServiceImpl) RequestOCR(image []byte) ([]byte, error) {
	slog.Info("RequestOCR started.")

	imageSize := len(image)
	slog.Info(fmt.Sprintf("DEBUG: Go received image data size: %d bytes", imageSize))

	if imageSize > 0 {
		endIndex := 20
		if imageSize < 20 {
			endIndex = imageSize
		}
		slog.Info(fmt.Sprintf("DEBUG: First bytes: %v", image[:endIndex]))
	} else {
		slog.Warn("DEBUG: Go received zero-length image data.")
	}

	url := "http://ml-service:8501/ocr"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "upload.png")
	if err != nil {
		return nil, err
	}
	written, err := part.Write(image)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("DEBUG: Written bytes to multipart form: %d", written))

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python server returned status %d", resp.StatusCode)
	}

	csvData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	slog.Info("RequestOCR Completed.")
	return csvData, nil
}

func (s *RequestConvertServiceImpl) RequestConvert(image []byte) ([]byte, error) {
	slog.Info("RequestConvert started.")

	imageSize := len(image)
	slog.Info(fmt.Sprintf("DEBUG: Go received image data size: %d bytes", imageSize))

	if imageSize > 0 {
		endIndex := 20
		if imageSize < 20 {
			endIndex = imageSize
		}
		slog.Info(fmt.Sprintf("DEBUG: First bytes: %v", image[:endIndex]))
	} else {
		slog.Warn("DEBUG: Go received zero-length image data.")
	}

	url := "http://ml-service:8501/convert"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "upload.png")
	if err != nil {
		return nil, err
	}
	written, err := part.Write(image)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("DEBUG: Written bytes to multipart form: %d", written))

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
