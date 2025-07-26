package db

import (
	"server/internal/models"
)

func FindByUID(uid string) (*models.User, bool, error) {
	var user models.User

	DB, err := Connect()
	if err != nil {
		return nil, false, err
	}

	result := DB.Where("firebase_uid = ?", uid).Find(&user)
	if result.Error != nil {
		return nil, false, result.Error
	}

	if result.RowsAffected > 0 {
		return &user, true, nil
	}

	return nil, false, nil
}

func CreateUser(uid string) (*models.User, error) {
	var user models.User
	user.DisplayName = "Amip"
	user.FirebaseUID = uid

	DB, err := Connect()
	if err != nil {
		return nil, err
	}

	result := DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func Update(uid string, user *models.User) error {
	DB, err := Connect()
	if err != nil {
		return err
	}

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
