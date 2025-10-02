package post

type MediaRequest struct {
	Image []byte `json:"image"`
}

type MediaResponse struct {
	Csv      []byte `json:"csv"`
	FileName string `json:"file_name"`
}
