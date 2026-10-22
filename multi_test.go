package main

import (
	"reflect"
	"testing"
)

func TestMultiplePaths(t *testing.T) {

	graph := map[string][]string{
		"A": {"B", "C"},
		"B": {"A", "D"},
		"C": {"A", "D"},
		"D": {"B", "C"},
	}

	paths := findDisjointPaths(graph, "A", "D")

	if len(paths) != 2 {
		t.Fatalf("Expected 2 paths, got %d", len(paths))
	}

	assignments := assignAnts(3, paths)

	expected := [][]int{
		{1, 3},
		{2},
	}

	if !reflect.DeepEqual(assignments, expected) {
		t.Fatalf("Expected %v, got %v", expected, assignments)
	}
}
