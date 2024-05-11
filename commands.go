package main

import (
	"fmt"
	"os"
)

func helpCallback(m map[string]*Command) func() error {
	str := "\n"

	for _, command := range m {
		str += fmt.Sprintf("%s: %s\n", command.name, command.description)
	}

	return func() error {
		fmt.Println(str)

		return nil
	}
}

func exitCallback() error {
	fmt.Println("\nGoodbye")
	os.Exit(0)

	return nil
}

func getCommandsMap() map[string]*Command {
	m := make(map[string]*Command)

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