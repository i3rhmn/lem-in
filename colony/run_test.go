package colony

import (
	"io"
	"os"
	"testing"
)

// Run the colony and capture its output.
func captureRun(t *testing.T, input string) (string, error) {
	t.Helper()

	old := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = writer

	runErr := Run([]byte(input))

	writer.Close()
	os.Stdout = old

	output, err := io.ReadAll(reader)
	reader.Close()

	if err != nil {
		t.Fatal(err)
	}

	return string(output), runErr
}

// Verify the exact output format.
func TestRunOutput(t *testing.T) {
	input := "3\n" +
		"##start\nS 0 0\n" +
		"A 1 0\n" +
		"##end\nE 2 0\n" +
		"S-A\nA-E\n"

	expected := input + "\n" +
		"L1-A\n" +
		"L1-E L2-A\n" +
		"L2-E L3-A\n" +
		"L3-E\n"

	output, err := captureRun(t, input)

	if err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Fatalf(
			"incorrect output\nExpected:\n%q\nGot:\n%q",
			expected, output,
		)
	}
}

// Verify input without a trailing newline.
func TestRunWithoutFinalNewline(t *testing.T) {
	input := "2\n" +
		"##start\nS 0 0\n" +
		"##end\nE 1 0\n" +
		"S-E"

	expected := input + "\n\n" +
		"L1-E\n" +
		"L2-E\n"

	output, err := captureRun(t, input)

	if err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Fatalf("unexpected output:\n%q", output)
	}
}

// Unknown commands and comments should be preserved.
func TestRunWithComments(t *testing.T) {
	input := "1\n" +
		"##unknown\n" +
		"##start\nS 0 0\n" +
		"# This is a comment\n" +
		"##end\nE 1 0\n" +
		"S-E\n"

	expected := input + "\nL1-E\n"

	output, err := captureRun(t, input)

	if err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Fatalf("unexpected output:\n%q", output)
	}
}

// Invalid colonies must return an error.
func TestRunInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"zero ants",
			"0\n##start\nS 0 0\n##end\nE 1 0\nS-E\n",
		},
		{
			"missing end",
			"2\n##start\nS 0 0\nA 1 0\nS-A\n",
		},
		{
			"duplicate link",
			"2\n##start\nS 0 0\n##end\nE 1 0\nS-E\nE-S\n",
		},
		{
			"invalid coordinates",
			"2\n##start\nS abc 0\n##end\nE 1 0\nS-E\n",
		},
		{
			"disconnected colony",
			"2\n##start\nS 0 0\nA 1 0\n" +
				"##end\nE 2 0\nB 3 0\nS-A\nB-E\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := captureRun(t, tt.input)

			if err == nil {
				t.Fatal("expected an error")
			}

			if output != "" {
				t.Fatalf(
					"invalid colony produced output: %q",
					output,
				)
			}
		})
	}
}
