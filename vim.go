package main

import (
	"log"
	"os"
	"os/exec"
	"fmt"
)

func EditResponseInVim(text string) (string, error) {
	// Create a temporary file. The pattern "vimedit*.txt" will have a random suffix.
	tmpFile, err := os.CreateTemp("", "vimedit*.txt")
	if err != nil {
		log.Fatalf("failed to create temp file: %v", err)
	}

	// store the file name for later
	fileName := tmpFile.Name()

	// Write the text to the temporary file.
	if _, err := tmpFile.Write([]byte(text)); err != nil {
		log.Fatalf("failed to write to temp file: %v", err)
	}

	// It's important to close the file so that Vim can open it.
	if err := tmpFile.Close(); err != nil {
		log.Fatalf("failed to close temp file: %v", err)
	}

	// Open Vim with the temporary file.
	cmd := exec.Command("vim", tmpFile.Name())
	// Make sure Vim is connected to your terminal.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run Vim.
	if err := cmd.Run(); err != nil {
		log.Fatalf("vim exited with error: %v", err)
	}

	// we want to read the file at the end, saving it's content to a variable
	// After Vim is closed, read the modified file content.
	modifiedBytes, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read file after editing: %w", err)
	}

	// Optionally, remove the temporary file after editing.
	// (You might want to preserve it if you need to retrieve the changes.)
	if err := os.Remove(tmpFile.Name()); err != nil {
		log.Printf("warning: failed to remove temp file: %v", err)
	}

	return string(modifiedBytes), nil
}
