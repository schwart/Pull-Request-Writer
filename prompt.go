package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"text/template"
)

//go:embed prompt-template.txt
var promptTemplate string

type PromptInputs struct {
	TargetBranch    string
	SourceBranch    string
	GitPatch        string
	JiraTicket      string
	JiraDescription string
}

func CreatePrompt(data PromptInputs) (string, error) {
	t := template.Must(template.New("prompt-output").Parse(promptTemplate))

	// save the templated output into a buffer
	var outputBytes bytes.Buffer
	err := t.Execute(&outputBytes, data)
	if err != nil {
		return "", fmt.Errorf("failed to create prompt: %v", err)
	}

	// convert buffer to string
	outputString := outputBytes.String()
	if outputString == "<nil>" {
		return "", errors.New("cannot template prompt, you're on your own")
	}

	return outputString, nil
}
