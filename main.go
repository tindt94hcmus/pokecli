package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/tindt94hcmus/pokedexcli/internal/pokeapi"
)

type config struct {
	Next    string
	Prev    string
	Pokedex map[string]pokeapi.PokemonData
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
		"catch": {
			name:        "catch",
			description: "Tries to catch a Pokemon by name",
			callback:    commandCatch(cfg),
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect details of a caught Pokémon",
			callback:    commandInspect(cfg),
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

func commandCatch(cfg *config) func(args []string) error {
	return func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please provide Pokemon name.")
		}

		pokemonName := args[0]
		fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
		pokemonData, err := pokeapi.FetchPokemonByName(pokemonName)

		if err != nil {
			return fmt.Errorf("error fetching Pokemon in area %s: %v", pokemonName, err)
		}

		rand.New(rand.NewSource(time.Now().UnixNano()))
		catchProbability := 100.0 / float64(pokemonData.BaseExperience)
		randomNumber := rand.Float64()
		fmt.Printf("catchProbability: %v\n", catchProbability)
		fmt.Printf("random: %v\n", randomNumber)
		if randomNumber <= catchProbability {
			fmt.Printf("You caught %s!\n", pokemonData.Name)
			cfg.Pokedex[pokemonData.Name] = pokemonData
		} else {
			fmt.Printf("%s escaped!\n", pokemonData.Name)
		}

		commandList(cfg)
		return nil
	}

}

func commandList(cfg *config) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("No Pokémon caught yet.")
		return nil
	}

	fmt.Println("Caught Pokémon:")
	for _, pokemon := range cfg.Pokedex {
		fmt.Printf("- %s (Base Experience: %d)\n", pokemon.Name, pokemon.BaseExperience)
	}
	return nil
}

func commandInspect(cfg *config) func(args []string) error {
	return func(args []string) error {

		if len(args) == 0 {
			fmt.Println("please provide a pokemon name")
		}

		pokemonName := args[0]
		pokemon, ok := cfg.Pokedex[pokemonName]

		if !ok {
			fmt.Printf("You have not caught %s yet.\n", pokemonName)
			return nil
		}

		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)

		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Println("Types:")
		for _, pokeType := range pokemon.Types {
			fmt.Printf("- %s\n", pokeType.Type.Name)
		}

		return nil
	}
}

func main() {
	cfg := &config{
		Pokedex: make(map[string]pokeapi.PokemonData),
	}
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
