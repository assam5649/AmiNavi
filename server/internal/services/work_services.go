package services

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"server/internal/db"
	"server/internal/dto/patch"
	"server/internal/dto/post"
	"server/internal/dto/put"
	"server/internal/models"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4/auth"
	"gorm.io/gorm"
)

type GetAllService interface {
	GetAllByUID(uid string) ([]models.Work, error)
}

type GetCompletedService interface {
	GetCompleted(uid string) ([]models.Work, error)
}

type CreateWorkService interface {
	CreateWork(uid string, request *post.WorksRequest) (*models.Work, error)
}

type GetService interface {
	GetByID(uid string, id int) (*models.Work, string, error)
}

type PutService interface {
	PutByID(uid string, id int, request *put.WorksIDRequest) (*models.Work, error)
}

type PatchService interface {
	PatchByID(uid string, id int, request *patch.WorksIDRequest) (*models.Work, error)
}

type DeleteService interface {
	DeleteByID(uid string, id int) error
}

type GetAllServiceImpl struct {
	DB *gorm.DB
}
type GetCompletedServiceImpl struct {
	DB *gorm.DB
}
type CreateWorkServiceImpl struct {
	DB *gorm.DB
}
type GetServiceImpl struct {
	DB      *gorm.DB
	Storage *storage.Client
}
type PutServiceImpl struct {
	DB *gorm.DB
}
type PatchServiceImpl struct {
	DB *gorm.DB
}
type DeleteServiceImpl struct {
	DB *gorm.DB
}

type WorkServices struct {
	DB           *gorm.DB
	FirebaseAuth *firebase.Client
	Storage      *storage.Client
	GetAll       GetAllService
	GetCompleted GetCompletedService
	CreateWork   CreateWorkService
	Get          GetService
	Put          PutService
	Patch        PatchService
	Delete       DeleteService
}

func NewWorkService(
	db *gorm.DB,
	firebaseAuthClient *firebase.Client,
	getAll GetAllService,
	getCompleted GetCompletedService,
	put PutService,
	patch PatchService,
	delete DeleteService,
) *WorkServices {
	return &WorkServices{
		DB:           db,
		FirebaseAuth: firebaseAuthClient,
		GetAll:       getAll,
		GetCompleted: getCompleted,
		Put:          put,
		Patch:        patch,
		Delete:       delete,
	}
}

var (
	ErrWorkNotFound = errors.New("work not found")
	ErrForbidden    = errors.New("unauthorized access to work")
)

type ServiceAccount struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

func (s *GetAllServiceImpl) GetAllByUID(uid string) ([]models.Work, error) {
	workModel, err := db.GetByUID(s.DB, uid)
	if err != nil {
		return nil, err
	}

	return workModel, nil
}

func (s *GetCompletedServiceImpl) GetCompleted(uid string) ([]models.Work, error) {
	workModel, err := db.GetCompleted(s.DB, uid)
	if err != nil {
		return nil, err
	}

	return workModel, nil
}

func (s *CreateWorkServiceImpl) CreateWork(uid string, request *post.WorksRequest) (*models.Work, error) {
	var work = models.Work{
		Author:      uid,
		Title:       request.Title,
		FileName:    request.FileName,
		Description: request.Description,
		RawIndex:    0,
		StitchIndex: 0,
		IsCompleted: false,
		Bookmark:    false,
		CompletedAt: nil,
	}
	if err := db.CreateWork(s.DB, &work); err != nil {
		return nil, err
	}
	return &work, nil
}

func (s *GetServiceImpl) GetByID(uid string, id int) (*models.Work, string, error) {
	work, err := db.GetByID(s.DB, uid, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrWorkNotFound
		}
		if work.Author != uid {
			return nil, "", ErrForbidden
		}

		return nil, "", err
	}

	if work == nil {
		return nil, "", ErrWorkNotFound
	}

	bucketName := "aminavi"
	gcsFileName := "csv/" + work.FileName

	data, err := os.ReadFile("/app/august-button-470219-d5-5cc12900c9b1.json")
	if err != nil {
		return nil, "", err
	}

	var sa ServiceAccount
	if err := json.Unmarshal(data, &sa); err != nil {
		return nil, "", err
	}

	privateKey := []byte(strings.ReplaceAll(sa.PrivateKey, `\n`, "\n"))

	opts := &storage.SignedURLOptions{
		GoogleAccessID: sa.ClientEmail,
		PrivateKey:     privateKey,
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
		Scheme:         storage.SigningSchemeV4,
	}

	signedURL, err := storage.SignedURL(bucketName, gcsFileName, opts)
	if err != nil {
		return nil, "", err
	}
	slog.Info("Completed GCS URL generate")

	return work, signedURL, nil
}

func (s *PutServiceImpl) PutByID(uid string, id int, request *put.WorksIDRequest) (*models.Work, error) {
	var work = models.Work{
		Title:       request.Title,
		RawIndex:    request.RawIndex,
		StitchIndex: request.StitchIndex,
		IsCompleted: request.IsCompleted,
		Description: request.Description,
	}

	if err := db.PutByID(s.DB, uid, id, &work); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkNotFound
		}
		if work.Author != uid {
			return nil, ErrForbidden
		}

		return nil, err
	}

	return &work, nil
}

func (s *PatchServiceImpl) PatchByID(uid string, id int, request *patch.WorksIDRequest) (*models.Work, error) {
	var work = models.Work{
		RawIndex:    request.RawIndex,
		StitchIndex: request.RawIndex,
		IsCompleted: request.IsCompleted,
	}

	if err := db.PatchByID(s.DB, uid, id, &work); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkNotFound
		}
		if work.Author != uid {
			return nil, ErrForbidden
		}

		return nil, err
	}

	return &work, nil
}

func (s DeleteServiceImpl) DeleteByID(uid string, id int) error {
	var work = models.Work{}

	result, err := db.DeleteByID(s.DB, uid, id, &work)
	if err != nil {
		return err
	}

	if work.Author != uid {
		return ErrForbidden
	}

	if result.RowsAffected == 0 {
		return ErrWorkNotFound
	}

	return nil
}
