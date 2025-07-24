package services

import (
	"server/internal/db"
	"server/internal/dto/patch"
	"server/internal/models"
	"time"
)

func RegisterIfNotExists(uid string) (models.User, bool, error) {
	userModel, exists, err := db.FindByUID(uid)
	if exists {
		return userModel, exists, err
	}

	registerModel, err := db.CreateUser(uid)
	if err != nil {
		return models.User{}, exists, err
	}

	return registerModel, exists, nil
}

func Update(uid string, request *patch.UpdateRequest) (int, error, time.Time) {
	id, err, date := db.Update(uid, request)
	if err != nil {
		return 0, err, time.Time{}
	}

	return id, nil, date
}
