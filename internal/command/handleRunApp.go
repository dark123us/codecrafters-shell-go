package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetPathDirs() []string {
	path := os.Getenv("PATH")
	return strings.Split(path, ":")
}

func isApp(name string) bool {
	paths := GetPathDirs()
	for _, path := range paths {
		fullPath := filepath.Join(path, name)
		if _, err := os.Stat(fullPath); err == nil {
			return true
		}
	}
	return false
}

func handleRunApp(command string, args []string) []byte {
	var result []byte
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		return result
	}
	return result
}
