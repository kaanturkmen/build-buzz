package github

import (
	"encoding/json"
	"fmt"
	"github.com/kaanturkmen/build-buzz/internal/config"
	"net/http"
	"strings"
)

const (
	githubAPIBaseURL = "https://api.github.com"
)

type GithubClient struct {
	Organization string
	Token        string
	HttpClient   *http.Client
}

func NewGitHubClient(organization string, token string) *GithubClient {
	return &GithubClient{
		Organization: organization,
		Token:        token,
		HttpClient:   &http.Client{},
	}
}

func (c *GithubClient) FetchRepositories() ([]Repository, error) {
	url := fmt.Sprintf("%s/orgs/%s/repos", githubAPIBaseURL, c.Organization)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.Token)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []Repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	return repos, err
}

func (c *GithubClient) FetchUserCommitEmail(username string, overrides config.EmailOverrides) (string, error) {
	repos, err := c.FetchRepositories()
	if err != nil {
		return "", err
	}

	for _, repo := range repos {
		url := fmt.Sprintf("%s/repos/%s/%s/commits?author=%s", githubAPIBaseURL, c.Organization, repo.Name, username)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Add("Authorization", "Bearer "+c.Token)

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var commits []Commit
		err = json.NewDecoder(resp.Body).Decode(&commits)
		if err != nil {
			continue
		}

		for _, commit := range commits {
			email := commit.Commit.Author.Email

			if strings.Contains(email, "@users.noreply.github.com") {
				continue
			}

			if overrideEmail, ok := overrides[email]; ok {
				return overrideEmail, nil
			}
			if email != "" {
				return email, nil
			}
		}
	}

	return "", fmt.Errorf("email address not found for user: %s", username)
}
