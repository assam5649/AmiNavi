package put

type Register struct {
	LoginID         string `json:"login_id" binding:"required"`
	DisplayName     string `json:"display_name" binding:"required,alphanum"`
	ProfileImageURL string `json:"profile_image_url"`
	Email           string `json:"email" binding:"required,email"`
}
