package db

import (
	"server/internal/dto/post"
	"server/internal/models"
	"time"
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

func GetByID(uid string, id int) (models.Work, error) {
	var work models.Work

	DB, err := Connect()
	if err != nil {
		return models.Work{}, err
	}

	result := DB.Where("id = ? AND author = ?", id, uid).Find(&work)
	if result.Error != nil {
		return models.Work{}, result.Error
	}

	return work, nil
}

func PutByID(uid string, id int, work *models.Work) (int, time.Time, error) {
	DB, err := Connect()
	if err != nil {
		return 0, time.Time{}, err
	}

	result := DB.Model(work).Where("id = ? AND author = ?", id, uid).Updates(work)
	if result.Error != nil {
		return 0, time.Time{}, result.Error
	}

	result = DB.First(work, work.ID)
	if result.Error != nil {
		return 0, time.Time{}, result.Error
	}

	return work.ID, work.UpdatedAt, nil
}

func PatchByID(uid string, id int, work *models.Work) (int, time.Time, error) {
	DB, err := Connect()
	if err != nil {
		return 0, time.Time{}, err
	}

	result := DB.Model(work).Where("id = ? AND author = ?", id, uid).Updates(work)
	if result.Error != nil {
		return 0, time.Time{}, result.Error
	}

	result = DB.First(work, work.ID)
	if result.Error != nil {
		return 0, time.Time{}, result.Error
	}

	return work.ID, work.UpdatedAt, nil
}

func DeleteByID(uid string, id int, work *models.Work) (int, string, time.Time, error) {
	DB, err := Connect()
	if err != nil {
		return 0, "", time.Time{}, err
	}

	result := DB.Where("id = ? AND author = ?", id, uid).Delete(work)
	if result.Error != nil {
		return 0, "", time.Time{}, result.Error
	}

	now := time.Now()

	return work.ID, work.Title, now, nil
}
