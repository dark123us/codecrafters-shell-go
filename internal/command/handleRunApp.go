package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type AppResult struct {
	Output      []byte
	ErrorOutput []byte
	Error       error
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
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

func FindAppPrefix(prefix string) ([]string, error) {
	paths := GetPathDirs()
	matches := []string{}
	for _, path := range paths {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasPrefix(file.Name(), prefix) {
				matches = append(matches, file.Name())
			}
		}
	}
	if len(matches) == 0 {
		return nil, errors.New("not found")
	}
	return matches, nil
}

func handleRunApp(command string, args []string) (AppResult, error) {
	result := AppResult{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	result.Output = stdout
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ErrorOutput = exitErr.Stderr
			// fmt.Print(string(result.ErrorOutput))
		} else {
			result.ErrorOutput = []byte(fmt.Sprintf("%v", err))
			// fmt.Print(string(result.ErrorOutput))
		}
		result.Error = err
		return result, err
	}
	return result, nil
}
