package main

import (
	"fmt"
	"strings"
	"os"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

var cliCommands = map[string]cliCommand{}
func init() {
	cliCommands = map[string]cliCommand{
		"exit": {
			name:		 "exit",
			description: "Exit the Pokedex",
			callback:	 commandExit,
		},
		"help": {
			name:		 "help",
			description: "Show help menu",
			callback:	 commandHelp,
		},
	}
}

func cleanInput(text string) []string {
	var stringSlice []string
	textSlice := strings.Fields(text)

	for _, str := range textSlice {
		str = strings.ToLower(str)
		stringSlice = append(stringSlice, str)
	}

	return stringSlice
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for name, cmd := range cliCommands {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}