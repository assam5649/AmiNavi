package config

import (
	"log"
	"os"
	"strconv"
)

// Config はアプリケーション全体の設値定義です。
// 各サービスやコンポーネメントが必要とする設定をここに集約します。
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig はHTTPサーバー関連の設定です。
type ServerConfig struct {
	Port string
}

// DatabaseConfig はデータベース接続関連の設定です。
type DatabaseConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

// Load は環境変数からアプリケーションの設定をロードし、Config構造体へのポインタを返します。
// 環境変数が設定されていない場合、デフォルト値を使用するか、致命的なエラーを発生させます。
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("WARN: PORT 環境変数が設定されていません。デフォルト値 '%s' を使用します。", port)
	}

	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		dbPortStr = "3306"
		log.Printf("WARN: DB_PORT 環境変数が設定されていません。デフォルト値 '%s' を使用します。", dbPortStr)
	}
	_, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("FATAL: DB_PORT 環境変数が無効な数値です: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatalf("FATAL: DB_USER 環境変数が設定されていません。")
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		log.Fatalf("FATAL: DB_PASS 環境変数が設定されていません。")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatalf("FATAL: DB_HOST 環境変数が設定されていません。")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalf("FATAL: DB_NAME 環境変数が設定されていません。")
	}

	// ロードした設定値をConfig構造体にまとめて返す
	return &Config{
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			User: dbUser,
			Pass: dbPass,
			Host: dbHost,
			Port: dbPortStr,
			Name: dbName,
		},
	}
}
