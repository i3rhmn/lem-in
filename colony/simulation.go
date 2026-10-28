package colony

import (
	"errors"
	"fmt"
	"strings"
)

// Give each next ant to the path where it can arrive earliest.
func assignAnts(ants int, paths [][]string) [][]int {
	assigned := make([][]int, len(paths))
	for ant := 1; ant <= ants; ant++ {
		best := 0
		for i := 1; i < len(paths); i++ {
			if len(paths[i])+len(assigned[i]) < len(paths[best])+len(assigned[best]) {
				best = i
			}
		}
		assigned[best] = append(assigned[best], ant)
	}
	// assigned = [][]int{
    //		{1, 2, 3},
	// }
	return assigned
}

// On turn t, ant j on a path moves to the room at index t-j.
// print aunts movments turn by turn
func simulateAnts(paths [][]string, assigned [][]int) {
	// paths = [][]string{
    // 		{"A", "B", "C"},
	// }

	// assigned = [][]int{
    //		{1, 2, 3},
	// }
	lastTurn := 0
	for i, ants := range assigned {
		if len(ants) > 0 && len(paths[i])+len(ants)-2 > lastTurn {
			lastTurn = len(paths[i]) + len(ants) - 2
		} // calc how many turns required
	}

	for turn := 1; turn <= lastTurn; turn++ {
		moves := []string{}
		for i, path := range paths {
			first := turn - len(path) + 1
			if first < 0 {
				first = 0
			}
			last := turn - 1
			if last >= len(assigned[i]) {
				last = len(assigned[i]) - 1
			}
			for j := first; j <= last; j++ {
				moves = append(moves, fmt.Sprintf("L%d-%s", assigned[i][j], path[turn-j]))
			}
		}
		fmt.Println(strings.Join(moves, " "))
	}
}	// 1	Ant 1 moves A → B	L1-B
	// 2	Ant 1 moves B → C; ant 2 enters B	L1-C L2-B
	// 3	Ant 2 reaches C; ant 3 enters B	L2-C L3-B
	// 4	Ant 3 reaches C	L3-C
	

// Run joins parsing, pathfinding, assignment, and printing.
func Run(data []byte) error {
	farm, err := parseColony(data)
	if err != nil {
		return err
	}
	paths := findOptimizedPaths(farm.Graph, farm.Start, farm.End, farm.Ants)
	if len(paths) == 0 {
		return errors.New("no path from start to end")
	}
	// Returned value: [["A", "B", "C"]]

	fmt.Print(string(data))
	if !strings.HasSuffix(string(data), "\n") {
		fmt.Println()
	}
	fmt.Println()
	simulateAnts(paths, assignAnts(farm.Ants, paths)) // two func in same file
	return nil
}
