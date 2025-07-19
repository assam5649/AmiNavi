package post

type Works struct {
	Title    string `json:"title"`
	WorkURL  string `json:"work_url"`
	Count    int    `json:"count"`
	Bookmark bool   `json:"bookmark"`
}
