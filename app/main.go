package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	builtins := map[string]bool{"exit": true, "echo": true, "type": true}

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(0)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			handleType(args, builtins)
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

// handleType manages the logic for identifying command locations
func handleType(args []string, builtins map[string]bool) {
	if len(args) == 0 {
		return
	}
	target := args[0]

	// 1. Check if it's a shell builtin
	if builtins[target] {
		fmt.Printf("%s is a shell builtin\n", target)
		return
	}

	// 2. Search for the executable in the PATH
	fullPath, found := findInPath(target)
	if found {
		fmt.Printf("%s is %s\n", target, fullPath)
	} else {
		fmt.Printf("%s: not found\n", target)
	}
}

// findInPath iterates through directories in the PATH environment variable
func findInPath(command string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	// filepath.SplitList automatically handles ':' on Unix and ';' on Windows
	paths := filepath.SplitList(pathEnv)

	for _, dir := range paths {
		fullPath := filepath.Join(dir, command)
		// Check if the file exists
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, true
		}
	}
	return "", false
}