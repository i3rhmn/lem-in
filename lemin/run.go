package lemin

import (
	"fmt"
	"os"
	"strings"
)

// Run reads the colony file and prints the ant movements.
func Run(args []string) {
	if len(args) == 0 {
		fmt.Println("ERROR: invalid data format, missing file argument; usage: go run . <file.txt>")
		return
	}
	if len(args) != 1 {
		fmt.Println("ERROR: invalid data format, too many arguments; provide exactly one file")
		return
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		switch {
		case os.IsNotExist(err):
			fmt.Printf("ERROR: invalid data format, file not found: %q\n", args[0])
		case os.IsPermission(err):
			fmt.Printf("ERROR: invalid data format, permission denied: %q\n", args[0])
		default:
			fmt.Printf("ERROR: invalid data format, cannot read file %q: %v\n", args[0], err)
		}
		return
	}

	colony, err := parseColony(data)
	if err != nil {
		fmt.Printf("ERROR: invalid data format, %v\n", err)
		return
	}
	paths := findOptimizedPaths(colony.Graph, colony.Start, colony.End, colony.Ants)
	if len(paths) == 0 {
		fmt.Printf("ERROR: invalid data format, no path from start room %q to end room %q\n", colony.Start, colony.End)
		return
	}

	assignments := assignAnts(colony.Ants, paths)
	fmt.Print(string(data))
	if !strings.HasSuffix(string(data), "\n") {
		fmt.Println()
	}
	fmt.Println()
	simulateAnts(paths, assignments)
}
