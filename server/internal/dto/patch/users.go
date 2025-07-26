package patch

import "time"

type UpdateRequest struct {
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

type UpdateResponse struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"display_name"`
	UpdatedAt   time.Time `json:"updated_at"`
}
