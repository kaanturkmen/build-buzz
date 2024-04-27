package github

type Repository struct {
	Name string `json:"name"`
}

type Commit struct {
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	Commit struct {
		Author struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"author"`
	} `json:"commit"`
}
