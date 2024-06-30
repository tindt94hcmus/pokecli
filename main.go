package main

import (
	"bufio"
	"fmt"
	"github.com/tindt94hcmus/pokedexcli/internal/pokeapi"
	"os"
)

type config struct {
	Next string
	Prev string
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func commandHelp() error {
	fmt.Println("Available commands: ")
	for _, cmd := range getCommands(&config{}) {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Exiting Pokedex...")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) func() error {
	return func() error {
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
func commandMapBack(cfg *config) func() error {
	return func() error {
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

		if cmd, exists := commands[input]; exists {
			err := cmd.callback()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}

		} else {
			fmt.Println("Unknown command. Type 'help' for a list of available commands.")
		}
	}
}
