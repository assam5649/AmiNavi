package models

import "time"

// Work はデータベースの 'works' テーブルに対応する構造体です。
type Work struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	WorkURL   string    `json:"work_url"`
	Count     int       `json:"count"`
	Bookmark  bool      `json:"bookmark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
