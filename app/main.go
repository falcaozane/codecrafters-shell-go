package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
			if len(args) > 0 {
				handleType(args[0], builtins)
			}
		default:
			// 1. Check if the command exists in PATH
			fullPath, found := findInPath(command)
			
			if found {
				// Use the original command name for the arguments list
				// but point the execution Path to the fullPath found.
				cmd := exec.Command(command, args...)
				cmd.Path = fullPath 
				
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				
				err := cmd.Run()
				if err != nil {
					fmt.Printf("%s: %v\n", command, err)
				}
			} else {
				fmt.Printf("%s: command not found\n", command)
			}
		}
	}
}

// handleType remains the same as before...
func handleType(target string, builtins map[string]bool) {
	if builtins[target] {
		fmt.Printf("%s is a shell builtin\n", target)
		return
	}
	if path, found := findInPath(target); found {
		fmt.Printf("%s is %s\n", target, path)
	} else {
		fmt.Printf("%s: not found\n", target)
	}
}

// findInPath remains the same, ensuring we only find executables
func findInPath(command string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)
	for _, dir := range paths {
		fullPath := filepath.Join(dir, command)
		info, err := os.Stat(fullPath)
		if err == nil && info.Mode().IsRegular() && info.Mode()&0111 != 0 {
			return fullPath, true
		}
	}
	return "", false
}