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
		"inspect": {
			name: 		 "inspect",
			description: "Inspect a pokemon in your pokedex by name (use example: inspect golbat)",
			callback:	 inspect,
		},
		"pokedex": {
			name:		 "pokedex",
			description: "inspect all catched pokemon",
			callback:	 pokedex,
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
	maxChance := 60.0

	pokemonName := args[0]
	pokemon, err := pokeapi.PokemonAPI(pokemonName)
	if err != nil {
		fmt.Println("There's no such pokemon...")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)

	normalized := (float64(pokemon.BaseExperience) - minXP) / (maxXP - minXP)
	inverted := 1 - normalized
	catchChance := inverted * maxChance

	if catchChance < minChance {
		catchChance = minChance
	}

	rand.Seed(time.Now().UnixNano())
	roll := rand.Float64() * 100
	if roll <= catchChance {
		fmt.Printf("%v was caught!\n", pokemon.Name)

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

func inspect(cfg *Config, args []string) error {
	pokemon := args[0]

	p, ok := Pokedex[pokemon]
	if ok == false {
		fmt.Println("You have yet to catch this Pokemon")
		fmt.Println("Or you misspelled. :)")
		return nil
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)

	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func pokedex(cfg *Config, args []string) error {

	if len(Pokedex) == 0 {
		fmt.Println("Your pokedex is empty")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for p, _ := range Pokedex {
		fmt.Printf(" - %v\n", p)
	}
	return nil
}

func startRepl(cfg *Config) {
	commandHelp(cfg, nil)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex > ")
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