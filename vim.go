package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// TODO: support $EDITOR variable
// TODO: sort out what we're gonna do on Windows
func EditResponseInVim(text string) (string, error) {
	// Create a temporary file. The pattern "vimedit*.txt" will have a random suffix.
	tmpFile, err := os.CreateTemp("", "vimedit*.txt")
	if err != nil {
		log.Fatalf("failed to create temp file: %v", err)
	}

	fileName := tmpFile.Name()

	if _, err := tmpFile.Write([]byte(text)); err != nil {
		log.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		log.Fatalf("failed to close temp file: %v", err)
	}

	cmd := exec.Command("vim", tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("vim exited with error: %v", err)
	}

	modifiedBytes, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read file after editing: %w", err)
	}

	if err := os.Remove(tmpFile.Name()); err != nil {
		log.Printf("warning: failed to remove temp file: %v", err)
	}

	return string(modifiedBytes), nil
}
