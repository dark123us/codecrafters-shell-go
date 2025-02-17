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
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	command := strings.TrimSpace(str)

	fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
}
