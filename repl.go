package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/deimerin/pkdxcli/internal/pokeapi"
)

type config struct {
	Next     *string
	Previous *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var commands map[string]cliCommand
var cfg config

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Print locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Print previous locations",
			callback:    commandMapB,
		},
	}
	url := "https://pokeapi.co/api/v2/location-area/"
	cfg = config{
		Next:     &url,
		Previous: nil,
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func start() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			words := cleanInput(scanner.Text())

			if len(words) > 0 {

				if command, exists := commands[words[0]]; exists {
					command.callback(&cfg)
				} else {
					fmt.Println("Unknown command")
				}

			}

		}

	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	locations, next, previous, err := pokeapi.FetchLocations(*cfg.Next)

	if err != nil {
		fmt.Println("Something went wrong on the API call")
		return err
	}

	for _, location := range locations {
		fmt.Println(location)
	}
	cfg.Previous = &previous
	cfg.Next = &next

	return nil

}

func commandMapB(cfg *config) error {
	if *cfg.Previous != "" {
		locations, next, previous, err := pokeapi.FetchLocations(*cfg.Previous)

		if err != nil {
			fmt.Println("Something went wrong on the API call")
			return err
		}

		for _, location := range locations {
			fmt.Println(location)
		}

		cfg.Previous = &previous
		cfg.Next = &next

	} else {
		fmt.Println("you're on the first page")
	}

	return nil

}
