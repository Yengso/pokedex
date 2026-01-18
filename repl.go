package main

import (
	"fmt"
	"bufio"
	"strings"
	"time"
	"os"
	"math/rand"
	"github.com/yengso/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name		string
	description	string
	callback	func(*Config, []string) error
}

type Config struct {
	Next		string
	Previous 	string
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
		"explore": {
			name: 		 "explore",
			description: "explore a location to find pokemon. (use example: explore eterna-city-area)",
			callback: 	 explore,
		},
		"catch": {
			name:		 "catch",
			description: "Try to catch named pokemon. (use example: catch tentacruel)",
			callback:	 catch,
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

func commandExit(cfg *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for name, cmd := range cliCommands {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

func commandMap(cfg *Config, args []string) error {
	url := cfg.Next
	if url == "" {
		url = ""
	}

	page, err := pokeapi.LocationsAPI(url)
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

func commandMapb(cfg *Config, args []string) error {
	url := cfg.Previous
	if url == "" {
		fmt.Println("you'r on the first page")
		url = ""
		return nil
	}

	page, err := pokeapi.LocationsAPI(url)
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

func explore(cfg *Config, args []string) error {
	if len(args) == 0 {
		fmt.Println("You must provide a location area name.")
		return nil
	}

	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	loc, err := pokeapi.ExploreAPI(areaName)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, enc := range loc.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}

	return nil
}

var Pokedex = make(map[string]pokeapi.Pokemon)

func catch(cfg *Config, args []string) error {
	minXP := 52.0
	maxXP := 608.0
	minChance := 0.1
	maxChance := 75.0

	pokemonName := args[0]
	pokemon, _ := pokeapi.PokemonAPI(pokemonName)

	normalized := (float64(pokemon.BaseExperience) - minXP) / (maxXP - minXP)
	inverted := 1 - normalized
	catchChance := inverted * maxChance

	if catchChance < minChance {
		catchChance = minChance
	}

	rand.Seed(time.Now().UnixNano())
	roll := rand.Float64() * 100

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)

	if roll <= catchChance {
		fmt.Printf("%v was caught!\nclea", pokemon.Name)

		if _, exists := Pokedex[pokemon.Name]; exists {
			fmt.Println("You alaready have this Pokemon, so you let this one go. :)")
		} else {
			Pokedex[pokemon.Name] = pokemon
			fmt.Println("A new Pokemon has been added to you Pokedex. :)")
		}
	}
	if roll > catchChance {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}

	return nil
}

func startRepl(cfg *Config) {
	commandHelp(cfg, nil)
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
		args := wordList[1:]

		command, exists := cliCommands[commandWord]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(cfg, args)
		if err != nil {
			fmt.Println(err)
		}
		
	}
}