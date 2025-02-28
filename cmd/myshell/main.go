package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/parseargs"
	"github.com/codecrafters-io/shell-starter-go/internal/readinput"
	"github.com/codecrafters-io/shell-starter-go/internal/redirect"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func compareAppPrefixes(old []string, new []string) bool {
	if len(old) != len(new) {
		return false
	}
	for i, v := range old {
		if v != new[i] {
			return false
		}
	}
	return true
}
func main() {

	appPrefixes := []string{}

	autocCompleteFunc := func(text string) (string, error) {
		if text == "ech" {
			return "echo ", nil
		} else if text == "exi" {
			return "exit ", nil
		} else if app, err := command.FindAppPrefix(text); err == nil {
			if len(app) == 1 {
				return app[0] + " ", nil
			}
			if compareAppPrefixes(appPrefixes, app) {
				// fmt.Fprint(os.Stdout, "\r\n", strings.Join(appPrefixes, "  "), "\r\n$ ", text)
				fmt.Printf("\r\n%s\r\n$ %s", strings.Join(appPrefixes, "  "), text)
				os.Stdout.Sync()
				return text, nil
			}
			appPrefixes = app
			return text, errors.New("found multiple")
		}
		return text, errors.New("not found")
	}

	reader := readinput.NewReader()
	defer reader.Close()

	for {
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
