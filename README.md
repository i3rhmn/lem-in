# lem-in — simple five-file Go version

A Go program that reads an ant colony file and moves the ants from `##start` to `##end` in as few turns as its selected routes allow.

## Structure

```text
lem-in/
├── main.go                 Start the program, read the file, report errors
├── go.mod
└── colony/
    ├── parser.go           Validate and understand ants, rooms and tunnels
    ├── graph.go            IN/OUT capacity network + augmenting path search
    ├── solver.go           Find multiple compatible paths and compare turns
    ├── simulation.go       Assign ants, output their movements and run workflow
    └── testdata/           Sample, error and audit input files
```

Exactly **five Go source files** (all standard library; Go 1.20+).

## Run

From the project's root folder:

```bash
go run . colony/testdata/test0.txt
```

Or build the executable:

```bash
go build -o lem-in .
./lem-in colony/testdata/audit/example00.txt
```

**Input:** first line is a positive ant count, followed by rooms (`name x y`), one `##start` command, one `##end` command, and tunnels (`name-name`). Comments beginning with `#` are ignored for parsing but printed with the original input. Room names may contain hyphens if tunnel names are unambiguous.

**Output:** original file, a blank line, then each turn of ant movements (`L1-room L2-room ...`). If the file is invalid or there is no route, output is `ERROR: invalid data format`.

## How it works (read files in this order)

1. **main.go** — `main()` reads the file and calls `colony.Run()`.
2. **parser.go** — `parseColony()` builds the room inventory and graph (adjacency list).
3. **graph.go** — splits every intermediate room into IN/OUT nodes with capacity 1. `augment()` searches the residual network for another route; reverse edges can undo an earlier choice.
4. **solver.go** — `findOptimizedPaths()` tries different path counts, `paths()` reconstructs selected routes, and `estimatedTurns()` computes how many turns are needed for that set.
5. **simulation.go** — `assignAnts()` distributes ants according to expected arrival time; `simulateAnts()` prints moves turn by turn. `Run()` joins all of the steps.

`##start` and `##end` allow multiple ants, while every intermediate room holds at most one. Each tunnel can be used only once per turn.

## Audit testing

Test the provided input files:

```bash
go run . colony/testdata/test0.txt
go run . colony/testdata/greedy-trap.txt
go run . colony/testdata/audit/example00.txt
go run . colony/testdata/unknown-room.txt
```

All six `colony/testdata/audit/example00.txt` through `example05.txt` were checked before packaging against the expected turn limits (6, 8, 11, 6, 6, 8). The program was also checked with separate automated tests for parser errors, movement legality, direct links, stress cases, rerouting, repeated output, and 4,000 small random graphs compared with exhaustive search. Those temporary test harness files are **not** included in this minimal five-Go-file package.

## Important limitation

The residual minimum-cost flow approach considers shortest-total-length sets for different path counts and then picks the fastest result among those sets. It is stronger than repeatedly running plain BFS, but the tests do **not** constitute a proof that every possible graph gets a globally minimum-turn schedule.

See [AUDIT_GUIDE.md](AUDIT_GUIDE.md) for simple explanations and common questions.
