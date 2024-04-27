package response

type NotifyResponse struct {
	GithubUsername string `json:"github_username"`
	SlackId        string `json:"slack_id"`
	SlackEmail     string `json:"slack_email"`
}
