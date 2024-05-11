package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PokemonSmallInformationType struct {
	Pokemon struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"pokemon"`
}

type AreaInformationType struct {
	PokemonEncounters []PokemonSmallInformationType `json:"pokemon_encounters"`
}

var areaCache = getNewCache[AreaInformationType](CacheConfig{
	DeleteInterval: 5 * time.Minute,
})

func exploreAreaCallback(areaInput string) error {
	fmt.Println("Exploring area " + areaInput + "...")

	var areaInfo AreaInformationType
	if areaCache.Get(areaInput) != nil {
		areaInfo = *areaCache.Get(areaInput)
	} else {
		url := fmt.Sprintf("%vlocation-area/%v/", baseUrl, areaInput)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("communication error occured. Error: %v", err.Error())
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			fmt.Println("Area " + areaInput + " was not found")
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("communication error occured. HTTP code: %v", resp.StatusCode)
		}

		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &areaInfo)
		if err != nil {
			return fmt.Errorf("failed to parse response. Error: %v", err.Error())
		}

		areaCache.Add(areaInput, &areaInfo)
	}

	fmt.Println("Found next pokemons:")
	for _, pokeInfo := range areaInfo.PokemonEncounters {
		fmt.Println(" - " + pokeInfo.Pokemon.Name)
	}

	return nil
}
