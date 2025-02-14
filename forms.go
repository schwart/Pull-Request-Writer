package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
)

type InitialFormOutputs struct {
	JiraTicket   string
	SourceBranch string
	TargetBranch string
}

func InitialForm() (*InitialFormOutputs, error) {
	var (
		jiraTicket   string
		sourceBranch string
		targetBranch string
	)

	// set up variables for the form
	gitBranchSuggestions, _ := getGitBranches()
	sourceBranch = gitBranchSuggestions[0]
	targetBranch = getDefaultTargetBranch(gitBranchSuggestions)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's the JIRA ticket number?").
				Value(&jiraTicket).
				Placeholder("AB-123"),
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
		return nil, err
	}

	return &InitialFormOutputs{
		JiraTicket:   jiraTicket,
		SourceBranch: sourceBranch,
		TargetBranch: targetBranch,
	}, nil
}

func EditResponseForm(response *PullRequestResponse) (*PullRequestResponse, error) {

	var editedResponse PullRequestResponse
	editedResponse.Body = response.Body
	editedResponse.Title = response.Title

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("PR Title").
				Value(&editedResponse.Title),
			huh.NewText().
				Title("PR Body").
				Value(&editedResponse.Body),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}
	return &editedResponse, nil
}

func CreatePrForm(sourceBranch string, targetBranch string) (bool, error) {
	targetToSource := fmt.Sprintf("%s -> %s", targetBranch, sourceBranch)
	var inputString string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("No existing PR for %s, create one?", targetToSource)).
				Options(
					huh.NewOption("Create", "Create"),
					huh.NewOption("Cancel", "Cancel"),
				).
				Value(&inputString),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return inputString == "Create", nil
}

func EditPrForm(sourceBranch string, targetBranch string) (bool, error) {
	targetToSource := fmt.Sprintf("%s -> %s", targetBranch, sourceBranch)
	var inputString string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("PR already exists for %s, edit it?", targetToSource)).
				Options(
					huh.NewOption("Edit", "Edit"),
					huh.NewOption("Cancel", "Cancel"),
				).
				Value(&inputString),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return inputString == "Edit", nil
}
