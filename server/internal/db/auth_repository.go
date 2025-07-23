package db

import (
	"server/internal/dto/patch"
	"server/internal/models"
	"time"
)

func FindByUID(uid string) (models.User, bool, error) {
	var user models.User

	DB, err := Connect()
	if err != nil {
		return models.User{}, false, err
	}

	result := DB.Where("firebase_uid = ?", uid).Find(&user)
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

func Update(uid string, request *patch.UpdateRequest) (int, error, time.Time) {
	var user models.User
	user.DisplayName = request.DisplayName
	user.ProfileImageURL = request.ProfileImageURL

	DB, err := Connect()
	if err != nil {
		return 0, err, time.Time{}
	}

	result := DB.Model(&user).Where("firebase_uid = ?", uid).Updates(models.User{DisplayName: user.DisplayName, ProfileImageURL: user.ProfileImageURL})
	if result.Error != nil {
		return 0, result.Error, time.Time{}
	}

	result = DB.First(&user, "firebase_uid = ?", uid)
	if result.Error != nil {
		return 0, result.Error, time.Time{}
	}
	return user.ID, nil, user.UpdatedAt
}
