package main

import (
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
