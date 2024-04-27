package main

import (
	"github.com/kaanturkmen/build-buzz/internal/config"
	"github.com/kaanturkmen/build-buzz/internal/handler"
	"github.com/kaanturkmen/build-buzz/internal/service"
	"github.com/kaanturkmen/build-buzz/pkg/github"
	"github.com/kaanturkmen/build-buzz/pkg/slack"
	"log"
	"net/http"
	"os"
)

func main() {
	overrides, err := config.LoadEmailOverrides("./config/email_overrides.json")

	if err != nil {
		log.Printf("error loading email overrides: %v\n", err)
		return
	}

	githubToken := os.Getenv("GITHUB_TOKEN")

	if githubToken == "" {
		log.Printf("please set a github token to the GITHUB_TOKEN env variable\n")
		return
	}

	slackToken := os.Getenv("SLACK_TOKEN")

	if slackToken == "" {
		log.Printf("please set a slack token to the SLACK_TOKEN env variable\n")
		return
	}

	githubOrganizationName := os.Getenv("GITHUB_ORG_NAME")

	if githubOrganizationName == "" {
		log.Printf("please set GITHUB_ORG_NAME env variable\n")
		return
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	githubClient := github.NewGitHubClient(githubOrganizationName, githubToken)
	slackClient := slack.NewSlackClient(slackToken)

	notifyService := service.NewNotifyService(overrides, githubClient, slackClient, githubToken, slackToken, githubOrganizationName)
	notifyHandler := handler.NewNotifyHandler(notifyService)

	http.HandleFunc("/notify", notifyHandler.Notify)

	log.Printf("server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
