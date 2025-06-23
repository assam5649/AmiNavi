package db

import (
	"server/models"
)

func Migrate() {
	DB.AutoMigrate(
		&models.Users{},
		&models.Works{},
	)
}
