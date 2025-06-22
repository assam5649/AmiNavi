package db_test

import (
	"os"
	"testing"

	"server/internal/db"
)

func TestConnect(t *testing.T) {
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASS", "testpass")
	os.Setenv("DB_HOST", "db")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "testdb")

	err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
}
