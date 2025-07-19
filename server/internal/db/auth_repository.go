package db

import (
	"server/internal/models"
)

func FindByUID(uid string) (models.User, bool, error) {
	var user models.User

	DB, err := Connect()
	if err != nil {
		return models.User{}, false, err
	}

	result := DB.Where("uid = ?", uid).Find(&user)
	if result.Error != nil {
		return models.User{}, false, result.Error
	}

	if result.RowsAffected > 0 {
		return user, true, nil
	}

	return models.User{}, false, nil
}

func Create(uid string) (models.User, error) {
	var user models.User
	user.DisplayName = "Amip"
	user.FirebaseUID = uid

	DB, err := Connect()
	if err != nil {
		return models.User{}, err
	}

	result := DB.Create(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}
