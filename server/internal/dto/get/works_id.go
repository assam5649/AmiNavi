package get

import "time"

type WorksIDResponse struct {
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	WorkUrl     string     `json:"work_url"`
	RawIndex    int        `json:"raw_index"`
	StitchIndex int        `json:"stitch_index"`
	IsCompleted bool       `json:"is_completed"`
	Description string     `json:"description"`
	CompletedAt *time.Time `json:"completed_at"`
}
