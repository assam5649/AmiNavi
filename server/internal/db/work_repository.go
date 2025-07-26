package db

import (
	"gorm.io/gorm"
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

func CreateWork(work *models.Work) error {
	DB, err := Connect()
	if err != nil {
		return err
	}

	result := DB.Create(work)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetByID(uid string, id int) (*models.Work, error) {
	var work models.Work

	DB, err := Connect()
	if err != nil {
		return nil, err
	}

	result := DB.Where("id = ? AND author = ?", id, uid).Find(&work)
	if result.Error != nil {
		return nil, result.Error
	}

	return &work, nil
}

func PutByID(uid string, id int, work *models.Work) error {
	DB, err := Connect()
	if err != nil {
		return err
	}

	result := DB.Model(work).Where("id = ? AND author = ?", id, uid).Updates(work)
	if result.Error != nil {
		return result.Error
	}

	result = DB.First(work, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func PatchByID(uid string, id int, work *models.Work) error {
	DB, err := Connect()
	if err != nil {
		return err
	}

	result := DB.Model(work).Where("id = ? AND author = ?", id, uid).Updates(work)
	if result.Error != nil {
		return result.Error
	}

	result = DB.First(work, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteByID(uid string, id int, work *models.Work) (*gorm.DB, error) {
	DB, err := Connect()
	if err != nil {
		return nil, err
	}

	result := DB.Where("id = ? AND author = ?", id, uid).Delete(work)
	if result.Error != nil {
		return nil, result.Error
	}

	return result, nil
}
