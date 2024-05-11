package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseUrl = "https://pokeapi.co/api/v2/"
)

type LocationType struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationSetType struct {
	Count   int            `json:"count"`
	Results []LocationType `json:"results"`
}

var (
	currentOffset           = 0
	numberOfElements        = 20
	maximumNumberOfElements = 32000

	savedLocations = make([]LocationType, 20)

	mapCache = getNewCache[LocationSetType](CacheConfig{
		DeleteInterval: 5 * time.Minute,
	})
)

func getLocations(isForward bool) (*LocationSetType, error) {
	offset := currentOffset
	if !isForward {
		offset -= numberOfElements
	}

	url := fmt.Sprintf("%slocation?offset=%v&limit=%v", baseUrl, offset, numberOfElements)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("communication error occured. Error: %v", err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("communication error occured. HTTP code: %v", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	var locs *LocationSetType
	err = json.Unmarshal(data, &locs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response. Error: %v", err.Error())
	}

	return locs, nil
}

func mapCallback(_ string) error {
	if currentOffset > maximumNumberOfElements {
		fmt.Println("No more data to fetch")
		return nil
	}

	var locs *LocationSetType

	if mapCache.Get(fmt.Sprint(currentOffset)) != nil {
		locs = mapCache.Get(fmt.Sprint(currentOffset))
	} else {
		l, err := getLocations(true)
		if err != nil {
			return err
		}
		mapCache.Add(fmt.Sprint(currentOffset), l)
		locs = l
	}

	maximumNumberOfElements = locs.Count
	savedLocations = locs.Results

	currentOffset += numberOfElements

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func mapBackCallback(_ string) error {
	if currentOffset < numberOfElements*2 {
		fmt.Println("No more data to fetch")
		return nil
	}

	var locs *LocationSetType

	currentOffset -= numberOfElements
	if mapCache.Get(fmt.Sprint(currentOffset-numberOfElements)) != nil {
		locs = mapCache.Get(fmt.Sprint(currentOffset - numberOfElements))
	} else {

		l, err := getLocations(false)
		if err != nil {
			currentOffset += numberOfElements
			return err
		}
		mapCache.Add(fmt.Sprint(currentOffset-numberOfElements), l)
		locs = l
	}

	maximumNumberOfElements = locs.Count
	savedLocations = locs.Results

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}

	return nil
}
