package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	commandHelp()
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

		err := command.callback()
		if err != nil {
			fmt.Println(err)
		}
		
	}
}