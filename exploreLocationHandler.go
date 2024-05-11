package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LocationInfomrationType struct {
	Areas []LocationType `json:"areas"`
}

var locCache = getNewCache[LocationInfomrationType](CacheConfig{
	DeleteInterval: 5 * time.Minute,
})

func exploreLocationCallback(locationInput string) error {
	fmt.Println("Exploring location " + locationInput + "...")

	var locInfo LocationInfomrationType
	if locCache.Get(locationInput) != nil {
		locInfo = *locCache.Get(locationInput)
	} else {
		url := fmt.Sprintf("%vlocation/%v/", baseUrl, locationInput)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("communication error occured. Error: %v", err.Error())
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			fmt.Println("Location " + locationInput + " was not found")
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("communication error occured. HTTP code: %v", resp.StatusCode)
		}

		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &locInfo)
		if err != nil {
			return fmt.Errorf("failed to parse response. Error: %v", err.Error())
		}

		locCache.Add(locationInput, &locInfo)
	}

	fmt.Println("Found next areas: ")
	for _, areaInfo := range locInfo.Areas {
		fmt.Println(" - " + areaInfo.Name)
	}

	return nil
}
