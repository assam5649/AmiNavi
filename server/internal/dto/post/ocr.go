package post

type OCRRequest struct {
	Image []byte `json:"image"`
}

type OCRResponse struct {
	Pattern string `json:"pattern"`
}
