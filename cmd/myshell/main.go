package main

import (
	"bufio"
	"fmt"
	"os"
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
)

var commandNames = map[string]CommandType{
	"exit": ExitCommand,
	"echo": EchoCommand,
	"type": TypeCommand,
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

func handleRunApp(command string, args []string) {
	countArgs := len(args)
	fmt.Fprintf(os.Stdout, "Program was passed %d args (including programm name).\n", countArgs+1)
	fmt.Fprintf(os.Stdout, "Arg #0 (programm name): %s\n", command)
	for i, arg := range args {
		fmt.Fprintf(os.Stdout, "Arg #%d: %s\n", i+1, arg)
	}
	sign := 5998595441
	fmt.Fprintf(os.Stdout, "Programm Signature: %d\n", sign)
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
		args := strings.Split(strings.TrimSpace(str), " ")
		runCommand(args[0], args[1:])
	}
}
