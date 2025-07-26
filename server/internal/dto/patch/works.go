package patch

import "time"

type WorksIDRequest struct {
	Title       string `json:"title"`
	WorkUrl     string `json:"work_url"`
	RawIndex    int    `json:"raw_index"`
	StitchIndex int    `json:"stitch_index"`
	IsCompleted bool   `json:"is_completed"`
	Description string `json:"description"`
}

type WorksIDResponse struct {
	Message   string    `json:"message"`
	ID        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}
