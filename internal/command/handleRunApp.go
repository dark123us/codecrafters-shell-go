package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
	matches := make(map[string]bool)
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
				matches[file.Name()] = true
			}
		}
	}
	if len(matches) == 0 {
		return nil, errors.New("not found")
	}
	file_names := []string{}
	for match := range matches {
		file_names = append(file_names, match)
	}
	slices.Sort(file_names)
	return file_names, nil
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
