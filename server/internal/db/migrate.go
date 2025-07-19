package db

import (
	"server/internal/models"
)

func Migrate() {
	DB.AutoMigrate(
		&models.User{},
		&models.Work{},
	)
}
