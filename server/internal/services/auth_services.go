package services

import (
	"server/internal/db"
	"server/internal/models"
)

func RegisterIfNotExists(uid string) (models.User, bool, error) {
	userModel, exists, err := db.FindByUID(uid)
	if exists {
		return userModel, exists, err
	}

	registerModel, err := db.Create(uid)
	if err != nil {
		return models.User{}, exists, err
	}

	return registerModel, exists, nil
}
