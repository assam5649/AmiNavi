package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect は環境変数を使用してデータベースへの接続を確立します。
// 最大10回までリトライを行い、接続が成功するかタイムアウトするまで待機します。
func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("DB connected")
			return DB, nil
		}
		log.Printf("DB connection failed. Retrying in 3 seconds... (%d/10): %v\n", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return DB, fmt.Errorf("failed to connect to database after retries: %w", err)
}
