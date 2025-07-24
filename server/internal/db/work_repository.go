package db

import (
	"server/internal/dto/post"
	"server/internal/models"
)

func GetByUID(uid string) ([]models.Work, error) {
	var work []models.Work

	DB, err := Connect()
	if err != nil {
		return []models.Work{}, err
	}

	result := DB.Where("author = ?", uid).Find(&work)
	if result.Error != nil {
		return []models.Work{}, result.Error
	}

	return work, nil
}

func GetCompleted(uid string) ([]models.Work, error) {
	var work []models.Work

	DB, err := Connect()
	if err != nil {
		return []models.Work{}, err
	}

	result := DB.Where("author = ? AND is_completed = ?", uid, true).Find(&work)
	if result.Error != nil {
		return []models.Work{}, result.Error
	}

	return work, nil
}

func CreateWork(uid string, request *post.WorksRequest) error {
	var work models.Work

	work.Author = uid
	work.Title = request.Title
	work.WorkURL = request.WorkUrl
	work.Description = request.Description
	work.IsCompleted = false
	work.Bookmark = false
	work.RawIndex = 0
	work.StitchIndex = 0
	work.CompletedAt = nil

	DB, err := Connect()
	if err != nil {
		return err
	}

	result := DB.Create(&work)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
