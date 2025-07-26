package models

import "time"

type Work struct {
	ID          int        `json:"id"`
	Author      string     `json:"author"`
	Title       string     `json:"title"`
	WorkURL     string     `json:"work_url"`
	RawIndex    int        `json:"raw_index"`
	StitchIndex int        `json:"stitch_index"`
	Bookmark    bool       `json:"bookmark"`
	IsCompleted bool       `json:"is_completed"`
	Description string     `json:"description"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
