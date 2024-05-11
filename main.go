package main

import (
	"bufio"
	"fmt"
	"os"
)

type Command struct {
	name        string
	description string
	callback    func() error
}

var commandsMap = getCommandsMap()

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome! Write \"help\" to get help message")
	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		str := scanner.Text()
		comm := commandsMap[str]

		if comm == nil {
			fmt.Println("Unknown command. Write \"help\" to see available commands")
			continue
		}

		err := comm.callback()
		if err != nil {
			fmt.Printf("Error occured: %s\n", err.Error())
		}
	}
}
