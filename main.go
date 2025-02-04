package main

import (
    "fmt"
    "log"
    "os"

    "github.com/charmbracelet/huh"
)

func main() {
    var (
        jiraTicket    string
        sourceBranch  string
        targetBranch  string
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
}

