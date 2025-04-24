package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
)

func checkDirectoryExists(directory string) bool {
	if stat, err := os.Stat(directory); err == nil && stat.IsDir() {
		return true
	}
	return false
}

func getWorkingDirectory(otherDirectory string) (*string, error) {
	if otherDirectory != "" {
		if checkDirectoryExists(otherDirectory) {
			return &otherDirectory, nil
		}
		return nil, fmt.Errorf("error: directory %s does not exist", otherDirectory)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error: can't get current directory. %v", err)
	}

	if !isGitRepository() {
		log.Fatalf("Current directory: %s is not a git repository", currentDir)
	}
	return &currentDir, nil
}

func main() {
	debugDirectory := flag.String("directory", "", "Directory to run the script in (if not the current one).")

	flag.Parse()

	workingDirectory, err := getWorkingDirectory(*debugDirectory)
	if err != nil {
		log.Fatal(err)
	}

	inputs, err := InitialForm(*workingDirectory)
	if err != nil {
		log.Fatal(err)
	}

	// get the patch between the input branches
	gitPatch, err := getGitPatch(inputs.TargetBranch, inputs.SourceBranch, *workingDirectory)
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
