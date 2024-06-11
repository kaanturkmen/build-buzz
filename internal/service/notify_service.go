package service

import (
	"fmt"
	"github.com/kaanturkmen/build-buzz/internal/config"
	"github.com/kaanturkmen/build-buzz/internal/request"
	"github.com/kaanturkmen/build-buzz/internal/response"
	"github.com/kaanturkmen/build-buzz/pkg/github"
	"github.com/kaanturkmen/build-buzz/pkg/slack"
	"strings"
)

type NotifyService struct {
	overrides              config.EmailOverrides
	githubClient           *github.GithubClient
	slackClient            *slack.SlackClient
	githubToken            string
	slackToken             string
	githubOrganizationName string
}

func NewNotifyService(overrides config.EmailOverrides, githubClient *github.GithubClient, slackClient *slack.SlackClient,
	githubToken string, slackToken string, githubOrganizationName string) *NotifyService {

	return &NotifyService{overrides: overrides, githubClient: githubClient, slackClient: slackClient,
		githubToken: githubToken, slackToken: slackToken, githubOrganizationName: githubOrganizationName}
}

func (s *NotifyService) Notify(notifyRequest request.NotifyRequest) (*response.NotifyResponse, error) {
	if strings.Contains(notifyRequest.Username, " ") {
		return nil, fmt.Errorf("username has space in it")
	}

	email, err := s.githubClient.FetchUserCommitEmail(notifyRequest.Username, s.overrides)

	if err != nil {
		return nil, fmt.Errorf("fetch user commit email had errors: %v", err)
	}

	userID, err := s.slackClient.UserLookupByEmail(email)

	if err != nil {
		return nil, fmt.Errorf("user lookup by email had errors: %v", err)
	}

	if err := s.slackClient.SendMessage(userID, notifyRequest.Message); err != nil {
		return nil, fmt.Errorf("send message had errors: %v", err)
	}

	return &response.NotifyResponse{GithubUsername: notifyRequest.Username, SlackId: userID, SlackEmail: email}, nil
}
