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
			os.Exit(0)
		}

		// 1. Trim and Split the input into parts
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		
		// parts[0] is the command, parts[1:] are the arguments
		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		// 2. Command Router (Switch Statement)
		switch command {
		case "exit":
			os.Exit(0)

		case "echo":
			// Join the arguments back together with a single space
			fmt.Println(strings.Join(args, " "))

		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}