package models

import "time"

type Works struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	WorkURL   string    `json:"work_url"`
	Count     int       `json:"count"`
	Bookmark  bool      `json:"bookmark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
