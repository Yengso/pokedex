package main

import (
	"fmt"
	"bufio"
	"strings"
	"os"
	"github.com/yengso/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name		string
	description	string
	callback	func(*Config) error
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
		"map": {
			name:		 "map",
			description: "Show the start/next 20 locations",
			callback:	 commandMap,
		},
		"mapb": {
			name: 		 "mapb",
			description: "Show the previous 20 locations",
			callback:	 commandMapb,
		},
	}
}

type Config struct {
	Next		string
	Previous 	string
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

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for name, cmd := range cliCommands {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

func commandMap(cfg *Config) error {
	url := cfg.Next
	if url == "" {
		url = ""
	}

	page, err := pokeapi.PokedexAPI(url)
	if err != nil {
		return err
	}

	for _, r := range page.Results {
		fmt.Println(r.Name)
	}

	cfg.Next = page.Next
	cfg.Previous = page.Previous

	return nil
}

func commandMapb(cfg *Config) error {
	url := cfg.Previous
	if url == "" {
		fmt.Println("you'r on the first page")
		url = ""
		return nil
	}

	page, err := pokeapi.PokedexAPI(url)
	if err != nil {
		return err
	}

	for _, r := range page.Results {
		fmt.Println(r.Name)
	}

	cfg.Previous = page.Previous
	cfg.Next = page.Next

	return nil
}

func explore(loc string) {
	
}

func startRepl(cfg *Config) {
	commandHelp(cfg)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Pokedex > ")
		scanner.Scan()
		
		userText := scanner.Text()
		wordList := cleanInput(userText)
		if len(wordList) == 0 {
			continue
		}
		
		commandWord := wordList[0]

		command, exists := cliCommands[commandWord]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(cfg)
		if err != nil {
			fmt.Println(err)
		}
		
	}
}