package colony

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditExamples(t *testing.T) {
	// Build the actual program.
	binary := filepath.Join(t.TempDir(), "lem-in")
	build := exec.Command("go", "build", "-o", binary, "..")

	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, output)
	}

	examples := []struct {
		name     string
		maxTurns int
	}{
		{"example00", 6},
		{"example01", 8},
		{"example02", 11},
		{"example03", 6},
		{"example04", 6},
		{"example05", 8},
	}

	for _, tc := range examples {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("testdata", "audit", tc.name+".txt")

			input, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			colony, err := parseColony(input)
			if err != nil {
				t.Fatal(err)
			}

			// Run the compiled program.
			output, err := exec.Command(binary, path).CombinedOutput()
			if err != nil {
				t.Fatalf("execution failed: %v\n%s", err, output)
			}

			// Check original input and blank separator.
			prefix := string(input)
			if !strings.HasSuffix(prefix, "\n") {
				prefix += "\n"
			}
			prefix += "\n"

			if !strings.HasPrefix(string(output), prefix) {
				t.Fatal("incorrect input or output formatting")
			}

			// Verify movements and completion time.
			moves := strings.TrimPrefix(string(output), prefix)

			turns := checkSimulation(
				t, colony.Graph, colony.Start,
				colony.End, colony.Ants, moves,
			)

			if turns > tc.maxTurns {
				t.Fatalf("got %d turns, limit %d", turns, tc.maxTurns)
			}

			// Same input must produce the same output.
			again, err := exec.Command(binary, path).CombinedOutput()
			if err != nil || string(again) != string(output) {
				t.Fatal("output changed between runs")
			}

			t.Logf("Passed: %d turns (limit %d)", turns, tc.maxTurns)
		})
	}
}
