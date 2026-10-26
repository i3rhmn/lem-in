package colony

import "testing"

func TestParserEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"empty file", "", false},
		{"negative ants", "-3\n", false},
		{"nonnumeric ants", "hello\n", false},
		{"overflowing ants", "9999999999999999999999999999\n", false},

		{
			"duplicate end",
			"2\n##start\nS 0 0\n##end\nE 1 0\n" +
				"##end\nX 2 0\nS-E\n",
			false,
		},
		{
			"missing command room",
			"2\n##start\n##end\nE 1 0\nS-E\n",
			false,
		},
		{
			"invalid room name",
			"2\n##start\nLroom 0 0\n##end\nE 1 0\nLroom-E\n",
			false,
		},
		{
			"no tunnels",
			"2\n##start\nS 0 0\n##end\nE 1 0\n",
			false,
		},
		{
			"malformed tunnel",
			"2\n##start\nS 0 0\n##end\nE 1 0\nS - E\n",
			false,
		},
		{
			"ambiguous tunnel",
			"2\n##start\nA 0 0\nA-B 1 0\n" +
				"B-C 2 0\n##end\nC 3 0\nA-B-C\n",
			false,
		},
		{
			"unknown command and comments",
			"2\n#comment\n##unknown\n##start\nS 0 0\n" +
				"##end\nE 1 0\nS-E\n",
			true,
		},
		{
			"hyphenated room name",
			"2\n##start\nS 0 0\nA-B 1 0\n" +
				"##end\nE 2 0\nS-A-B\nA-B-E\n",
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseColony([]byte(tc.input))

			if (err == nil) != tc.valid {
				t.Fatalf(
					"expected valid=%v, got error=%v",
					tc.valid, err,
				)
			}
		})
	}
}
