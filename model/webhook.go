package model

//GitlabRepository represents repository information from the webhook
type Repository struct {
	Name      string
	GitSshUrl string `json:"git_ssh_url"`
}

//Author represents author information from the webhook
type Author struct {
	Name, Email string
}

//Commit represents commit information from the webhook
type Commit struct {
	Id, Message, Timestamp, Url string
	Author                      Author
}

//Webhook represents push information from the webhook
type Webhook struct {
	ObjectKind        string `json:"object_kind"`
	EventName         string `json:"event_name"`
	Ref               string
	UserName          string
	Repository        Repository
	Commits           []Commit
	TotalCommitsCount int
}
