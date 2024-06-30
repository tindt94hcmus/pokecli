package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tindt94hcmus/pokedexcli/internal/pokeapi"
)

type config struct {
	Next string
	Prev string
}

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func getCommands(cfg *config) map[string]cliCommand {
	return map[string]cliCommand{
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
			description: "Display the name of 20 location areas in the Pokedex world",
			callback:    commandMap(cfg),
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world",
			callback:    commandMapBack(cfg),
		},
		"explore": {
			name:        "explore",
			description: "Displays the list of Pokémon in a given location area",
			callback:    commandExplore,
		},
	}
}

func commandHelp(args []string) error {
	fmt.Println("Available commands: ")
	for _, cmd := range getCommands(&config{}) {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(args []string) error {
	fmt.Println("Exiting Pokedex...")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) func(args []string) error {
	return func(args []string) error {
		response, err := pokeapi.FetchLocationAreas(cfg.Next)
		if err != nil {
			return fmt.Errorf("error fetching location areas: %v", err)
		}

		for _, area := range response.Results {
			fmt.Println(area.Name)
		}

		if response.Next != nil {
			cfg.Next = *response.Next
		}

		if response.Previous != nil {
			cfg.Prev = *response.Previous
		}

		return nil
	}
}
func commandMapBack(cfg *config) func(args []string) error {
	return func(args []string) error {
		if cfg.Prev == "" {
			fmt.Println("No previous locations available.")
			return nil
		}

		response, err := pokeapi.FetchLocationAreas(cfg.Prev)
		if err != nil {
			return fmt.Errorf("error fetching location areas: %v", err)
		}

		for _, area := range response.Results {
			fmt.Println(area.Name)
		}

		if response.Next != nil {
			cfg.Next = *response.Next
		}

		if response.Previous != nil {
			cfg.Prev = *response.Previous
		} else {
			cfg.Prev = ""
		}

		return nil
	}
}

func commandExplore(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a location area name")
	}

	areaName := args[0]
	response, err := pokeapi.FetchPokemonInArea(areaName)
	if err != nil {
		return fmt.Errorf("error fetching Pokemon in area %s: %v", areaName, err)
	}

	for _, pokemon := range response.Pokemon {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func main() {
	cfg := &config{}
	commands := getCommands(cfg)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		commandName := parts[0]
		args := parts[1:]

		if cmd, exists := commands[commandName]; exists {
			err := cmd.callback(args)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}

		} else {
			fmt.Println("Unknown command. Type 'help' for a list of available commands.")
		}
	}
}
