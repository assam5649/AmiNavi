package put

type WorksId struct {
	Title	string	`json:"title"`,
	WorkURL	string	`json:"work_url"`,
	Count	int	`json:"count"`,
	Bookmark	string	`json:"bookmark"`
}