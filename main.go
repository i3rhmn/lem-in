package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	// Require exactly one file argument.
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Read the input file.
	data, err := os.ReadFile(os.Args[1])

	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Parse and validate the colony.
	colony, err := parseColony(data)

	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Find optimized paths.
	paths := findOptimizedPaths(
		colony.Graph,
		colony.Start,
		colony.End,
		colony.Ants,
	)

	if len(paths) == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Distribute ants across the paths.
	assignments := assignAnts(colony.Ants, paths)

	// Display the original input.
	fmt.Print(string(data))

	if !strings.HasSuffix(string(data), "\n") {
		fmt.Println()
	}

	fmt.Println()

	// Simulate ant movement.
	simulateAnts(paths, assignments)
}
