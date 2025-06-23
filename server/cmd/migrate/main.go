package migrate

import (
	"log"
	"server/internal/db"
)

// Migrateはアプリケーションのデータベーススキーマを最新の状態に移行します。
// データベースへの接続を行い、その後にマイグレーション処理を実行します。
func Migrate() {
	err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.Migrate()
}
