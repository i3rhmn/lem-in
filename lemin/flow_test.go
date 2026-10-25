package lemin

import "testing"

func TestFlowCanReroute(t *testing.T) {

	graph := map[string][]string{
		"S": {"A", "B"},
		"A": {"S", "X", "Y"},
		"B": {"S", "X"},
		"X": {"A", "B", "E"},
		"Y": {"A", "E"},
		"E": {"X", "Y"},
	}

	paths := findOptimizedPaths(graph, "S", "E", 4)

	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %v", paths)
	}

	if estimatedTurns(paths, 4) != 4 {
		t.Fatalf(
			"expected 4 turns, got %d",
			estimatedTurns(paths, 4),
		)
	}
}
