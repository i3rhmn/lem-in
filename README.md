# Lem-in 🐜

Lem-in is an ant colony simulation written in Go.

The goal is to move all ants from the starting room to the
ending room in as few turns as possible while respecting
room and tunnel restrictions.

## Project Structure

lem-in/
├── main.go
├── go.mod
└── colony/
    ├── run.go
    ├── parser.go
    ├── graph.go
    ├── flow.go
    ├── solver.go
    ├── simulation.go
    ├── *_test.go
    └── testdata/

## How It Works

1. Read the input file.
2. Validate ants, rooms, and tunnels.
3. Construct a graph of the colony.
4. Find suitable paths using minimum-cost flow.
5. Distribute ants across selected paths.
6. Simulate each turn and print movements.

## Algorithms

### Graph Representation

An adjacency list stores the connections between rooms.

Example:

S -> [A, B]
A -> [S, E]
B -> [S, E]
E -> [A, B]

### Minimum-Cost Flow

The solver uses a residual flow network.

Intermediate rooms are split into input and output nodes
with capacity one, preventing paths from sharing them.

Reverse edges allow the solver to reconsider previous
path selections.

The program compares completion times for different
numbers of selected paths.

### Ant Distribution

Ants are assigned based on path length and the number
of ants already assigned to each path.

Shorter and less congested paths are preferred.

### Movement Simulation

Ants move one tunnel per turn.

Each intermediate room can contain only one ant.
Each tunnel can only be used once per turn.

The simulation prints only ants that move.

## Running

Run directly:

    go run . colony/testdata/test0.txt

Or build the executable:

    go build -o lem-in .

    ./lem-in colony/testdata/test0.txt

## Testing

Run all tests:

    go test ./... -count=1

Check code quality:

    go vet ./...

Run race detection:

    CGO_ENABLED=1 go test -race ./...

Run audit examples:

    go test ./colony -run TestAuditExamples -v

Run the performance benchmark:

    go test ./colony -run '^$' -bench BenchmarkOptimizedPaths -benchmem

## Error Handling

Invalid inputs produce:

    ERROR: invalid data format

Examples include missing start/end rooms, invalid ant
counts, duplicate rooms, invalid tunnels, and disconnected
colonies.

## Technologies

- Go
- Standard library only
- Graph theory
- Residual networks
- Minimum-cost flow
- Unit and integration testing
