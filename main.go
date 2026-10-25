package main

import (
	"fmt"
	"os"

	"lem-in/colony"
)

func main() {
	// Validate arguments
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Read input file
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Execute the colony simulation
	if err := colony.Run(data); err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}
}
