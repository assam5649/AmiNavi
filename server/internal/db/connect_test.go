package db

import (
	"os"
	"testing"
)

// TestConnect はデータベース接続関数 Connect() の動作をテストします。
// テスト実行時に必要な環境変数を一時的に設定し、接続が成功するか検証します。
func TestConnect(t *testing.T) {
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASS", "testpass")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "testdb")

	err := Connect()
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
}
