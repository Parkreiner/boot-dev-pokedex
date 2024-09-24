package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const pokedexLocationsEndpoint = "https://pokeapi.co/api/v2/location-area"

type PokedexCommand struct {
	command     string
	description string
	operation   func() error
}

type PokemonLocation struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonLocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`

	Results []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func PokedexCommands() map[string]PokedexCommand {
	const mapPageSize = 20
	mapOffset := 0

	exit := func() error {
		os.Exit(0)
		return nil
	}

	makeRequest := func() (PokemonLocationResponse, error) {
		endpoint := pokedexLocationsEndpoint + "?offset=" + strconv.Itoa(mapOffset)

		res, err := http.Get(endpoint)
		if err != nil {
			return PokemonLocationResponse{}, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return PokemonLocationResponse{}, err
		}

		var pokeRes PokemonLocationResponse
		err = json.Unmarshal(body, &pokeRes)
		if err != nil {
			return PokemonLocationResponse{}, err
		}

		return pokeRes, nil
	}

	var cmdMap map[string]PokedexCommand
	cmdMap = map[string]PokedexCommand{
		"map": {
			command:     "map",
			description: "List 20 regions. The next call to map will list the next 20 regions (or however many remain).",
			operation: func() error {
				res, err := makeRequest()
				if err != nil {
					return err
				}

				for _, r := range res.Results {
					fmt.Println(r.Name)
				}

				mapOffset += mapPageSize
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

				mapOffset -= mapPageSize
				res, err := makeRequest()
				if err != nil {
					return err
				}

				for _, r := range res.Results {
					fmt.Println(r.Name)
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
