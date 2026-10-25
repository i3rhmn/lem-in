package lemin

import (
	"fmt"
	"strings"
)

// Resolve hyphenated room names only when the tunnel is unambiguous.
func resolveTunnel(line string, rooms map[string]Room) (string, string, error) {
	var a, b string
	matches := 0
	for i := 1; i < len(line)-1; i++ {
		if line[i] != '-' {
			continue
		}
		left, right := line[:i], line[i+1:]
		if _, ok := rooms[left]; !ok {
			continue
		}
		if _, ok := rooms[right]; !ok {
			continue
		}
		a, b = left, right
		matches++
	}
	if matches > 1 {
		return "", "", fmt.Errorf("ambiguous tunnel %q; room names allow multiple interpretations", line)
	}
	if matches == 0 {
		if strings.Count(line, "-") == 1 {
			parts := strings.SplitN(line, "-", 2)
			if parts[0] == "" || parts[1] == "" {
				return "", "", fmt.Errorf("invalid tunnel %q; both room names are required", line)
			}
			for _, name := range parts {
				if _, ok := rooms[name]; !ok {
					return "", "", fmt.Errorf("unknown room %q in tunnel %q", name, line)
				}
			}
		}
		if !strings.Contains(line, "-") {
			return "", "", fmt.Errorf("invalid tunnel %q; expected room1-room2", line)
		}
		return "", "", fmt.Errorf("invalid tunnel %q; cannot resolve two existing room names", line)
	}
	if a == b {
		return "", "", fmt.Errorf("self-link %q; a room cannot connect to itself", line)
	}
	return a, b, nil
}

func addTunnel(graph map[string][]string, seen map[[2]string]bool, a, b string) error {
	if a > b {
		a, b = b, a
	}
	key := [2]string{a, b}
	if seen[key] {
		return fmt.Errorf("duplicate tunnel between %q and %q; reversed links count as duplicates", a, b)
	}
	seen[key] = true
	graph[a] = append(graph[a], b)
	graph[b] = append(graph[b], a)
	return nil
}
