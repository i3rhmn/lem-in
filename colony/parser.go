package colony

import (
	"errors"
	"strconv"
	"strings"
)

// A graph is an adjacency list: room name -> neighboring room names.
type Room struct {
	Name string
	X, Y int
}

type Colony struct {
	Ants       int
	Rooms      map[string]Room
	Start, End string
	Graph      map[string][]string
}

func parseColony(data []byte) (*Colony, error) {
	lines := strings.Split(string(data), "\n")
	ants, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil || ants <= 0 {
		return nil, errors.New("invalid ant count")
	}
	farm := &Colony{Ants: ants, Rooms: make(map[string]Room), Graph: make(map[string][]string)}
	seen := make(map[[2]string]bool)
	pending, links, linkCount := "", false, 0

	for _, raw := range lines[1:] {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if line == "##start" || line == "##end" {
			if pending != "" || links {
				return nil, errors.New("invalid command")
			}
			if line == "##start" {
				if farm.Start != "" {
					return nil, errors.New("duplicate start")
				}
				pending = "start"
			} else {
				if farm.End != "" {
					return nil, errors.New("duplicate end")
				}
				pending = "end"
			}
			continue
		}
		if strings.HasPrefix(line, "#") { // Comments and unknown commands.
			continue
		}

		parts := strings.Fields(line)
		if !links && len(parts) == 3 {
			name := parts[0]
			if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
				return nil, errors.New("invalid room name")
			}
			if _, exists := farm.Rooms[name]; exists {
				return nil, errors.New("duplicate room")
			}
			x, xerr := strconv.Atoi(parts[1])
			y, yerr := strconv.Atoi(parts[2])
			if xerr != nil || yerr != nil {
				return nil, errors.New("invalid coordinates")
			}
			farm.Rooms[name] = Room{name, x, y}
			farm.Graph[name] = []string{}
			if pending == "start" {
				farm.Start = name
			}
			if pending == "end" {
				farm.End = name
			}
			pending = ""
			continue
		}
		if pending != "" || len(parts) != 1 {
			return nil, errors.New("invalid room or tunnel")
		}
		// A tunnel splits at a '-' that separates two known room names.
		a, b, matches := "", "", 0
		for i := 1; i < len(line)-1; i++ {
			if line[i] != '-' {
				continue
			}
			left, right := line[:i], line[i+1:]
			_, leftOK := farm.Rooms[left]
			_, rightOK := farm.Rooms[right]
			if leftOK && rightOK {
				a, b, matches = left, right, matches+1
			}
		}
		if matches != 1 || a == b {
			return nil, errors.New("invalid tunnel")
		}
		if a > b {
			a, b = b, a
		}
		key := [2]string{a, b}
		if seen[key] {
			return nil, errors.New("duplicate tunnel")
		}
		seen[key] = true
		farm.Graph[a] = append(farm.Graph[a], b)
		farm.Graph[b] = append(farm.Graph[b], a)
		links, linkCount = true, linkCount+1
	}
	if pending != "" || farm.Start == "" || farm.End == "" || linkCount == 0 {
		return nil, errors.New("missing start, end, or tunnels")
	}
	return farm, nil
}
