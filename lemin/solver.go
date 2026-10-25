package lemin

import "errors"

// findShortestPath uses BFS to find a path
// from the start room to the end room.
func findShortestPath(graph map[string][]string, start, end string) ([]string, error) {

	// Queue stores rooms waiting to be explored
	queue := []string{start}

	// Track visited rooms to avoid infinite loops
	visited := make(map[string]bool)
	visited[start] = true

	// Remember where we came from
	parent := make(map[string]string)

	// Process the queue
	for len(queue) > 0 {

		// Take the first room
		current := queue[0]
		queue = queue[1:]

		// Stop when we reach the end
		if current == end {
			break
		}

		// Explore connected rooms
		for _, next := range graph[current] {

			if !visited[next] {

				visited[next] = true
				parent[next] = current

				queue = append(queue, next)
			}
		}
	}

	// Check whether the end was reached
	if !visited[end] {
		return nil, errors.New("no path found")
	}

	// Reconstruct the path backwards
	path := []string{}

	for room := end; room != start; room = parent[room] {
		path = append([]string{room}, path...)
	}

	// Add the starting room
	path = append([]string{start}, path...)

	return path, nil
}

// Find multiple paths without sharing intermediate rooms.
func findDisjointPaths(graph map[string][]string, start, end string) [][]string {

	paths := [][]string{}
	usedRooms := make(map[string]bool)
	directUsed := false

	for {
		// Create a graph without already-used rooms.
		available := make(map[string][]string)

		for room, neighbors := range graph {

			if usedRooms[room] {
				continue
			}

			for _, next := range neighbors {

				if usedRooms[next] {
					continue
				}

				// Don't use the direct start-end tunnel twice.
				if directUsed && ((room == start && next == end) ||
					(room == end && next == start)) {
					continue
				}

				available[room] = append(available[room], next)
			}
		}

		// Use our existing BFS function.
		path, err := findShortestPath(available, start, end)

		if err != nil {
			break
		}

		paths = append(paths, path)

		if len(path) == 2 {
			directUsed = true
		}

		// Prevent other paths from using these intermediate rooms.
		for _, room := range path[1 : len(path)-1] {
			usedRooms[room] = true
		}
	}

	return paths
}

// Assign ants to paths based on earliest arrival.
func assignAnts(antCount int, paths [][]string) [][]int {

	assignments := make([][]int, len(paths))

	for ant := 1; ant <= antCount; ant++ {

		best := 0

		for i := 1; i < len(paths); i++ {

			arrival := len(paths[i]) - 1 + len(assignments[i])

			bestArrival := len(paths[best]) - 1 + len(assignments[best])

			if arrival < bestArrival {
				best = i
			}
		}

		assignments[best] = append(assignments[best], ant)
	}

	return assignments
}
