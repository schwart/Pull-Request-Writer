package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"github.com/charmbracelet/huh"
)

func runGitLog(targetBranch, sourceBranch string) (string, error) {
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

func main() {
	var (
		jiraTicket   string
		sourceBranch string
		targetBranch string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's the JIRA ticket link?").
				Value(&jiraTicket),
			huh.NewInput().
				Title("What's the source branch?").
				Value(&sourceBranch),
			huh.NewInput().
				Title("What's the target branch?").
				Value(&targetBranch),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nJIRA Ticket Link:", jiraTicket)
	fmt.Println("Source Branch:", sourceBranch)
	fmt.Println("Target Branch:", targetBranch)
	fmt.Println("Current Directory:", currentDir)

	if !isGitRepository() {
		log.Fatal("%s is not a git repository", currentDir)
	}
	diffOutput, err := runGitLog(targetBranch, sourceBranch)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(diffOutput)

}
