package patch

import "time"

type UpdateRequest struct {
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

type UpdateResponse struct {
	Message   string    `json:"message"`
	ID        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}
