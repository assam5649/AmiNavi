package put

import "time"

type WorksIDRequest struct {
	Title       string `json:"title"`
	RawIndex    int    `json:"raw_index"`
	StitchIndex int    `json:"stitch_index"`
	IsCompleted bool   `json:"is_completed"`
	Description string `json:"description"`
}

type WorksIDResponse struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	RawIndex    int        `json:"raw_index"`
	StitchIndex int        `json:"stitch_index"`
	IsCompleted bool       `json:"is_completed"`
	Description string     `json:"description"`
	CompletedAt *time.Time `json:"completed_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
