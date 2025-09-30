package services

import (
	"errors"
	"server/internal/db"
	"server/internal/dto/patch"
	"server/internal/dto/post"
	"server/internal/dto/put"
	"server/internal/models"

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
	GetByID(uid string, id int) (*models.Work, error)
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
	DB *gorm.DB
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
		WorkURL:     request.WorkUrl,
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

func (s *GetServiceImpl) GetByID(uid string, id int) (*models.Work, error) {
	work, err := db.GetByID(s.DB, uid, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkNotFound
		}
		if work.Author != uid {
			return nil, ErrForbidden
		}

		return nil, err
	}

	if work == nil {
		return nil, ErrWorkNotFound
	}

	return work, nil
}

func (s *PutServiceImpl) PutByID(uid string, id int, request *put.WorksIDRequest) (*models.Work, error) {
	var work = models.Work{
		Title:       request.Title,
		WorkURL:     request.WorkUrl,
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

	if work.Author != uid {
		return ErrForbidden
	}

	result, err := db.DeleteByID(s.DB, uid, id, &work)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ErrWorkNotFound
	}

	return nil
}
