package searchtrace

import "tucil3/game"

// Step mencatat satu ekspansi node untuk kebutuhan visualisasi pencarian.
type Step struct {
	Expanded Node
	Edges    []Edge
}

type Node struct {
	State    graph.Node
	GCost    int
	HCost    int
	Priority int
}

type Edge struct {
	From      graph.Node
	To        graph.Node
	Direction graph.Dir
	StepCost  int
	GCost     int
	HCost     int
	Priority  int
	Accepted  bool
}
