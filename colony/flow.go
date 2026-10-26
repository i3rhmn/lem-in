package colony

import (
	"slices"
	"sort"
)

type flowEdge struct {
	to, reverse, capacity, cost int
	fromRoom, toRoom            string
}

type flowNetwork struct {
	edges [][]*flowEdge
}

// Each forward edge has a reverse edge that can undo a previous choice.
func (f *flowNetwork) addEdge(from, to, cost int, fromRoom, toRoom string) {
	forward := &flowEdge{
		to: to, reverse: len(f.edges[to]),
		capacity: 1, cost: cost, fromRoom: fromRoom, toRoom: toRoom,
	}
	backward := &flowEdge{
		to: from, reverse: len(f.edges[from]),
		capacity: 0, cost: -cost,
	}
	f.edges[from] = append(f.edges[from], forward)
	f.edges[to] = append(f.edges[to], backward)
}

// Find another route through the residual graph and update its flow.
func (f *flowNetwork) augment(start, end int) bool {
	const infinity = int(^uint(0)>>1) / 4
	dist := make([]int, len(f.edges))
	parent := make([]int, len(f.edges))
	parentEdge := make([]int, len(f.edges))
	inQueue := make([]bool, len(f.edges))

	for i := range dist {
		dist[i] = infinity
		parent[i] = -1
	}
	dist[start] = 0
	queue := []int{start}
	inQueue[start] = true

	// Queue-based shortest path handles negative-cost reverse edges.
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		inQueue[current] = false

		for i, edge := range f.edges[current] {
			if edge.capacity == 0 || dist[current]+edge.cost >= dist[edge.to] {
				continue
			}
			dist[edge.to] = dist[current] + edge.cost
			parent[edge.to] = current
			parentEdge[edge.to] = i

			if !inQueue[edge.to] {
				queue = append(queue, edge.to)
				inQueue[edge.to] = true
			}
		}
	}
	if parent[end] == -1 {
		return false
	}

	for node := end; node != start; node = parent[node] {
		from := parent[node]
		edge := f.edges[from][parentEdge[node]]
		edge.capacity--
		f.edges[node][edge.reverse].capacity++
	}
	return true
}

// Build ordinary room paths from the currently selected flow edges.
func (f *flowNetwork) usedPaths(start, end string) [][]string {
	next := make(map[string][]string)

	for _, edges := range f.edges {
		for _, edge := range edges {
			if edge.fromRoom != "" && edge.capacity == 0 {
				next[edge.fromRoom] = append(next[edge.fromRoom], edge.toRoom)
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

	// Shorter paths first; alphabetical ordering makes ties reproducible.
	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}
		return slices.Compare(paths[i], paths[j]) < 0
	})
	return paths
}

// Binary-search the minimum turns needed for the given paths and ants.
func estimatedTurns(paths [][]string, ants int) int64 {
	low, high := int64(0), int64(ants)+int64(len(paths[0]))-2

	for low < high {
		mid := low + (high-low)/2
		capacity := int64(0)

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

// Use minimum-cost flow to choose among different path counts.
func findOptimizedPaths(graph map[string][]string, start, end string, ants int) [][]string {
	rooms := make([]string, 0, len(graph))
	for room := range graph {
		rooms = append(rooms, room)
	}
	sort.Strings(rooms)

	// Split intermediate rooms into an in-node and out-node (capacity 1).
	in, out := make(map[string]int), make(map[string]int)
	nodes := 0

	for _, room := range rooms {
		in[room] = nodes
		nodes++
		out[room] = in[room]

		if room != start && room != end {
			out[room] = nodes
			nodes++
		}
	}

	network := &flowNetwork{edges: make([][]*flowEdge, nodes)}

	for _, room := range rooms {
		if room != start && room != end {
			network.addEdge(in[room], out[room], 0, "", "")
		}

		neighbors := append([]string(nil), graph[room]...)
		sort.Strings(neighbors)

		for _, neighbor := range neighbors {
			if room != end && neighbor != start {
				network.addEdge(out[room], in[neighbor], 1, room, neighbor)
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

	var bestPaths [][]string
	var bestTurns int64

	for count := 0; count < limit; count++ {
		if !network.augment(out[start], in[end]) {
			break
		}

		paths := network.usedPaths(start, end)
		if len(paths) != count+1 {
			continue
		}

		turns := estimatedTurns(paths, ants)
		if bestPaths == nil || turns < bestTurns {
			bestPaths, bestTurns = paths, turns
		}
	}
	return bestPaths
}
