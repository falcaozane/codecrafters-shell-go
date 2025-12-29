package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// List of builtins to check against
	builtins := map[string]bool{
		"exit": true,
		"echo": true,
		"type": true,
	}

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
			if len(args) == 0 {
				continue
			}
			target := args[0]

			// 1. Check if it's a builtin
			if builtins[target] {
				fmt.Printf("%s is a shell builtin\n", target)
			} else {
				// 2. Check if it's in the PATH
				path, found := getPath(target)
				if found {
					fmt.Printf("%s is %s\n", target, path)
				} else {
					fmt.Printf("%s: not found\n", target)
				}
			}

		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

// getPath searches for the executable in the system's PATH environment variable
func getPath(command string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)

	for _, dir := range paths {
		fullPath := filepath.Join(dir, command)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, true
		}
	}
	return "", false
}