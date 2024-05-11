package main

import "fmt"

var pokedex = make([]string, 0)

func showPokedexCallback(_ string) error {
	fmt.Println("Your pokedex: ")
	for _, pokemon := range pokedex {
		fmt.Println(" - " + pokemon)
	}

	return nil
}
