package db

import (
	"server/internal/models"
)

// Migrate はデータベースのスキーマを最新の状態に自動移行します。
// ここでは、定義されたモデルに基づいてテーブルを作成または更新します。
func Migrate() {
	DB.AutoMigrate(
		&models.User{},
		&models.Work{},
	)
}
