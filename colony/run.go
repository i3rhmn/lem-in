package colony

import (
	"errors"
	"fmt"
	"strings"
)

// Run executes the complete ant colony simulation.
func Run(data []byte) error {

	// Parse the input
	farm, err := parseColony(data)
	if err != nil {
		return err
	}

	// Find optimized paths
	paths := findOptimizedPaths(
		farm.Graph,
		farm.Start,
		farm.End,
		farm.Ants,
	)

	if len(paths) == 0 {
		return errors.New("no valid path")
	}

	// Assign ants to the paths
	assignments := assignAnts(farm.Ants, paths)

	// Display the original input
	fmt.Print(string(data))

	if !strings.HasSuffix(string(data), "\n") {
		fmt.Println()
	}

	fmt.Println()

	// Run the simulation
	simulateAnts(paths, assignments)

	return nil
}
