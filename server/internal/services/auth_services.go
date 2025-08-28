package services

import (
	firebase "firebase.google.com/go/v4/auth"
	"gorm.io/gorm"
	"server/internal/db"
	"server/internal/dto/patch"
	"server/internal/models"
)

type RegisterService interface {
	RegisterIfNotExists(uid string) (*models.User, bool, error)
}

type UpdateService interface {
	Update(uid string, request *patch.UpdateRequest) (*models.User, error)
}

type RegisterServiceImpl struct {
	DB *gorm.DB
}
type UpdateServiceImpl struct {
	DB *gorm.DB
}

type AuthServices struct {
	DB           *gorm.DB
	FirebaseAuth *firebase.Client
	Register     RegisterService
	Update       UpdateService
}

func NewAuthService(db *gorm.DB, firebaseAuthClient *firebase.Client, register RegisterService, update UpdateService) *AuthServices {
	return &AuthServices{
		DB:           db,
		FirebaseAuth: firebaseAuthClient,
		Register:     register,
		Update:       update,
	}
}

func (s *RegisterServiceImpl) RegisterIfNotExists(uid string) (*models.User, bool, error) {
	user, exists, err := db.FindByUID(s.DB, uid)
	if exists {
		return user, exists, err
	}

	registerModel, err := db.CreateUser(s.DB, uid)
	if err != nil {
		return nil, exists, err
	}

	return registerModel, exists, nil
}

func (s *UpdateServiceImpl) Update(uid string, request *patch.UpdateRequest) (*models.User, error) {
	var user models.User
	user.DisplayName = request.DisplayName
	user.ProfileImageURL = request.ProfileImageURL

	if err := db.Update(s.DB, uid, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
