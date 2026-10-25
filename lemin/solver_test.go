package lemin

import (
	"reflect"
	"testing"
)

func TestFindShortestPath(t *testing.T) {

	graph := map[string][]string{
		"A": {"B", "C"},
		"B": {"A", "D"},
		"C": {"A", "D"},
		"D": {"B", "C"},
	}

	path, err := findShortestPath(graph, "A", "D")

	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"A", "B", "D"}

	if !reflect.DeepEqual(path, expected) {
		t.Fatalf("Expected %v, got %v", expected, path)
	}

	t.Logf("Shortest path: %v", path)
}

func TestNoPath(t *testing.T) {

	graph := map[string][]string{
		"A": {"B"},
		"B": {"A"},
		"C": {},
	}

	_, err := findShortestPath(graph, "A", "C")

	if err == nil {
		t.Fatal("Expected an error because no path exists")
	}
}
