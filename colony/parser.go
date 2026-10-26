package colony

import (
	"errors"
	"strconv"
	"strings"
)

// Room represents a room in the colony.
type Room struct {
	Name string
	X    int
	Y    int
}

// Colony stores all information about the ant farm.
type Colony struct {
	Ants  int
	Rooms map[string]Room
	Start string
	End   string
	Graph map[string][]string
}

// parseColony reads and validates the entire input.
func parseColony(data []byte) (*Colony, error) {

	lines := strings.Split(string(data), "\n")

	// The first line must contain a positive number of ants.
	ants, err := strconv.Atoi(strings.TrimSpace(lines[0]))

	if err != nil || ants <= 0 {
		return nil, errors.New("invalid ant count")
	}

	colony := &Colony{
		Ants:  ants,
		Rooms: make(map[string]Room),
		Graph: make(map[string][]string),
	}

	seenTunnels := make(map[[2]string]bool)

	pending := ""
	readingLinks := false
	linkCount := 0

	// Skip the ant count and read the remaining lines.
	for _, raw := range lines[1:] {

		line := strings.TrimSpace(raw)

		if line == "" {
			continue
		}

		// Handle start and end commands.
		if line == "##start" || line == "##end" {

			if readingLinks || pending != "" {
				return nil, errors.New("invalid command position")
			}

			if line == "##start" {

				if colony.Start != "" {
					return nil, errors.New("duplicate start")
				}

				pending = "start"

			} else {

				if colony.End != "" {
					return nil, errors.New("duplicate end")
				}

				pending = "end"
			}

			continue
		}

		// Ignore comments and unknown commands.
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)

		// Read rooms before the tunnel section.
		if !readingLinks && len(parts) == 3 {

			name := parts[0]

			// Validate the room name.
			if name == "" ||
				strings.HasPrefix(name, "L") ||
				strings.HasPrefix(name, "#") {

				return nil, errors.New("invalid room name")
			}

			// Room names must be unique.
			if _, exists := colony.Rooms[name]; exists {
				return nil, errors.New("duplicate room")
			}

			// Validate coordinates.
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])

			if errX != nil || errY != nil {
				return nil, errors.New("invalid coordinates")
			}

			// Store the room.
			colony.Rooms[name] = Room{
				Name: name,
				X:    x,
				Y:    y,
			}

			colony.Graph[name] = []string{}

			// Assign the start or end room.
			if pending == "start" {
				colony.Start = name
			} else if pending == "end" {
				colony.End = name
			}

			pending = ""
			continue
		}

		// A command must be followed by a room.
		// Tunnel definitions must contain one token.
		if pending != "" || len(parts) != 1 {
			return nil, errors.New("invalid room or tunnel")
		}

		// Identify the two connected rooms.
		a, b, err := resolveTunnel(line, colony.Rooms)

		if err != nil {
			return nil, err
		}

		// Validate and store the tunnel.
		err = addTunnel(colony.Graph, seenTunnels, a, b)

		if err != nil {
			return nil, err
		}

		readingLinks = true
		linkCount++
	}

	// Check required colony information.
	if pending != "" ||
		colony.Start == "" ||
		colony.End == "" ||
		linkCount == 0 {

		return nil, errors.New("incomplete colony")
	}

	return colony, nil
}
