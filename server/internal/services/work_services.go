package services

import (
	"errors"
	"gorm.io/gorm"
	"server/internal/db"
	"server/internal/dto/patch"
	"server/internal/dto/post"
	"server/internal/dto/put"
	"server/internal/models"
)

var (
	ErrWorkNotFound = errors.New("work not found")
	ErrForbidden    = errors.New("unauthorized access to work")
)

func GetAllByUID(uid string) ([]models.Work, error) {
	workModel, err := db.GetByUID(uid)
	if err != nil {
		return nil, err
	}

	return workModel, nil
}

func GetCompleted(uid string) ([]models.Work, error) {
	workModel, err := db.GetCompleted(uid)
	if err != nil {
		return nil, err
	}

	return workModel, nil
}

func CreateWork(uid string, request *post.WorksRequest) (*models.Work, error) {
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
	if err := db.CreateWork(&work); err != nil {
		return nil, err
	}
	return &work, nil
}

func GetByID(uid string, id int) (*models.Work, error) {
	work, err := db.GetByID(uid, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkNotFound
		}
		if work.Author != uid {
			return nil, ErrForbidden
		}

		return nil, err
	}

	return work, nil
}

func PutByID(uid string, id int, request *put.WorksIDRequest) (*models.Work, error) {
	var work = models.Work{
		Title:       request.Title,
		WorkURL:     request.WorkUrl,
		RawIndex:    request.RawIndex,
		StitchIndex: request.StitchIndex,
		IsCompleted: request.IsCompleted,
		Description: request.Description,
	}

	if err := db.PutByID(uid, id, &work); err != nil {
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

func PatchByID(uid string, id int, request *patch.WorksIDRequest) (*models.Work, error) {
	var work = models.Work{
		RawIndex:    request.RawIndex,
		StitchIndex: request.RawIndex,
		IsCompleted: request.IsCompleted,
	}

	if err := db.PatchByID(uid, id, &work); err != nil {
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

func DeleteByID(uid string, id int) error {
	var work = models.Work{}

	if work.Author != uid {
		return ErrForbidden
	}

	result, err := db.DeleteByID(uid, id, &work)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ErrWorkNotFound
	}

	return nil
}
