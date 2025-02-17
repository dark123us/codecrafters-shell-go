package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	command := scanner.Text()

	fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
}
