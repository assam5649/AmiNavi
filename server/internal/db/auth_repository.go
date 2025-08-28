package db

import (
	"gorm.io/gorm"
	"server/internal/models"
)

func FindByUID(DB *gorm.DB, uid string) (*models.User, bool, error) {
	var user models.User

	result := DB.Where("firebase_uid = ?", uid).Find(&user)
	if result.Error != nil {
		return nil, false, result.Error
	}

	if result.RowsAffected > 0 {
		return &user, true, nil
	}

	return nil, false, nil
}

func CreateUser(DB *gorm.DB, uid string) (*models.User, error) {
	var user models.User
	user.DisplayName = "Amip"
	user.FirebaseUID = uid

	result := DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func Update(DB *gorm.DB, uid string, user *models.User) error {
	result := DB.Model(user).Where("firebase_uid = ?", uid).Updates(models.User{DisplayName: user.DisplayName, ProfileImageURL: user.ProfileImageURL})
	if result.Error != nil {
		return result.Error
	}

	result = DB.First(&user, "firebase_uid = ?", uid)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
