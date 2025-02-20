package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type AppResult struct {
	Output []byte
	Error  error
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

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

func handleRunApp(command string, args []string) (AppResult, error) {
	result := AppResult{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	cmd := exec.Command(command, args...)
	cmd.Stdin = result.Stdin
	cmd.Stdout = result.Stdout
	cmd.Stderr = result.Stderr
	err := cmd.Run()
	if err != nil {
		result.Output = []byte(fmt.Sprintf("%v", result.Stderr))
		result.Error = err
		return result, err
	}
	result.Output = []byte(fmt.Sprintf("%v", result.Stdout))
	return result, nil
}
