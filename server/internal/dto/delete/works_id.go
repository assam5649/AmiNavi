package delete

import "time"

type WorksIDResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	DeletedAt time.Time `json:"deleted_at"`
}
