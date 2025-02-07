package main

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenInVim(text string) {
	// take the text input, pipe it into vim
	cmd := exec.Command("vim")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fmt.Println(err)
}
