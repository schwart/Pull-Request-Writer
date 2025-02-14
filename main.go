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

	fmt.Println("Completed templating prompt, calling Gemini")

	response := CallGemini(prompt, "gemini-2.0-flash", config.GeminiApiKey)

	editedTitle, err := EditResponseInVim(response.Title)
	if err != nil {
		log.Fatal(err)
	}
	editedBody, err := EditResponseInVim(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	EditOrUpdatePr(editedTitle, editedBody, inputs.TargetBranch)
}
