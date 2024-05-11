package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	name        string
	description string
	callback    func(input string) error
}

var commandsMap = getCommandsMap()

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome! Write \"help\" to get help message")
	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		inputs := strings.Split(scanner.Text(), " ")

		if len(inputs) < 1 || len(inputs) > 3 {
			fmt.Println("Unknown command. Write \"help\" to see available commands")
			continue
		}

		comm := commandsMap[inputs[0]]

		if comm == nil {
			fmt.Println("Unknown command. Write \"help\" to see available commands")
			continue
		}

		var err error
		if len(inputs) == 1 {
			err = comm.callback("")
		} else {
			err = comm.callback(inputs[1])
		}

		if err != nil {
			fmt.Printf("Error occured: %s\n", err.Error())
		}
	}
}
