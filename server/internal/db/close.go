package db

import (
	"fmt"
	"log"
)

// Close はデータベース接続プールを安全にクローズします。
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying SQL DB instance: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		log.Println("INFO: Database connection closed.")
	}
	return nil
}
