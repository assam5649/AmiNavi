package post

import "time"

type WorksRequest struct {
	Title       string `json:"title"`
	WorkUrl     string `json:"work_url"`
	Description string `json:"description"`
}

type WorksResponse struct {
	Message   string    `json:"message"`
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
