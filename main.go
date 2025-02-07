package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/charmbracelet/huh"
)

// embed the prompt file
//go:embed prompt-template.txt
var promptTemplate string


func getGitPatch(targetBranch, sourceBranch string) (string, error) {
	targetToSource := fmt.Sprintf("%s..%s", targetBranch, sourceBranch)
	fmt.Println(targetToSource)
	// eg: git log -p --full-diff master..RC-001-some-branch
	cmd := exec.Command("git", "log", "-p", "--full-diff", targetToSource)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error running git log: %v", err)
	}
	return string(output), nil
}

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	return err == nil
}

func getGitBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--sort=-committerdate")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error getting git branches: %v", err)
	}
	lines := strings.Split(string(output), "\n")

	var branches []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		trimmed = strings.TrimPrefix(trimmed, "* ")
		branches = append(branches, trimmed)
	}

	return branches, nil

} 

func getDefaultTargetBranch(branches []string) string {
	// try to find "master" first
	for _, branch := range branches {
		if branch == "master" {
			return "master"
		}
	}
	// then try to find "main"
	for _, branch := range branches {
		if branch == "main" {
			return "main"
		}
	}
	// if we can't find either, just return an empty string
	return ""
}

func main() {
	// check the current directory is a git repository
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if !isGitRepository() {
		log.Fatal("Current directory: %s is not a git repository", currentDir)
	}

	// set up variables for the form
	var (
		jiraTicket   string
		sourceBranch string
		targetBranch string
	)

	gitBranchSuggestions, _ := getGitBranches()
	sourceBranch = gitBranchSuggestions[0]
	targetBranch = getDefaultTargetBranch(gitBranchSuggestions)

	// create the form with some suggestions / default values 
	// jira ticket: no suggestion or default value cos I cba with the API
	// source branch:
	//	default: the most recent branch that was committed to
	//	suggestions: all branches
	// target branch:
	//	default: master, main or nothing. Whichever is found first
	//	suggestions: all branches
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's the JIRA ticket link?").
				Value(&jiraTicket).
				Placeholder("https://example.atlassian.net/browse/AB-123"),
			huh.NewInput().
				Title("What's the source branch?").
				Value(&sourceBranch).
				Suggestions(gitBranchSuggestions),
			huh.NewInput().
				Title("What's the target branch?").
				Value(&targetBranch).
				Suggestions(gitBranchSuggestions),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	// get the patch between the input branches
	gitPatch, err := getGitPatch(targetBranch, sourceBranch)
	if err != nil {
		log.Fatal(err)
	}

	// set up the data for the template
	data := map[string]string {
		"TargetBranch": targetBranch,
		"SourceBranch": sourceBranch,
		"GitPatch": gitPatch,
		"JiraTicket": jiraTicket,
	}

	t := template.Must(template.New("prompt-output").Parse(promptTemplate))

	// save the templated output into a buffer
	var outputBytes bytes.Buffer
	t.Execute(&outputBytes, data)

	// convert buffer to string
	outputString := outputBytes.String()
	if outputString == "<nil>" {
		log.Fatal("Cannot template prompt, you're on your own.")
	}

	fmt.Println("Completed templating prompt, calling Gemini")

	response := CallGemini(outputString, "gemini-2.0-flash")
	// save response to clipboard
	editedTitle, err := EditResponseInVim(response.Title)
	if err != nil {
		log.Fatal(err)
	}
	editedBody, err := EditResponseInVim(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// clipboard.Write(clipboard.FmtText, []byte(editedTitle))
	// fmt.Println("Saved response to clipboard")
	EditOrUpdatePr(editedTitle, editedBody, targetBranch)
}
