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
		result := command.RunCommand(args.Command, args.Args)
		if args.IsRedirect {
			redirect.RedirectFile(args.RedirectFile, result)
		}
	}
}
