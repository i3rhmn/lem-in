package colony

// Every room (except start/end) will be split into IN and OUT nodes.
// A capacity-one IN -> OUT edge lets a path reserve that room once.
type edge struct {
	to, reverse, capacity, cost int
	fromRoom, toRoom            string
}

type network struct{ edges [][]*edge }

// Adds forward and reverse connections
func (f *network) add(from, to, cost int, a, b string) {
	forward := &edge{to: to, reverse: len(f.edges[to]), capacity: 1, cost: cost, fromRoom: a, toRoom: b}
	backward := &edge{to: from, reverse: len(f.edges[from]), cost: -cost}
	f.edges[from] = append(f.edges[from], forward)
	f.edges[to] = append(f.edges[to], backward)
}

// A BFS-like shortest-path search finds an augmenting route.
// Reverse edges allow an earlier route to be rearranged.

//finds route and reserves it.
func (f *network) augment(start, end int) bool {
	const inf = int(^uint(0)>>1) / 4
	dist := make([]int, len(f.edges))
	parent, parentEdge := make([]int, len(f.edges)), make([]int, len(f.edges))
	queued := make([]bool, len(f.edges))
	for i := range dist {
		dist[i], parent[i] = inf, -1
	}
	dist[start] = 0
	queue := []int{start}
	queued[start] = true
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		queued[u] = false
		for i, e := range f.edges[u] {
			if e.capacity == 0 || dist[u]+e.cost >= dist[e.to] {
				continue
			}
			dist[e.to], parent[e.to], parentEdge[e.to] = dist[u]+e.cost, u, i
			if !queued[e.to] {
				queue = append(queue, e.to)
				queued[e.to] = true
			}
		}
	}
	if parent[end] == -1 {
		return false
	}
	for v := end; v != start; v = parent[v] {
		u := parent[v]
		e := f.edges[u][parentEdge[v]]
		e.capacity--
		f.edges[v][e.reverse].capacity++
	}
	return true
}
