package main

import "errors"

// Identify two room names from a tunnel definition.
func resolveTunnel(line string, rooms map[string]Room) (string, string, error) {

	var a, b string
	matches := 0

	// Find a separator between two existing room names.
	for i := 1; i < len(line)-1; i++ {

		if line[i] != '-' {
			continue
		}

		left := line[:i]
		right := line[i+1:]

		if _, ok := rooms[left]; !ok {
			continue
		}

		if _, ok := rooms[right]; !ok {
			continue
		}

		a = left
		b = right
		matches++
	}

	// There must be exactly one valid interpretation.
	if matches != 1 {
		return "", "", errors.New("invalid tunnel")
	}

	// Prevent self-links.
	if a == b {
		return "", "", errors.New("self tunnel")
	}

	return a, b, nil
}

// Add a validated tunnel to the graph.
func addTunnel(
	graph map[string][]string,
	seen map[[2]string]bool,
	a, b string,
) error {

	// Normalize the room order.
	if a > b {
		a, b = b, a
	}

	key := [2]string{a, b}

	// Prevent duplicate tunnels.
	if seen[key] {
		return errors.New("duplicate tunnel")
	}

	seen[key] = true

	// Tunnels connect rooms in both directions.
	graph[a] = append(graph[a], b)
	graph[b] = append(graph[b], a)

	return nil
}
