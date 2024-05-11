package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type PokemonInfomrationType struct {
	CaptureRate int    `json:"capture_rate"`
	Name        string `json:"name"`
}

var pokemonCache = getNewCache[PokemonInfomrationType](CacheConfig{
	DeleteInterval: 5 * time.Minute,
})

func catchPokemonCallback(pokemonInput string) error {
	fmt.Println("Trying to catch " + pokemonInput + "...")

	var pokeInfo PokemonInfomrationType
	if pokemonCache.Get(pokemonInput) != nil {
		pokeInfo = *pokemonCache.Get(pokemonInput)
	} else {
		url := fmt.Sprintf("%vpokemon-species/%v/", baseUrl, pokemonInput)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("communication error occured. Error: %v", err.Error())
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			fmt.Println("Pokemon " + pokemonInput + " was not found")
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("communication error occured. HTTP code: %v", resp.StatusCode)
		}

		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &pokeInfo)
		if err != nil {
			return fmt.Errorf("failed to parse response. Error: %v", err.Error())
		}

		pokemonCache.Add(pokemonInput, &pokeInfo)
	}

	captureRateNormalized := pokeInfo.CaptureRate * 100 / 255
	fmt.Println(captureRateNormalized)
	if int32(captureRateNormalized) > rand.Int31n(101) {
		fmt.Println("Captured " + pokemonInput + "!")
		pokedex = append(pokedex, pokeInfo.Name)
	} else {
		fmt.Println("Failed to capture " + pokemonInput + "...")
	}

	return nil
}
