package colony

import (
	"fmt"
	"strings"
)

// Simulate ants moving along multiple paths.
func simulateAnts(paths [][]string, assignments [][]int) {

	lastTurn := 0

	// Calculate when the final ant will arrive.
	for i, ants := range assignments {

		if len(ants) == 0 {
			continue
		}

		finish := len(paths[i]) + len(ants) - 2

		if finish > lastTurn {
			lastTurn = finish
		}
	}

	// Simulate one turn at a time.
	for turn := 1; turn <= lastTurn; turn++ {

		moves := []string{}

		for i, path := range paths {

			ants := assignments[i]

			// Find which ants are moving this turn.
			first := turn - (len(path) - 1)

			if first < 0 {
				first = 0
			}

			last := turn - 1

			if last >= len(ants) {
				last = len(ants) - 1
			}

			// Move every active ant one room forward.
			for j := first; j <= last; j++ {

				step := turn - j

				moves = append(moves,
					fmt.Sprintf("L%d-%s", ants[j], path[step]))
			}
		}

		// Print all moves for this turn.
		fmt.Println(strings.Join(moves, " "))
	}
}
