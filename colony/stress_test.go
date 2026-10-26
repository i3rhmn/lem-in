package colony

import (
	"fmt"
	"strings"
	"testing"
)

// Create a colony with one long path.
func makeCorridor(rooms, ants int) string {
	var b strings.Builder

	fmt.Fprintf(&b, "%d\n##start\nS 0 0\n", ants)

	for i := 1; i <= rooms; i++ {
		fmt.Fprintf(&b, "R%d %d 0\n", i, i)
	}

	fmt.Fprintf(&b, "##end\nE %d 0\n", rooms+1)

	fmt.Fprintln(&b, "S-R1")

	for i := 1; i < rooms; i++ {
		fmt.Fprintf(&b, "R%d-R%d\n", i, i+1)
	}

	fmt.Fprintf(&b, "R%d-E\n", rooms)

	return b.String()
}

// Create multiple independent paths.
func makeParallelPaths(lanes, ants int) string {
	var b strings.Builder

	fmt.Fprintf(&b, "%d\n##start\nS 0 0\n", ants)

	for i := 0; i < lanes; i++ {
		fmt.Fprintf(&b, "R%d 1 %d\n", i, i)
	}

	fmt.Fprintln(&b, "##end")
	fmt.Fprintln(&b, "E 2 0")

	for i := 0; i < lanes; i++ {
		fmt.Fprintf(&b, "S-R%d\n", i)
		fmt.Fprintf(&b, "R%d-E\n", i)
	}

	return b.String()
}

// Test 200 intermediate rooms with 30 ants.
func TestLargeCorridor(t *testing.T) {
	const rooms = 200
	const ants = 30

	colony, err := parseColony(
		[]byte(makeCorridor(rooms, ants)),
	)
	if err != nil {
		t.Fatal(err)
	}

	paths := findOptimizedPaths(
		colony.Graph,
		colony.Start,
		colony.End,
		colony.Ants,
	)

	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}

	expected := int64(230)
	actual := estimatedTurns(paths, ants)

	if actual != expected {
		t.Fatalf("expected %d turns, got %d", expected, actual)
	}
}

// Test 30 parallel paths with 300 ants.
func TestManyParallelPaths(t *testing.T) {
	const lanes = 30
	const ants = 300

	colony, err := parseColony(
		[]byte(makeParallelPaths(lanes, ants)),
	)
	if err != nil {
		t.Fatal(err)
	}

	paths := findOptimizedPaths(
		colony.Graph,
		colony.Start,
		colony.End,
		colony.Ants,
	)

	if len(paths) != lanes {
		t.Fatalf("expected %d paths, got %d", lanes, len(paths))
	}

	expected := int64(11)
	actual := estimatedTurns(paths, ants)

	if actual != expected {
		t.Fatalf("expected %d turns, got %d", expected, actual)
	}
}

// Test a direct tunnel with an alternative path.
func TestDirectTunnelWithAlternative(t *testing.T) {
	input := `5
##start
S 0 0
A 1 1
##end
E 2 0
S-E
S-A
A-E
`

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

	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %v", paths)
	}

	if estimatedTurns(paths, 5) != 3 {
		t.Fatal("expected completion in 3 turns")
	}
}

// Measure solver performance.
func BenchmarkOptimizedPaths(b *testing.B) {
	colony, err := parseColony(
		[]byte(makeParallelPaths(30, 300)),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		paths := findOptimizedPaths(
			colony.Graph,
			colony.Start,
			colony.End,
			colony.Ants,
		)

		if len(paths) == 0 {
			b.Fatal("no paths found")
		}
	}
}
