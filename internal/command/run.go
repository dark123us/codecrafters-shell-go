package command

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

type CommandResult struct {
	Output      []byte
	ErrorOutput []byte
	Error       error
}

func getCommandResult(output []byte, errorOutput []byte, err error) CommandResult {
	return CommandResult{
		Output:      output,
		ErrorOutput: errorOutput,
		Error:       err,
	}
}

func RunCommand(command string, args []string) (CommandResult, error) {
	var result CommandResult
	if isApp(command) {
		result, err := handleRunApp(command, args)
		if err != nil {
			return getCommandResult(result.Output, result.ErrorOutput, err), nil
		}
		return getCommandResult(result.Output, result.ErrorOutput, nil), nil
	}
	commandType, ok := commandNames[command]
	if !ok {
		fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
		return result, nil
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
			return result, nil
		}
		handleTypeCommand(args[0])
	case PwdCommand:
		curDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			return result, nil
		}
		fmt.Fprintf(os.Stdout, "%s\n", curDir)
	case CdCommand:
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "cd: missing argument\n")
			return result, nil
		}
		if args[0] == "~" {
			args[0] = os.Getenv("HOME")
		}
		err := os.Chdir(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", args[0])
			return result, nil
		}
	}
	return result, nil
}
