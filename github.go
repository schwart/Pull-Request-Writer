package main

import (
	"encoding/json"
	"github.com/cli/go-gh/v2"
	"log"
	"strconv"
)

type PullRequest struct {
	Id    int    `json:"number"`
	State string `json:"state"`
}

type PullRequestUrl struct {
	Url string `json:"url"`
}

func GetPrId() (*PullRequest, error) {
	response, _, err := gh.Exec("pr", "view", "--json", "number,state")
	if err != nil {
		// will error if there's no PRs for the current branch
		return nil, err
	}

	id := PullRequest{}
	if err := json.Unmarshal(response.Bytes(), &id); err != nil {
		return nil, err
	}
	return &id, nil
}

func GetPrUrl() (string, error) {
	response, _, err := gh.Exec("pr", "view", "--json", "url")
	if err != nil {
		// will error if there's no PRs for the current branch
		return "", err
	}

	url := PullRequestUrl{}
	if err := json.Unmarshal(response.Bytes(), &url); err != nil {
		return "", err
	}
	return url.Url, nil
}

func EditPr(title string, body string, pr PullRequest) {
	// check if the PR is closed, if it is, re-open it
	if pr.State == "CLOSED" {
		_, _, err := gh.Exec("pr", "reopen", strconv.Itoa(pr.Id))
		if err != nil {
			log.Fatalf("Failed to reopen PR: %d %v", pr.Id, err)
		}
	}

	_, _, err := gh.Exec("pr", "edit", strconv.Itoa(pr.Id), "--body", body, "--title", title)
	if err != nil {
		log.Fatalf("Failed to edit PR: %d, %v", pr.Id, err)
	}
}

func CreatePr(title string, body string, targetBranch string) {
	_, _, err := gh.Exec("pr", "create", "--body", body, "--title", title, "--base", targetBranch)
	if err != nil {
		log.Fatalf("Failed to create PR for: %s, %v", targetBranch, err)
	}
}
