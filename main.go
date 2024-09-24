package main

import (
	"errors"
	"fmt"
	"os"
)

type PokedexCommand struct {
	command     string
	description string
	operation   func() error
}

func PokedexCommands() map[string]PokedexCommand {
	mapOffset := 0
	exit := func() error {
		os.Exit(0)
		return nil
	}

	var cmdMap map[string]PokedexCommand
	cmdMap = map[string]PokedexCommand{
		"map": {
			command:     "map",
			description: "List 20 regions. The next call to map will list the next 20 regions (or however many remain).",
			operation: func() error {
				return nil
			},
		},
		"mapb": {
			command:     "mapb",
			description: "List the previous 20 regions.",
			operation: func() error {
				if mapOffset == 0 {
					return errors.New("can't go back any further in map list")
				}

				return nil
			},
		},
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

		cmd, ok := cmdMap[input]
		if !ok {
			fmt.Println("Sorry, don't know that command")
			continue
		}

		err := cmd.operation()
		if err != nil {
			fmt.Printf("error - %v\n", err)
			if input == "exit" || input == "quit" {
				os.Exit(1)
			}
		}
	}
}
