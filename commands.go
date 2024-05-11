package main

import (
	"fmt"
	"os"
)

func helpCallback(m map[string]*Command) func(_ string) error {
	str := ""

	for _, command := range m {
		str += fmt.Sprintf("%s: %s\n", command.name, command.description)
	}

	return func(_ string) error {
		fmt.Println(str)

		return nil
	}
}

func exitCallback(_ string) error {
	fmt.Println("Goodbye!")
	os.Exit(0)

	return nil
}

func getCommandsMap() map[string]*Command {
	m := make(map[string]*Command)

	m["map"] = &Command{
		name:        "map",
		description: "Gets next 20 locations",
		callback:    mapCallback,
	}

	m["mapb"] = &Command{
		name:        "mapb",
		description: "Gets previous 20 locations",
		callback:    mapBackCallback,
	}

	m["explore-location"] = &Command{
		name:        "explore-location",
		description: "Explores input location. Usage example: \"explore-location example-location\"",
		callback:    exploreLocationCallback,
	}

	m["explore-area"] = &Command{
		name:        "explore-area",
		description: "Explores input area. Usage example: \"explore-area example-area\"",
		callback:    exploreAreaCallback,
	}

	m["catch"] = &Command{
		name:        "explore-pokemon",
		description: "Explores species of given pokemon. Usage example: \"explore-area example-pokemon\"",
		callback:    catchPokemonCallback,
	}

	m["pokedex"] = &Command{
		name:        "Pokedex",
		description: "Show pokedex",
		callback:    showPokedexCallback,
	}

	m["help"] = &Command{
		name:        "help",
		description: "Displays help message",
		callback:    nil,
	}

	m["exit"] = &Command{
		name:        "exit",
		description: "Exits application",
		callback:    exitCallback,
	}

	m["help"].callback = helpCallback(m)

	return m
}
