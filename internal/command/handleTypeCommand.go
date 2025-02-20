package command

import (
	"fmt"
	"os"
	"path/filepath"
)

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
