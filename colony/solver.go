package colony

// assignAnts distributes ants to minimize their estimated arrival time.
func assignAnts(antCount int, paths [][]string) [][]int {
	assignments := make([][]int, len(paths))

	for ant := 1; ant <= antCount; ant++ {
		best := 0

		for i := 1; i < len(paths); i++ {
			arrival := len(paths[i]) + len(assignments[i])
			bestArrival := len(paths[best]) + len(assignments[best])

			if arrival < bestArrival {
				best = i
			}
		}

		assignments[best] = append(assignments[best], ant)
	}

	return assignments
}
