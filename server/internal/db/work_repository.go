package db

import "server/internal/models"

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
