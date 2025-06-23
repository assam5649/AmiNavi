package migrate

import (
	"server/db"
)

// Migrateはアプリケーションのデータベーススキーマを最新の状態に移行します。
// データベースへの接続を行い、その後にマイグレーション処理を実行します。
func Migrate() {
	db.Connect()
	db.Migrate()
}
