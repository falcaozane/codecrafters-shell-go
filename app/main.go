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
			if len(args) > 0 {
				handleType(args[0], builtins)
			}
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func handleType(target string, builtins map[string]bool) {
	// 1. Check Builtins first
	if builtins[target] {
		fmt.Printf("%s is a shell builtin\n", target)
		return
	}

	// 2. Search PATH for Executable
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)

	for _, dir := range paths {
		fullPath := filepath.Join(dir, target)
		
		info, err := os.Stat(fullPath)
		if err != nil {
			continue // Skip directories that don't exist
		}

		// Check if it's a regular file AND has execute permissions
		// 0111 is the bitmask for --x--x--x
		if info.Mode().IsRegular() && info.Mode()&0111 != 0 {
			fmt.Printf("%s is %s\n", target, fullPath)
			return
		}
	}

	// 3. Not found anywhere
	fmt.Printf("%s: not found\n", target)
}