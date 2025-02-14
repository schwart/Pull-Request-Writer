package main

import (
	"fmt"
	"os/exec"
	"strings"
)

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
