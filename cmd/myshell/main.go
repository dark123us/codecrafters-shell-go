package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/parseargs"
	"github.com/codecrafters-io/shell-starter-go/internal/redirect"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		str, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			break
		}
		args := parseargs.TrimString(str)
		result, err := command.RunCommand(args.Command, args.Args)

		if err != nil {
			if args.IsRedirectError {
				//output := append(result.ErrorOutput, result.Output...)
				output := result.ErrorOutput
				redirect.RedirectFile(args.RedirectFile, output)
			} else {
				fmt.Fprint(os.Stdout, string(result.ErrorOutput))
			}
		}
		if args.IsRedirect {
			redirect.RedirectFile(args.RedirectFile, result.Output)
		} else {
			fmt.Fprint(os.Stdout, string(result.Output))
		}
	}
}
