package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

type CommandType int

const (
	ExitCommand CommandType = iota
	EchoCommand
	TypeCommand
	PwdCommand
	CdCommand
)

var commandNames = map[string]CommandType{
	"exit": ExitCommand,
	"echo": EchoCommand,
	"type": TypeCommand,
	"pwd":  PwdCommand,
	"cd":   CdCommand,
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

func handleTypeCommand(arg string) {
	_, ok := commandNames[arg]
	if ok {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", arg)
		return
	}

	paths := GetPathDirs()
	for _, path := range paths {
		fullPath := filepath.Join(path, arg)
		if _, err := os.Stat(fullPath); err == nil {
			fmt.Fprintf(os.Stdout, "%s is %s\n", arg, fullPath)
			return
		}
	}
	fmt.Fprintf(os.Stdout, "%s: not found\n", arg)
}

type State int

const (
	StateNormal = iota
	StateSingleQuote
	StateDoubleQuote
)

func trimString(argin string) []string {
	var result []string
	var state State = StateNormal
	var cur int = 0
	quote := strings.ReplaceAll(argin, "''", "")
	arg := strings.ReplaceAll(quote, "\"\"", "")
	for i, c := range arg {
		switch state {
		case StateNormal:
			if c == '\n' || c == ' ' {
				if cur < i {
					result = append(result, arg[cur:i])
				}
				cur = i + 1
			} else if c == '\'' {
				if cur < i {
					result = append(result, arg[cur:i])
				}
				state = StateSingleQuote
				cur = i + 1
			} else if c == '"' {
				if cur < i {
					result = append(result, arg[cur:i])
				}
				state = StateDoubleQuote
				cur = i + 1
			}
		case StateSingleQuote:
			if c == '\'' {
				if cur < i {
					result = append(result, arg[cur:i])
				}
				state = StateNormal
				cur = i + 1
			}
		case StateDoubleQuote:
			if c == '"' {
				if cur < i {
					result = append(result, arg[cur:i])
				}
				state = StateNormal
				cur = i + 1
			}
		}

	}
	if cur < len(arg)-1 {
		result = append(result, arg[cur:])
	}
	return result
}

func handleRunApp(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Output()

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		return
	}
}

func runCommand(command string, args []string) {
	if isApp(command) {
		handleRunApp(command, args)
		return
	}
	commandType, ok := commandNames[command]
	if !ok {
		fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
		return
	}
	switch commandType {
	case ExitCommand:
		n := 0
		if len(args) > 0 {
			n, _ = strconv.Atoi(args[0])
		}
		os.Exit(n)
	case EchoCommand:
		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(args, " "))
	case TypeCommand:
		if len(args) == 0 {
			fmt.Fprintf(os.Stdout, "not found args\n")
			return
		}
		handleTypeCommand(args[0])
	case PwdCommand:
		curDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", curDir)
	case CdCommand:
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "cd: missing argument\n")
			return
		}
		if args[0] == "~" {
			args[0] = os.Getenv("HOME")
		}
		err := os.Chdir(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", args[0])
			return
		}
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		str, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			break
		}
		args := trimString(str)
		runCommand(args[0], args[1:])
	}
}
