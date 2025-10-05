package post

type FixRequest struct {
	Csv []byte `json:"csv"`
}

type FixResponse struct {
	Csv      []byte `json:"csv"`
	FileName string `json:"file_name"`
}
