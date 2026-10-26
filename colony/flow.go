package colony

import (
	"sort"
	"strings"
)

type flowEdge struct {
	to       int
	reverse  int
	capacity int
	cost     int
	fromRoom string
	toRoom   string
}

type flowNetwork struct {
	edges [][]*flowEdge
}

// Add a forward edge and a backward edge
func (f *flowNetwork) addEdge(from, to, cost int, fromRoom, toRoom string) {

	forward := &flowEdge{
		to:       to,
		reverse:  len(f.edges[to]),
		capacity: 1,
		cost:     cost,
		fromRoom: fromRoom,
		toRoom:   toRoom,
	}

	backward := &flowEdge{
		to:       from,
		reverse:  len(f.edges[from]),
		capacity: 0,
		cost:     -cost,
	}

	f.edges[from] = append(f.edges[from], forward)
	f.edges[to] = append(f.edges[to], backward)
}

// Find another path, allowing previous paths to be adjusted
func (f *flowNetwork) augment(start, end int) bool {

	const infinity = int(^uint(0)>>1) / 4

	dist := make([]int, len(f.edges))
	previousNode := make([]int, len(f.edges))
	previousEdge := make([]int, len(f.edges))
	waiting := make([]bool, len(f.edges))

	for i := range dist {
		dist[i] = infinity
		previousNode[i] = -1
	}

	dist[start] = 0

	queue := []int{start}
	waiting[start] = true

	// Search the residual graph
	for head := 0; head < len(queue); head++ {

		current := queue[head]
		waiting[current] = false

		for i, edge := range f.edges[current] {

			if edge.capacity == 0 ||
				dist[current]+edge.cost >= dist[edge.to] {
				continue
			}

			dist[edge.to] = dist[current] + edge.cost
			previousNode[edge.to] = current
			previousEdge[edge.to] = i

			if !waiting[edge.to] {
				queue = append(queue, edge.to)
				waiting[edge.to] = true
			}
		}
	}

	if previousNode[end] == -1 {
		return false
	}

	// Update the flow along the selected path
	for node := end; node != start; node = previousNode[node] {

		from := previousNode[node]
		edge := f.edges[from][previousEdge[node]]

		edge.capacity--
		f.edges[node][edge.reverse].capacity++
	}

	return true
}

// Convert the current flow into actual room paths
func (f *flowNetwork) usedPaths(start, end string) [][]string {

	nextRooms := make(map[string][]string)

	for _, edges := range f.edges {

		for _, edge := range edges {

			if edge.fromRoom != "" && edge.capacity == 0 {
				nextRooms[edge.fromRoom] = append(
					nextRooms[edge.fromRoom],
					edge.toRoom,
				)
			}
		}
	}

	paths := [][]string{}

	for _, first := range nextRooms[start] {

		path := []string{start, first}

		seen := map[string]bool{
			start: true,
			first: true,
		}

		current := first

		for current != end {

			if len(nextRooms[current]) != 1 {
				return nil
			}

			current = nextRooms[current][0]

			if seen[current] {
				return nil
			}

			seen[current] = true
			path = append(path, current)
		}

		paths = append(paths, path)
	}

	// Sort paths by length
	sort.Slice(paths, func(i, j int) bool {

		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}

		return strings.Join(paths[i], "\x00") <
			strings.Join(paths[j], "\x00")
	})

	return paths
}

// Calculate the minimum turns for a set of paths
func estimatedTurns(paths [][]string, ants int) int64 {

	low := int64(0)
	high := int64(ants) + int64(len(paths[0])) - 2

	for low < high {

		mid := low + (high-low)/2
		places := int64(0)

		for _, path := range paths {

			available := mid - int64(len(path)-1) + 1

			if available > 0 {
				places += available
			}

			if places >= int64(ants) {
				break
			}
		}

		if places >= int64(ants) {
			high = mid
		} else {
			low = mid + 1
		}
	}

	return low
}

// Find optimized paths using minimum-cost flow
func findOptimizedPaths(
	graph map[string][]string,
	start, end string,
	ants int,
) [][]string {

	rooms := make([]string, 0, len(graph))

	for room := range graph {
		rooms = append(rooms, room)
	}

	sort.Strings(rooms)

	in := make(map[string]int)
	out := make(map[string]int)

	count := 0

	// Give each room graph node IDs
	for _, room := range rooms {

		in[room] = count
		count++

		if room != start && room != end {

			out[room] = count
			count++

		} else {
			out[room] = in[room]
		}
	}

	network := &flowNetwork{
		edges: make([][]*flowEdge, count),
	}

	// Build the flow network
	for _, room := range rooms {

		if room != start && room != end {
			network.addEdge(
				in[room],
				out[room],
				0,
				"",
				"",
			)
		}

		neighbors := append([]string(nil), graph[room]...)
		sort.Strings(neighbors)

		for _, neighbor := range neighbors {

			if room != end && neighbor != start {

				network.addEdge(
					out[room],
					in[neighbor],
					1,
					room,
					neighbor,
				)
			}
		}
	}

	// Determine the maximum paths to explore
	limit := ants

	if limit > len(graph[start]) {
		limit = len(graph[start])
	}

	if limit > len(graph[end]) {
		limit = len(graph[end])
	}

	var bestPaths [][]string
	var bestTurns int64

	// Try increasing numbers of paths
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
			bestPaths = paths
			bestTurns = turns
		}
	}

	return bestPaths
}
