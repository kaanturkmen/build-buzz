package request

type NotifyRequest struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
