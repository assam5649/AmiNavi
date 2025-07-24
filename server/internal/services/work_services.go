package services

import (
	"server/internal/db"
	"server/internal/dto/post"
	"server/internal/models"
)

func GetAllByUID(uid string) ([]models.Work, error) {
	workModel, err := db.GetByUID(uid)
	if err != nil {
		return []models.Work{}, err
	}

	return workModel, nil
}

func GetCompleted(uid string) ([]models.Work, error) {
	workModel, err := db.GetCompleted(uid)
	if err != nil {
		return []models.Work{}, err
	}

	return workModel, nil
}

func CreateWork(uid string, work *post.WorksRequest) error {
	if err := db.CreateWork(uid, work); err != nil {
		return err
	}
	return nil
}
