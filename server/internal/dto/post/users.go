package post

import "time"

type RegisterResponse struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}
