package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tindt94hcmus/pokedexcli/internal/pokecache"
)

const ENDPOINT = "https://pokeapi.co/api/v2"

type LocationArea struct {
	Name string `json:"name"`
}

type LocationAreaResponse struct {
	Results  []LocationArea `json:"results"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
}

var cache = pokecache.NewCache(5 * time.Minute)

func FetchLocationAreas(url string) (LocationAreaResponse, error) {
	if url == "" {
		url = fmt.Sprintf("%s/%s", ENDPOINT, "location-area")
	}

	if data, found := cache.Get(url); found {
		var locationAreas LocationAreaResponse
		err := json.Unmarshal(data, &locationAreas)
		if err == nil {
			return locationAreas, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	var locationAreas LocationAreaResponse
	err = json.Unmarshal(body, &locationAreas)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	cache.Add(url, body)

	return locationAreas, nil
}

type Pokemon struct {
	Name string `json:"name"`
}

type PokemonEntry struct {
	Pokemon Pokemon `json:"pokemon"`
}

type PokemonResponse struct {
	Pokemon []PokemonEntry `json:"pokemon_encounters"`
}

func FetchPokemonInArea(areaName string) (PokemonResponse, error) {
	url := fmt.Sprintf("%s/location-area/%s", ENDPOINT, areaName)

	if data, found := cache.Get(url); found {
		var pokemonResponse PokemonResponse
		err := json.Unmarshal(data, &pokemonResponse)
		if err == nil {
			return pokemonResponse, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return PokemonResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonResponse{}, err
	}

	var pokemonResponse PokemonResponse
	err = json.Unmarshal(body, &pokemonResponse)
	if err != nil {
		return PokemonResponse{}, err
	}

	cache.Add(url, body)

	return pokemonResponse, nil
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type Type struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type PokemonData struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []Stat `json:"stats"`
	Types          []Type `json:"types"`
}

func FetchPokemonByName(pokemonName string) (PokemonData, error) {
	url := fmt.Sprintf("%s/pokemon/%s", ENDPOINT, pokemonName)

	if data, found := cache.Get(url); found {
		var pokemonData PokemonData
		err := json.Unmarshal(data, &pokemonData)
		if err == nil {
			return pokemonData, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return PokemonData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonData{}, err
	}

	var pokemonData PokemonData
	err = json.Unmarshal(body, &pokemonData)
	if err != nil {
		return PokemonData{}, err
	}

	fmt.Printf("BaseExperience: %v\n", pokemonData.BaseExperience)

	cache.Add(url, body)

	return pokemonData, nil
}
