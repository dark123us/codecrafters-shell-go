package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		args := strings.Split(strings.TrimSpace(str), " ")

		command := args[0]

		if command == "exit" {
			break
		}

		fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
	}
}
