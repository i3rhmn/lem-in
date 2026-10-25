package lemin

import "testing"

func TestParseColony(t *testing.T) {

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name: "valid colony",
			input: "3\n##start\nA 0 0\n" +
				"##end\nB 1 1\nA-B\n",
			valid: true,
		},
		{
			name:  "invalid ants",
			input: "0\n",
			valid: false,
		},
		{
			name: "missing start",
			input: "3\nA 0 0\n" +
				"##end\nB 1 1\nA-B\n",
			valid: false,
		},
		{
			name: "duplicate room",
			input: "3\n##start\nA 0 0\n" +
				"A 1 1\n##end\nB 2 2\nA-B\n",
			valid: false,
		},
		{
			name: "invalid coordinates",
			input: "3\n##start\nA hello 0\n" +
				"##end\nB 1 1\nA-B\n",
			valid: false,
		},
		{
			name: "unknown room",
			input: "3\n##start\nA 0 0\n" +
				"##end\nB 1 1\nA-C\n",
			valid: false,
		},
		{
			name: "duplicate tunnel",
			input: "3\n##start\nA 0 0\n" +
				"##end\nB 1 1\nA-B\nB-A\n",
			valid: false,
		},
		{
			name: "self tunnel",
			input: "3\n##start\nA 0 0\n" +
				"##end\nB 1 1\nA-A\n",
			valid: false,
		},
		{
			name: "room after tunnel",
			input: "3\n##start\nA 0 0\n" +
				"##end\nB 1 1\nA-B\nC 2 2\n",
			valid: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			_, err := parseColony([]byte(tt.input))

			if (err == nil) != tt.valid {
				t.Fatalf(
					"expected valid=%v, got error=%v",
					tt.valid,
					err,
				)
			}
		})
	}
}
