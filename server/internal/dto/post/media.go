package post

type MediaRequest struct {
	Image []byte `json:"image"`
}

type MediaResponse struct {
	Csv    []byte `json:"csv"`
	CsvUrl string `json:"csv_url"`
}
