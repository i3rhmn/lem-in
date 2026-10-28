package colony

import (
	"sort"
	"strings"
)

// Extract chosen paths by following the used room-to-room edges.

// It converts the selected network edges into complete paths(actual room names).
func (f *network) paths(start, end string) [][]string {
	next := make(map[string][]string)
	for _, edges := range f.edges {
		for _, e := range edges {
			if e.fromRoom != "" && e.capacity == 0 {
				next[e.fromRoom] = append(next[e.fromRoom], e.toRoom)
			}
		}
	}
	var paths [][]string
	for _, first := range next[start] {
		path := []string{start, first}
		seen := map[string]bool{start: true, first: true}
		for current := first; current != end; {
			if len(next[current]) != 1 {
				return nil
			}
			current = next[current][0]
			if seen[current] {
				return nil
			}
			seen[current] = true
			path = append(path, current)
		}
		paths = append(paths, path)
	}
	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}
		return strings.Join(paths[i], "\x00") < strings.Join(paths[j], "\x00")
	})
	return paths
}

// paths = [][]string {
//	{"A", "B", "C"},
//		}

// How many turns are needed if the ants use these paths?
func estimatedTurns(paths [][]string, ants int) int64 {
	low, high := int64(0), int64(ants)+int64(len(paths[0]))-2
	for low < high {
		mid, capacity := low+(high-low)/2, int64(0)
		for _, path := range paths {
			if available := mid - int64(len(path)) + 2; available > 0 {
				capacity += available
			}
			if capacity >= int64(ants) {
				break
			}
		}
		if capacity >= int64(ants) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

//	// decide which path is better //
//
// Find several compatible paths; prefer the set with fewer total turns.
func findOptimizedPaths(graph map[string][]string, start, end string, ants int) [][]string {
	rooms := make([]string, 0, len(graph))
	for room := range graph {
		rooms = append(rooms, room)
	}
	sort.Strings(rooms)
	// rooms = []string{"A", "B", "C"} Ex

	// make a network
	in, out, count := make(map[string]int), make(map[string]int), 0
	for _, room := range rooms {
		in[room], out[room] = count, count
		count++
		if room != start && room != end {
			out[room] = count
			count++
		}
	}
	//	A         = node 0
	//	B IN      = node 1
	//	B OUT     = node 2
	//	C         = node 3

	f := &network{edges: make([][]*edge, count)}
	for _, room := range rooms {
		if room != start && room != end {
			f.add(in[room], out[room], 0, "", "") //graph.go
		}
		neighbors := append([]string(nil), graph[room]...)
		sort.Strings(neighbors)
		for _, to := range neighbors {
			if room != end && to != start {
				f.add(out[room], in[to], 1, room, to)
			}
		}
	}

	limit := ants
	if limit > len(graph[start]) {
		limit = len(graph[start])
	}
	if limit > len(graph[end]) {
		limit = len(graph[end])
	}
	var best [][]string
	var bestTurns int64

	for i := 0; i < limit && f.augment(out[start], in[end]); i++ { //graph.go
		paths := f.paths(start, end) //solver.go
		if len(paths) != i+1 {
			continue
		}
		turns := estimatedTurns(paths, ants) // solver.go
		if best == nil || turns < bestTurns {
			best, bestTurns = paths, turns
		}
	}
	return best
}

// best = paths
// bestTurns = 4