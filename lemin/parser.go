package lemin

import (
	"errors"
	"fmt"
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

func parseColony(data []byte) (*Colony, error) {
	if strings.TrimSpace(string(data)) == "" {
		return nil, errors.New("empty input file")
	}
	lines := strings.Split(string(data), "\n")
	ants, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if errors.Is(err, strconv.ErrRange) {
		return nil, errors.New("line 1: ant count exceeds the supported integer range")
	}
	if err != nil {
		return nil, errors.New("line 1: invalid ant count; expected a positive integer")
	}
	if ants <= 0 {
		return nil, errors.New("line 1: ant count must be greater than zero")
	}

	colony := &Colony{Ants: ants, Rooms: make(map[string]Room), Graph: make(map[string][]string)}
	seenTunnels := make(map[[2]string]bool)
	pending := ""
	readingLinks := false
	linkCount := 0

	for index, raw := range lines[1:] {
		lineNumber := index + 2
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if line == "##start" || line == "##end" {
			if readingLinks {
				return nil, fmt.Errorf("line %d: %s must appear before tunnels", lineNumber, line)
			}
			if pending != "" {
				return nil, fmt.Errorf("line %d: expected a room after ##%s, found %s", lineNumber, pending, line)
			}
			if line == "##start" {
				if colony.Start != "" {
					return nil, fmt.Errorf("line %d: duplicate ##start command", lineNumber)
				}
				pending = "start"
			} else {
				if colony.End != "" {
					return nil, fmt.Errorf("line %d: duplicate ##end command", lineNumber)
				}
				pending = "end"
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if !readingLinks && len(parts) == 3 {
			name := parts[0]
			if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
				return nil, fmt.Errorf("line %d: invalid room name %q; names cannot start with L or #", lineNumber, name)
			}
			if _, exists := colony.Rooms[name]; exists {
				return nil, fmt.Errorf("line %d: duplicate room %q", lineNumber, name)
			}
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				return nil, fmt.Errorf("line %d: invalid coordinates for room %q; X and Y must fit in an int", lineNumber, name)
			}
			colony.Rooms[name] = Room{Name: name, X: x, Y: y}
			colony.Graph[name] = []string{}
			if pending == "start" {
				colony.Start = name
			} else if pending == "end" {
				colony.End = name
			}
			pending = ""
			continue
		}
		if pending != "" {
			return nil, fmt.Errorf("line %d: expected a room after ##%s; format: name X Y", lineNumber, pending)
		}
		if readingLinks && len(parts) == 3 {
			return nil, fmt.Errorf("line %d: room definition after tunnels; define all rooms first", lineNumber)
		}
		if len(parts) != 1 {
			return nil, fmt.Errorf("line %d: malformed room or tunnel; use name X Y or room1-room2", lineNumber)
		}
		a, b, err := resolveTunnel(line, colony.Rooms)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		if err := addTunnel(colony.Graph, seenTunnels, a, b); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		readingLinks = true
		linkCount++
	}
	if pending != "" {
		return nil, fmt.Errorf("missing room definition after ##%s", pending)
	}
	if colony.Start == "" {
		return nil, errors.New("no start room found; add ##start before its room")
	}
	if colony.End == "" {
		return nil, errors.New("no end room found; add ##end before its room")
	}
	if linkCount == 0 {
		return nil, errors.New("no tunnels found; add room1-room2 links")
	}
	return colony, nil
}
