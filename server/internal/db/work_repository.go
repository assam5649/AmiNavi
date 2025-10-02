package db

import (
	"errors"
	"server/internal/models"

	"gorm.io/gorm"
)

func GetByUID(DB *gorm.DB, uid string) ([]models.Work, error) {
	var work []models.Work

	result := DB.Where("author = ?", uid).Find(&work)
	if result.Error != nil {
		return []models.Work{}, result.Error
	}

	return work, nil
}

func GetCompleted(DB *gorm.DB, uid string) ([]models.Work, error) {
	var work []models.Work

	result := DB.Where("author = ? AND is_completed = ?", uid, true).Find(&work)
	if result.Error != nil {
		return []models.Work{}, result.Error
	}

	return work, nil
}

func CreateWork(DB *gorm.DB, work *models.Work) error {
	result := DB.Create(work)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetByID(DB *gorm.DB, uid string, id int) (*models.Work, error) {
	var work models.Work

	result := DB.Where("id = ? AND author = ?", id, uid).Find(&work)
	if result.Error != nil {
		return nil, result.Error
	}

	if work.ID == 0 && work.Author == "" {
		return nil, nil
	}

	return &work, nil
}

func SaveFileName(DB *gorm.DB, work *models.Work, id int, uid int, fileName string) error {
	updates := map[string]interface{}{
		"file_name": fileName,
	}

	result := DB.Model(&work).Where("id = ? AND author = ?", id, uid).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func PutByID(DB *gorm.DB, uid string, id int, work *models.Work) error {
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

func PatchByID(DB *gorm.DB, uid string, id int, work *models.Work) error {
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

func DeleteByID(DB *gorm.DB, uid string, id int, work *models.Work) (*gorm.DB, error) {
	result := DB.Where("id = ? AND author = ?", id, uid).First(&work)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	result = DB.Delete(&work)
	if result.Error != nil {
		return nil, result.Error
	}

	return result, nil
}
