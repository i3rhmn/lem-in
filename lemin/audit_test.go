package lemin

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAuditCLI(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "lem-in")
	build := exec.Command("go", "build", "-o", binary, "..")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, output)
	}

	run := func(t *testing.T, args ...string) string {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		output, err := exec.CommandContext(ctx, binary, args...).CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("execution timed out: %v", ctx.Err())
		}
		if err != nil {
			t.Fatalf("execution failed: %v\n%s", err, output)
		}
		return string(output)
	}

	valid := []struct {
		name  string
		turns int
	}{
		{"example00", 6},
		{"example01", 8},
		{"example02", 11},
		{"example03", 6},
		{"example04", 6},
		{"example05", 8},
	}

	for _, tc := range valid {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("..", "testdata", "audit", tc.name+".txt")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			colony, err := parseColony(data)
			if err != nil {
				t.Fatal(err)
			}

			output := run(t, path)
			prefix := string(data) + "\n"
			if !strings.HasPrefix(output, prefix) {
				t.Fatalf("original input or blank separator is incorrect:\n%s", output)
			}

			moves := strings.TrimPrefix(output, prefix)
			turns := checkSimulation(
				t, colony.Graph, colony.Start, colony.End,
				colony.Ants, moves,
			)

			if turns > tc.turns {
				t.Fatalf("got %d turns; audit limit is %d", turns, tc.turns)
			}

			for repeat := 0; repeat < 3; repeat++ {
				if again := run(t, path); again != output {
					t.Fatal("output changed between identical runs")
				}
			}

			t.Logf("PASS: %d turns, limit %d; valid movements and repeatable output",
				turns, tc.turns)
		})
	}

	validInput := "3\n##start\nS 0 0\n##end\nE 1 1\nS-E\n"

	t.Run("unknown_commands_and_comments", func(t *testing.T) {
		input := strings.Replace(
			validInput, "##start\n",
			"#comment\n##unknown\n##start\n", 1,
		)
		path := filepath.Join(t.TempDir(), "unknown-command.txt")
		if err := os.WriteFile(path, []byte(input), 0600); err != nil {
			t.Fatal(err)
		}

		output := run(t, path)
		prefix := input + "\n"
		if !strings.HasPrefix(output, prefix) {
			t.Fatalf("unknown command was not ignored correctly:\n%s", output)
		}

		graph := map[string][]string{"S": {"E"}, "E": {"S"}}
		turns := checkSimulation(t, graph, "S", "E", 3,
			strings.TrimPrefix(output, prefix))
		if turns != 3 {
			t.Fatalf("direct tunnel: expected 3 turns, got %d", turns)
		}
	})

	invalid := []struct {
		name  string
		input string
	}{
		{"empty_file", ""},
		{"zero_ants", "0\n"},
		{"negative_ants", "-3\n"},
		{"nonnumeric_ants", "ants\n"},
		{"overflowing_ants", "999999999999999999999999999999\n"},
		{"missing_start", "3\nS 0 0\n##end\nE 1 1\nS-E\n"},
		{"missing_end", "3\n##start\nS 0 0\nE 1 1\nS-E\n"},
		{"duplicate_start", "3\n##start\nS 0 0\n##start\nA 2 2\n##end\nE 1 1\nS-E\n"},
		{"duplicate_end", "3\n##start\nS 0 0\n##end\nE 1 1\n##end\nA 2 2\nS-E\n"},
		{"duplicate_room", "3\n##start\nS 0 0\nS 2 2\n##end\nE 1 1\nS-E\n"},
		{"invalid_coordinates", "3\n##start\nS abc 0\n##end\nE 1 1\nS-E\n"},
		{"invalid_name", "3\n##start\nLroom 0 0\n##end\nE 1 1\nLroom-E\n"},
		{"unknown_room", "3\n##start\nS 0 0\n##end\nE 1 1\nS-X\n"},
		{"self_link", "3\n##start\nS 0 0\n##end\nE 1 1\nS-S\n"},
		{"duplicate_link", validInput + "S-E\n"},
		{"reversed_duplicate_link", validInput + "E-S\n"},
		{"room_after_links", validInput + "A 2 2\n"},
		{"no_links", "3\n##start\nS 0 0\n##end\nE 1 1\n"},
		{"no_route", "3\n##start\nS 0 0\nA 1 0\n##end\nE 2 0\nS-A\n"},
		{"unfinished_command", "3\n##start\n"},
	}

	wantReasons := map[string]string{
		"empty_file":              "empty input file",
		"zero_ants":               "ant count must be greater than zero",
		"negative_ants":           "ant count must be greater than zero",
		"nonnumeric_ants":         "invalid ant count; expected a positive integer",
		"overflowing_ants":        "ant count exceeds the supported integer range",
		"missing_start":           "no start room found",
		"missing_end":             "no end room found",
		"duplicate_start":         "duplicate ##start command",
		"duplicate_end":           "duplicate ##end command",
		"duplicate_room":          "duplicate room \"S\"",
		"invalid_coordinates":     "invalid coordinates for room \"S\"",
		"invalid_name":            "invalid room name \"Lroom\"",
		"unknown_room":            "unknown room \"X\"",
		"self_link":               "self-link \"S-S\"",
		"duplicate_link":          "duplicate tunnel between \"E\" and \"S\"",
		"reversed_duplicate_link": "duplicate tunnel between \"E\" and \"S\"",
		"room_after_links":        "room definition after tunnels",
		"no_links":                "no tunnels found",
		"no_route":                "no path from start room \"S\" to end room \"E\"",
		"unfinished_command":      "missing room definition after ##start",
	}

	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "invalid.txt")
			if err := os.WriteFile(path, []byte(tc.input), 0600); err != nil {
				t.Fatal(err)
			}
			output := run(t, path)
			if !strings.HasPrefix(output, "ERROR: invalid data format, ") || !strings.Contains(output, wantReasons[tc.name]) {
				t.Fatalf("unexpected error output: %q", output)
			}
		})
	}

	t.Run("missing_argument", func(t *testing.T) {
		if output := run(t); strings.TrimSpace(output) != "ERROR: invalid data format, missing file argument; usage: go run . <file.txt>" {
			t.Fatalf("unexpected output: %q", output)
		}
	})

	t.Run("extra_arguments", func(t *testing.T) {
		if output := run(t, "one.txt", "two.txt"); strings.TrimSpace(output) != "ERROR: invalid data format, too many arguments; provide exactly one file" {
			t.Fatalf("unexpected output: %q", output)
		}
	})

	t.Run("missing_file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "does-not-exist.txt")
		if output := run(t, path); !strings.HasPrefix(output, "ERROR: invalid data format, file not found: ") || !strings.Contains(output, path) {
			t.Fatalf("unexpected output: %q", output)
		}
	})
}
