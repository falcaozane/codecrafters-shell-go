package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 1. Initialize the reader outside the loop for efficiency
	reader := bufio.NewReader(os.Stdin)

	for {
		// PRINT: Display the prompt
		fmt.Print("$ ")

		// READ: Wait for user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue // Don't exit, just wait for the next attempt
		}

		// Clean the input (remove newline characters)
		command := strings.TrimSpace(input)

		// EXIT CONDITION: Allow the user to quit the REPL
		if command == "exit" {
			break
		}

		// EVAL & PRINT: Logic to handle the command
		if command != "" {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}