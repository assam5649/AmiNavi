package get

import "time"

type WorkResponse struct {
	Title       string     `json:"title"`
	IsCompleted bool       `json:"is_completed"`
	Description string     `json:"description"`
	CompletedAt *time.Time `json:"completed_at"`
}
