package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
)

func initialChecks() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: can't get current directory. %v", err)
	}

	if !isGitRepository() {
		log.Fatalf("Current directory: %s is not a git repository", currentDir)
	}
}

func main() {
	initialChecks()

	inputs, err := InitialForm()
	if err != nil {
		log.Fatal(err)
	}

	// get the patch between the input branches
	gitPatch, err := getGitPatch(inputs.TargetBranch, inputs.SourceBranch)
	if err != nil {
		log.Fatal(err)
	}

	config := GetConfig()

	jiraIssueLink := GetJiraLink(config.JiraDomain, inputs.JiraTicket)
	jiraDescription, err := GetJiraDescription(
		config.JiraEmail,
		config.JiraApiKey,
		config.JiraDomain,
		inputs.JiraTicket,
	)

	if err != nil {
		log.Fatal(err)
	}

	// set up the data for the template
	data := PromptInputs{
		TargetBranch:    inputs.TargetBranch,
		SourceBranch:    inputs.SourceBranch,
		GitPatch:        gitPatch,
		JiraTicket:      jiraIssueLink,
		JiraDescription: jiraDescription,
	}

	prompt, err := CreatePrompt(data)
	if err != nil {
		log.Fatal(err)
	}

	response := CallGemini(prompt, "gemini-2.0-flash", config.GeminiApiKey)
	editedResponse, err := EditResponseForm(response)
	if err != nil {
		log.Fatal(err)
	}

	pr, err := GetPrId()
	// if this errors, then the pr doesn't exist
	if err != nil {
		// confirm with the user that they want to create it
		shouldCreate, err := CreatePrForm(inputs.TargetBranch, inputs.SourceBranch)
		if err != nil {
			log.Fatal(err)
		}
		if shouldCreate {
			CreatePr(editedResponse.Title, editedResponse.Body, inputs.TargetBranch)
		}
	} else {
		// confirm with the user that they want to edit it
		shouldEdit, err := EditPrForm(inputs.TargetBranch, inputs.SourceBranch)
		if err != nil {
			log.Fatal(err)
		}
		if shouldEdit {
			EditPr(editedResponse.Title, editedResponse.Body, *pr)
		}
	}

	// finally, get the PR URL and show it to the user
	prUrl, err := GetPrUrl()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PR URL: %s", prUrl)
}
