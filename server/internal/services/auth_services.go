package services

import (
	"server/internal/db"
	"server/internal/dto/patch"
	"server/internal/models"
)

func RegisterIfNotExists(uid string) (*models.User, bool, error) {
	user, exists, err := db.FindByUID(uid)
	if exists {
		return user, exists, err
	}

	registerModel, err := db.CreateUser(uid)
	if err != nil {
		return nil, exists, err
	}

	return registerModel, exists, nil
}

func Update(uid string, request *patch.UpdateRequest) (*models.User, error) {
	var user models.User
	user.DisplayName = request.DisplayName
	user.ProfileImageURL = request.ProfileImageURL

	if err := db.Update(uid, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
