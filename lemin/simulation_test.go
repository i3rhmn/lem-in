package lemin

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Capture the output of our simulation.
func captureSimulation(paths [][]string, assignments [][]int) (string, error) {
	old := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = writer
	simulateAnts(paths, assignments)
	writer.Close()
	os.Stdout = old

	output, err := io.ReadAll(reader)
	reader.Close()

	return string(output), err
}

// Check whether every ant movement follows the rules.
func checkSimulation(
	t *testing.T,
	graph map[string][]string,
	start, end string,
	ants int,
	output string,
) int {

	t.Helper()

	// Initially, every ant is in the starting room.
	positions := make([]string, ants+1)

	for ant := 1; ant <= ants; ant++ {
		positions[ant] = start
	}

	turns := strings.Split(strings.TrimSpace(output), "\n")

	for turnIndex, line := range turns {

		if line == "" {
			t.Fatalf("empty turn at %d", turnIndex+1)
		}

		movedAnts := make(map[int]bool)
		usedTunnels := make(map[[2]string]bool)

		// Check every movement in the current turn.
		for _, move := range strings.Fields(line) {

			if !strings.HasPrefix(move, "L") {
				t.Fatalf("invalid movement: %s", move)
			}

			parts := strings.SplitN(move[1:], "-", 2)

			if len(parts) != 2 {
				t.Fatalf("invalid movement: %s", move)
			}

			ant, err := strconv.Atoi(parts[0])

			if err != nil || ant < 1 || ant > ants {
				t.Fatalf("invalid ant: %s", move)
			}

			// An ant cannot move twice in one turn.
			if movedAnts[ant] {
				t.Fatalf("ant %d moved twice", ant)
			}

			movedAnts[ant] = true

			from := positions[ant]
			to := parts[1]

			// Finished ants cannot move again.
			if from == end {
				t.Fatalf("ant %d moved after finishing", ant)
			}

			// Verify that the tunnel exists.
			linked := false

			for _, next := range graph[from] {
				if next == to {
					linked = true
					break
				}
			}

			if !linked {
				t.Fatalf("invalid tunnel: %s -> %s", from, to)
			}

			// Check that the tunnel is not reused this turn.
			a, b := from, to

			if a > b {
				a, b = b, a
			}

			tunnel := [2]string{a, b}

			if usedTunnels[tunnel] {
				t.Fatalf("tunnel %s-%s used twice", a, b)
			}

			usedTunnels[tunnel] = true

			// Update the ant's position.
			positions[ant] = to
		}

		// Check room occupancy after simultaneous movement.
		occupied := make(map[string]int)

		for ant := 1; ant <= ants; ant++ {

			room := positions[ant]

			if room != start && room != end {

				occupied[room]++

				if occupied[room] > 1 {
					t.Fatalf("multiple ants in room %s", room)
				}
			}
		}
	}

	// Make sure every ant reached the destination.
	for ant := 1; ant <= ants; ant++ {

		if positions[ant] != end {
			t.Fatalf("ant %d did not reach the end", ant)
		}
	}

	return len(turns)
}

// Test different types of colonies.
func TestCompleteSimulations(t *testing.T) {

	tests := []struct {
		name          string
		input         string
		expectedTurns int
	}{
		{
			"direct tunnel",
			"5\n##start\nS 0 0\n##end\nE 1 1\nS-E\n",
			5,
		},
		{
			"single path",
			"3\n##start\nS 0 0\nA 1 0\nB 2 0\n" +
				"##end\nE 3 0\nS-A\nA-B\nB-E\n",
			5,
		},
		{
			"two routes",
			"3\n##start\nS 0 0\nA 1 1\nB 1 2\n" +
				"##end\nE 2 0\nS-A\nS-B\nA-E\nB-E\n",
			3,
		},
		{
			"greedy trap",
			"4\n##start\nS 0 0\nA 1 1\nB 1 2\n" +
				"X 2 1\nY 2 2\n##end\nE 3 0\n" +
				"S-A\nS-B\nA-X\nA-Y\nB-X\nX-E\nY-E\n",
			4,
		},
		{
			"hyphen room",
			"2\n##start\nS 0 0\nA-1 1 0\n" +
				"##end\nE 2 0\nS-A-1\nA-1-E\n",
			3,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			colony, err := parseColony([]byte(tt.input))
			if err != nil {
				t.Fatal(err)
			}

			paths := findOptimizedPaths(
				colony.Graph,
				colony.Start,
				colony.End,
				colony.Ants,
			)

			if len(paths) == 0 {
				t.Fatal("no paths found")
			}

			assignments := assignAnts(colony.Ants, paths)

			output, err := captureSimulation(paths, assignments)
			if err != nil {
				t.Fatal(err)
			}

			turns := checkSimulation(
				t,
				colony.Graph,
				colony.Start,
				colony.End,
				colony.Ants,
				output,
			)

			if turns != tt.expectedTurns {
				t.Fatalf(
					"expected %d turns, got %d\n%s",
					tt.expectedTurns, turns, output,
				)
			}

			if int64(turns) != estimatedTurns(paths, colony.Ants) {
				t.Fatal("simulation and estimator disagree")
			}
		})
	}
}

// A disconnected colony must not produce paths.
func TestDisconnectedColony(t *testing.T) {

	input := "2\n##start\nS 0 0\nA 1 0\n" +
		"##end\nE 2 0\nB 3 0\nS-A\nB-E\n"

	colony, err := parseColony([]byte(input))
	if err != nil {
		t.Fatal(err)
	}

	paths := findOptimizedPaths(
		colony.Graph,
		colony.Start,
		colony.End,
		colony.Ants,
	)

	if len(paths) != 0 {
		t.Fatalf("expected no route, got %v", paths)
	}
}

// Generate random graphs and verify their simulations.
func TestRandomGraphSimulations(t *testing.T) {

	rng := rand.New(rand.NewSource(42))

	for trial := 0; trial < 80; trial++ {

		rooms := []string{"S", "E"}

		for i := 0; i < 7; i++ {
			rooms = append(rooms, fmt.Sprintf("R%d", i))
		}

		graph := make(map[string][]string)

		for _, room := range rooms {
			graph[room] = nil
		}

		added := make(map[[2]string]bool)

		// Add a bidirectional tunnel if it does not exist.
		add := func(a, b string) {

			key := [2]string{a, b}

			if a > b {
				key = [2]string{b, a}
			}

			if added[key] {
				return
			}

			added[key] = true

			graph[a] = append(graph[a], b)
			graph[b] = append(graph[b], a)
		}

		// Guarantee at least one valid path.
		add("S", "R0")
		add("R0", "E")

		// Add random tunnels.
		for i := range rooms {
			for j := i + 1; j < len(rooms); j++ {

				if rng.Intn(4) == 0 {
					add(rooms[i], rooms[j])
				}
			}
		}

		ants := 1 + rng.Intn(10)

		paths := findOptimizedPaths(graph, "S", "E", ants)

		if len(paths) == 0 {
			t.Fatalf("trial %d: no path found", trial)
		}

		assignments := assignAnts(ants, paths)

		output, err := captureSimulation(paths, assignments)

		if err != nil {
			t.Fatal(err)
		}

		turns := checkSimulation(
			t, graph, "S", "E", ants, output,
		)

		if int64(turns) != estimatedTurns(paths, ants) {
			t.Fatalf("trial %d: incorrect turn count", trial)
		}
	}
}
