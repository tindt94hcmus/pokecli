package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tindt94hcmus/pokedexcli/internal/pokecache"
)

const ENDPOINT = "https://pokeapi.co/api/v2/location-area"

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
		url = ENDPOINT
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
	url := fmt.Sprintf("%s/%s", ENDPOINT, areaName)

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
