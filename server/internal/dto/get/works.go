package get

import "time"

type WorkResponse struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	IsCompleted bool       `json:"is_completed"`
	Description string     `json:"description"`
	CompletedAt *time.Time `json:"completed_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
