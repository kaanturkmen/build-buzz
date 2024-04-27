package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackClient struct {
	httpClient *http.Client
	token      string
}

func NewSlackClient(token string) *SlackClient {
	return &SlackClient{
		httpClient: &http.Client{},
		token:      token,
	}
}

func (c *SlackClient) UserLookupByEmail(email string) (string, error) {
	url := fmt.Sprintf("https://slack.com/api/users.lookupByEmail?email=%s", email)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData struct {
		OK   bool `json:"ok"`
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		Error string `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", err
	}

	if !respData.OK {
		return "", fmt.Errorf("slack API error: %s", respData.Error)
	}

	return respData.User.ID, nil
}

func (c *SlackClient) SendMessage(userID, message string) error {
	url := "https://slack.com/api/chat.postMessage"

	payload := map[string]interface{}{
		"channel": userID,
		"text":    message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respData struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&respData)

	if err != nil {
		return err
	}

	if !respData.OK {
		return fmt.Errorf("slack API error: %s", respData.Error)
	}

	return nil
}
