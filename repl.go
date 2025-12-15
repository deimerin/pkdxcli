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
	callback    func(*config, []string) error
}

var commands map[string]cliCommand
var cfg config
var baseLocationAreaURL string = "https://pokeapi.co/api/v2/location-area/"
var pokedex map[string]pokeapi.Pokemon

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
		"explore": {
			name:        "explore",
			description: "explore the area, etc",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "try to catch a Pokemon",
			callback:    commandCatch,
		},
	}

	pokedex = make(map[string]pokeapi.Pokemon)

	cfg = config{
		Next:     &baseLocationAreaURL,
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
					if err := command.callback(&cfg, words[1:]); err != nil {
						fmt.Println("Error:", err)
					}
				} else {
					fmt.Println("Unknown command")
				}

			}

		}

	}
}

func commandExit(cfg *config, param []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, param []string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config, param []string) error {
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

func commandMapB(cfg *config, param []string) error {
	if cfg.Previous != nil && *cfg.Previous != "" {
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

func commandExplore(cfg *config, param []string) error {
	if len(param) > 0 {

		areaName := strings.Join(param, "-")

		newURL := baseLocationAreaURL + areaName

		pokeList, err := pokeapi.FetchLocationArea(newURL)

		if err != nil {
			fmt.Println("Something went wrong fetching the area data")
			return err
		}

		fmt.Printf("Exploring %s...\n", param[0])
		fmt.Println("Found Pokemon:")
		for _, pkm := range pokeList {
			fmt.Printf(" - %s\n", pkm)
		}
	}

	if len(param) == 0 {
		fmt.Println("You must provide an area name. Example: explore canalave-city-area ")
	}

	return nil
}

func commandCatch(cfg *config, param []string) error {

	if len(param) == 0 {
		fmt.Println("Provide a valid pokemon name. Example: catch pikachu")
		return nil
	}

	if len(param) > 0 {

		wildPokemon := param[0]

		pokemon, err := pokeapi.GetPokemon(wildPokemon)

		if err != nil {
			fmt.Println("Invalid pokemon name/id")
			return err
		}

		fmt.Printf("Throwing a Pokeball at %s...\n", wildPokemon)

		if pokemon.Catch() {
			fmt.Printf("%s was caught!\n", wildPokemon)
			pokedex[wildPokemon] = pokemon
		} else {
			fmt.Printf("%s escaped!\n", wildPokemon)
		}
	}
	return nil
}
