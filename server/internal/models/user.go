package models

import "time"

// User はデータベースの 'users' テーブルに対応する構造体です。
type User struct {
	ID              int       `json:"id"`
	FirebaseUID     string    `json:"firebase_uid"`
	LoginID         string    `json:"login_id"`
	DisplayName     string    `json:"display_name"`
	ProfileImageURL string    `json:"profile_image_url"`
	EMail           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
