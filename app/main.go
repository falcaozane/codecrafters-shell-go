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
	builtins := map[string]bool{"exit": true, "echo": true, "type": true, "pwd": true, "cd": true}

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(0)
		}

		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")


		if input == "" {
			continue
		}

		parts := parseInput(input)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "pwd":
			presentWorkingDirectory()
		case "cd":
			if len(args) > 0 {
				changeDirectory(args[0])
			}
		case "type":
			if len(args) > 0 {
				handleType(args[0], builtins)
			}
		default:
			var stdoutFile string
			var cleanedArgs []string

			for i:=0;i<len(args);i++ {
				if args[i] == ">" || args[i] == "1>"{
					if i+1 < len(args) {
						stdoutFile = args[i+1]
						break
					}
				}
				cleanedArgs = append(cleanedArgs,args[i])
			}

			// 1. Check if the command exists in PATH
			fullPath, found := findInPath(command)
			
			if found {
				// Use the original command name for the arguments list
				// but point the execution Path to the fullPath found.
				cmd := exec.Command(command, cleanedArgs...)
				cmd.Path = fullPath 
				
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if stdoutFile != "" {
					file, err := os.OpenFile(stdoutFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
					if err != nil {
						fmt.Printf("Error opening file %s: %v\n", stdoutFile, err)
						continue
					}
					defer file.Close()
					cmd.Stdout = file
				}
				
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

func presentWorkingDirectory() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error retrieving current directory:", err)
		return
	}
	fmt.Println(dir)
}

func changeDirectory(path string) {
	var targetPath string

	// Handle Tilde Expansion
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("cd: could not determine home directory\n")
			return
		}
		targetPath = home
	} else if strings.HasPrefix(path, "~/") {
		// Handle paths like ~/Downloads or ~/some/other/path
		home, _ := os.UserHomeDir()
		targetPath = filepath.Join(home, path[2:])
	} else {
		targetPath = path
	}

	// Attempt to change the directory
	err := os.Chdir(targetPath)
	if err != nil {
		// Ensure we print the original 'path' in the error message, 
		// not the expanded internal 'targetPath'
		fmt.Printf("cd: %s: No such file or directory\n", path)
	}
}

// parseInput manages the State Machine for parsing command line input
func parseInput(input string) []string {
	var args []string
	var currentArg strings.Builder
	inSingle := false
	inDouble := false
	hasContent := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		// 1. Backslash OUTSIDE of any quotes
		if char == '\\' && !inSingle && !inDouble {
			if i+1 < len(input) {
				currentArg.WriteByte(input[i+1])
				i++
				hasContent = true
			}
			continue
		}

		// 2. Backslash INSIDE Double Quotes (The "Selective" rule)
		if char == '\\' && inDouble {
			if i+1 < len(input) {
				next := input[i+1]
				if next == '"' || next == '\\' || next == '$' || next == '\n' {
					currentArg.WriteByte(next)
					i++
				} else {
					currentArg.WriteByte('\\')
				}
			} else {
				currentArg.WriteByte('\\')
			}
			hasContent = true
			continue
		}

		// 3. Handle Quotes
		if char == '\'' && !inDouble {
			inSingle = !inSingle
			hasContent = true
			continue
		}
		if char == '"' && !inSingle {
			inDouble = !inDouble
			hasContent = true
			continue
		}

		// 4. Handle Delimiters
		if (char == ' ' || char == '\t') && !inSingle && !inDouble {
			if hasContent {
				args = append(args, currentArg.String())
				currentArg.Reset()
				hasContent = false
			}
			continue
		}

		currentArg.WriteByte(char)
		hasContent = true
	}

	if hasContent {
		args = append(args, currentArg.String())
	}
	return args
}