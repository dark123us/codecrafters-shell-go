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

func handleTypeCommand(arg string) {
	// commandType, ok := commandNames[arg]
	// if !ok {
	// 	fmt.Fprintf(os.Stdout, "%s: not found\n", arg)
	// 	return
	// }
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

func runCommand(command string, args []string) {
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
