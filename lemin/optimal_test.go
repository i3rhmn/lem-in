package lemin

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

// Find the true minimum turns by exploring all legal moves.
// Only suitable for small graphs.
func shortestAntSchedule(
	graph map[string][]string,
	ants int,
) int {

	rooms := []string{"S", "E"}

	for room := range graph {
		if room != "S" && room != "E" {
			rooms = append(rooms, room)
		}
	}

	sort.Strings(rooms[2:])

	// Give each room a numeric ID.
	ids := make(map[string]int)

	for i, room := range rooms {
		ids[room] = i
	}

	adj := make([][]int, len(rooms))

	for i, room := range rooms {
		for _, next := range graph[room] {
			adj[i] = append(adj[i], ids[next])
		}
	}

	// All ants start in room S (index 0).
	initial := make([]byte, ants)

	// The target is room E (index 1).
	goal := make([]byte, ants)

	for i := range goal {
		goal[i] = 1
	}

	queue := []string{string(initial)}
	distance := map[string]int{string(initial): 0}

	for head := 0; head < len(queue); head++ {

		state := queue[head]
		turn := distance[state]

		if state == string(goal) {
			return turn
		}

		next := make([]byte, ants)
		usedTunnels := make(map[[2]int]bool)

		var explore func(int, bool)

		explore = func(ant int, moved bool) {

			// All ants have made their decisions.
			if ant == ants {

				if !moved {
					return
				}

				// Check room capacity.
				occupied := make(map[byte]bool)

				for _, room := range next {

					if room != 0 && room != 1 {

						if occupied[room] {
							return
						}

						occupied[room] = true
					}
				}

				newState := string(next)

				if _, exists := distance[newState]; !exists {
					distance[newState] = turn + 1
					queue = append(queue, newState)
				}

				return
			}

			current := int(state[ant])

			// Option 1: Stay in the current room.
			next[ant] = byte(current)
			explore(ant+1, moved)

			// Finished ants cannot move.
			if current == 1 {
				return
			}

			// Option 2: Move through a tunnel.
			for _, destination := range adj[current] {

				// Returning to the start is unnecessary.
				if destination == 0 {
					continue
				}

				a, b := current, destination

				if a > b {
					a, b = b, a
				}

				tunnel := [2]int{a, b}

				if usedTunnels[tunnel] {
					continue
				}

				usedTunnels[tunnel] = true
				next[ant] = byte(destination)

				explore(ant+1, true)

				delete(usedTunnels, tunnel)
			}
		}

		explore(0, false)
	}

	return -1
}

func TestOptimalSmallGraphs(t *testing.T) {

	rng := rand.New(rand.NewSource(5555))

	for trial := 0; trial < 100; trial++ {

		rooms := []string{"S", "E"}

		for i := 0; i < 4; i++ {
			rooms = append(rooms, fmt.Sprintf("R%d", i))
		}

		graph := make(map[string][]string)

		for _, room := range rooms {
			graph[room] = nil
		}

		// Generate random bidirectional tunnels.
		for i := 0; i < len(rooms); i++ {

			for j := i + 1; j < len(rooms); j++ {

				if rng.Intn(100) < 40 {

					a := rooms[i]
					b := rooms[j]

					graph[a] = append(graph[a], b)
					graph[b] = append(graph[b], a)
				}
			}
		}

		ants := 1 + rng.Intn(3)

		// Independent exhaustive search.
		optimal := shortestAntSchedule(graph, ants)

		// Our lem-in solver.
		paths := findOptimizedPaths(graph, "S", "E", ants)

		// Handle disconnected graphs.
		if optimal == -1 {

			if len(paths) != 0 {
				t.Fatal("expected no path")
			}

			continue
		}

		if len(paths) == 0 {
			t.Fatal("solver missed a valid path")
		}

		actual := int(estimatedTurns(paths, ants))

		if actual != optimal {

			t.Fatalf(
				"trial %d: expected %d turns, got %d\n"+
					"graph: %v\npaths: %v",
				trial, optimal, actual, graph, paths,
			)
		}
	}
}
