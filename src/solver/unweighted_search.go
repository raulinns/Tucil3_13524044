package solver

import (
	"time"

	graph "tucil3/game"
	"tucil3/parser"
	"tucil3/searchtrace"
)

func unweightedGraphSearch(
	board *parser.BoardConfig,
	algorithm Algorithm,
	trace *[]searchtrace.Step,
) Result {
	startTime := time.Now()
	start := graph.Node{
		Row:        board.StartRow,
		Col:        board.StartCol,
		NextTarget: 0,
	}

	startNode := &SearchNode{
		Node:     start,
		GCost:    0,
		Priority: 0,
		Depth:    0,
	}

	frontier := []*SearchNode{startNode}
	visited := map[graph.Node]bool{start: true}
	iterations := 0

	for len(frontier) > 0 {
		var current *SearchNode
		if algorithm == AlgorithmDFS {
			last := len(frontier) - 1
			current = frontier[last]
			frontier = frontier[:last]
		} else {
			current = frontier[0]
			frontier = frontier[1:]
		}

		iterations++

		if graph.IsGoal(board, current.Node) {
			result := reconstructResult(current)
			result.Found = true
			result.TotalCost = current.GCost
			result.Iterations = iterations
			result.ExecutionMs = time.Since(startTime).Milliseconds()
			return result
		}

		edges := graph.GetSuccessors(board, current.Node)
		traceStep := buildUnweightedTraceStep(trace, current)
		pending := make([]*SearchNode, 0, len(edges))

		for _, edge := range edges {
			nextNode := edge.State
			nextDepth := current.Depth + 1
			nextGCost := current.GCost + edge.Cost
			accepted := !visited[nextNode]

			appendTraceEdge(traceStep, searchtrace.Edge{
				From:      current.Node,
				To:        nextNode,
				Direction: edge.Direction,
				StepCost:  edge.Cost,
				GCost:     nextGCost,
				HCost:     0,
				Priority:  nextDepth,
				Accepted:  accepted,
			})

			if !accepted {
				continue
			}

			visited[nextNode] = true
			pending = append(pending, &SearchNode{
				Node:     nextNode,
				Parent:   current,
				Move:     edge.Direction,
				GCost:    nextGCost,
				HCost:    0,
				Priority: nextDepth,
				Depth:    nextDepth,
				StepPath: edge.Path,
			})
		}

		if algorithm == AlgorithmDFS {
			for i := len(pending) - 1; i >= 0; i-- {
				frontier = append(frontier, pending[i])
			}
		} else {
			frontier = append(frontier, pending...)
		}

		if traceStep != nil {
			*trace = append(*trace, *traceStep)
		}
	}

	return Result{
		Found:       false,
		Iterations:  iterations,
		ExecutionMs: time.Since(startTime).Milliseconds(),
	}
}

func buildUnweightedTraceStep(trace *[]searchtrace.Step, current *SearchNode) *searchtrace.Step {
	if trace == nil {
		return nil
	}

	return &searchtrace.Step{
		Expanded: searchtrace.Node{
			State:    current.Node,
			GCost:    current.GCost,
			HCost:    0,
			Priority: current.Depth,
		},
	}
}
