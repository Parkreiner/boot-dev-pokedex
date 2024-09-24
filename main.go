package main

import (
	"fmt"
	"os"
)

type PokedexCommand struct {
	command     string
	description string
	operation   func() error
}

func PokedexCommands() map[string]PokedexCommand {
	var cmdMap map[string]PokedexCommand

	exit := func() error {
		os.Exit(0)
		return nil
	}

	cmdMap = map[string]PokedexCommand{
		"exit": {
			command:     "exit",
			description: "Exit the program",
			operation:   exit,
		},
		"quit": {
			command:     "quit",
			description: "Exit the program (alias for exit)",
			operation:   exit,
		},
		"help": {
			command:     "help",
			description: "List full list of all operations",
			operation: func() error {
				fmt.Println("\nWelcome to the pokedex!")
				fmt.Println("Usage: ")

				for _, cmd := range cmdMap {
					fmt.Printf("%s - %s\n", cmd.command, cmd.description)
				}

				return nil
			},
		},
	}

	return cmdMap
}

func main() {
	var input string
	cmdMap := PokedexCommands()

	for {
		fmt.Print(("pokedex> "))
		fmt.Scanln(&input)

		if cmd, ok := cmdMap[input]; ok {
			err := cmd.operation()
			if err != nil {
				os.Exit(1)
			}
		} else {
			fmt.Println("Sorry, don't know that command")
		}
	}
}
