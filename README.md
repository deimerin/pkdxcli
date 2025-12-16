# Pokedex CLI

A lightweight command-line Pokedex application built in Go that lets you explore Pokemon locations and catch Pokemon using the PokeAPI.

## Features

- **Explore Locations**: Browse Pokemon location areas from the PokeAPI
- **Catch Pokemon**: Attempt to catch wild Pokemon with a randomized catch rate
- **Pokedex**: View all Pokemon you've caught
- **Inspect**: View detailed stats and information about caught Pokemon
- **Navigation**: Browse forward and backward through location pages

## Commands

- `help` - Display available commands
- `map` - Print next set of location areas
- `mapb` - Print previous set of location areas
- `explore <area>` - Explore a specific Pokemon area and see available Pokemon
- `catch <pokemon>` - Attempt to catch a Pokemon
- `inspect <pokemon>` - View details about a caught Pokemon
- `pokedex` - List all caught Pokemon
- `exit` - Exit the application

## Running

```bash
go run .
```

## Project Structure

- `main.go` - Entry point
- `repl.go` - REPL (Read-Eval-Print Loop) and command implementations
- `internal/pokeapi/` - PokeAPI wrapper and data structures
- `internal/pokecache/` - Caching layer for API responses