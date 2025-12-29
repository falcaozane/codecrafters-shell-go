package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		input, err := reader.ReadString('\n')
		if err != nil {
			// Handle EOF (Ctrl+D) gracefully
			os.Exit(0)
		}

		// Clean the input to remove \n or \r\n
		command := strings.TrimSpace(input)

		// Check for known commands
		if command == "exit" {
			os.Exit(0)
		} else if command == "" {
			continue // Handle empty enter key press
		} else {
			// Print error for unknown commands
			fmt.Printf("%s: command not found\n", command)
		}
	}
}