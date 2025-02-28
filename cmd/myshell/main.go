package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/parseargs"
	"github.com/codecrafters-io/shell-starter-go/internal/readinput"
	"github.com/codecrafters-io/shell-starter-go/internal/redirect"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	autocCompleteFunc := func(text string) string {
		// fmt.Println("autocCompleteFunc text", text)
		if text == "ech" {
			return "echo "
		} else if text == "exi" {
			return "exit "
		}
		return text
	}

	reader := readinput.NewReader()
	defer reader.Close()

	for {
		// fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		// str, err := bufio.NewReader(os.Stdin).ReadString('\n')

		str, err := reader.ReadLine(autocCompleteFunc)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			break
		}
		args := parseargs.TrimString(str)
		result, err := command.RunCommand(args.Command, args.Args)

		if args.IsRedirectError {
			output := result.ErrorOutput
			if args.IsAppend {
				redirect.RedirectFileAppend(args.RedirectFile, output)
			} else {
				redirect.RedirectFile(args.RedirectFile, output)
			}
		} else if err != nil {
			fmt.Fprint(os.Stdout, string(result.ErrorOutput))
		}

		if args.IsRedirect {
			if args.IsAppend {
				redirect.RedirectFileAppend(args.RedirectFile, result.Output)
			} else {
				redirect.RedirectFile(args.RedirectFile, result.Output)
			}
		} else {
			fmt.Fprint(os.Stdout, string(result.Output))
		}
	}
}
