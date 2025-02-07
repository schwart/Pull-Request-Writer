package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"github.com/cli/go-gh/v2"
)

type PullRequest struct {
	Id int `json:"number"`
	State string `json:"state"`
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
	log.Println(response.String())
	return &id, nil
}

func EditPr(title string, body string, pr PullRequest) {

	// check if the PR is closed, if it is, re-open it
	if pr.State == "CLOSED" {
		_, _, err := gh.Exec("pr", "reopen", strconv.Itoa(pr.Id))
		if err != nil {
			log.Printf("Failed to reopen PR: %d", pr.Id)
			log.Fatal(err)
		}
	}

	response, _, err := gh.Exec("pr", "edit", strconv.Itoa(pr.Id), "--body", body, "--title", title)
	if err != nil {
		log.Printf("Failed to edit PR: %d", pr.Id)
		log.Fatal(err)
	}
	fmt.Println(response.String())
}

func CreatePr(title string, body string) {
	_, _, err := gh.Exec("pr", "create", "--body", body, "--title", title)
	if err != nil {
		log.Println("Failed to create PR.")
		log.Fatal(err)
	}
}

func EditOrUpdatePr(title string, body string) {
	// if we have an pullRequest, then we need to edit the current PR
	// if we don't, then we need to update it
	log.Println("Getting PR info")
	pullRequest, err := GetPrId()
	// using an error to indicate that there's no PR for this branch
	// it might indicate an error in unmarshalling the JSON but I think that's unlikely hehe
	// if there's an error, it's likely there's no PR
	if err != nil {
		log.Println("Creating PR as non-existed prior.")
		CreatePr(title, body)
		return
	}
	log.Println("Editing PR with new title and body.")
	EditPr(title, body, *pullRequest)
}
