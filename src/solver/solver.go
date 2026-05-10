package solver

import (
	"container/heap"
	"time"

	"tucil3/game"
	"tucil3/parser"
	"tucil3/searchtrace"
)

type Algorithm int

const (
	AlgorithmUCS Algorithm = iota
	AlgorithmGBFS
	AlgorithmAStar
)

// Node pada search tree
// graph.Node menyimpan state permainan, searchNode menyimpan node untuk pathfinding
type SearchNode struct {
	Node     graph.Node
	Parent   *SearchNode // SearchNode sebelumnya
	Move     graph.Dir
	GCost    int
	HCost    int
	Priority int
	StepPath []graph.Pos

	Index int
}

// Menyimpan hasil pencarian
type Result struct {
	Found bool

	Moves     []graph.Dir
	PathNodes []graph.Node
	StepPaths [][]graph.Pos

	TotalCost   int
	Iterations  int
	ExecutionMs int64
}

func graphSearch(board *parser.BoardConfig, algorithm Algorithm, heuristicMode string) Result {
	return graphSearchWithTrace(board, algorithm, heuristicMode, nil)
}

func graphSearchWithTrace(
	board *parser.BoardConfig,
	algorithm Algorithm,
	heuristicMode string,
	trace *[]searchtrace.Step,
) Result {
	startTime := time.Now()

	start := graph.Node{
		Row:        board.StartRow,
		Col:        board.StartCol,
		NextTarget: 0,
	}

	startH := 0
	if algorithm == AlgorithmGBFS || algorithm == AlgorithmAStar {
		startH = Heuristic(board, start, heuristicMode)
	}

	startNode := &SearchNode{
		Node:     start,
		Parent:   nil,
		GCost:    0,
		HCost:    startH,
		Priority: calculatePriority(algorithm, 0, startH),
		StepPath: nil,
	}

	agenda := &PriorityQueue{}
	heap.Init(agenda)
	heap.Push(agenda, startNode)

	bestCost := make(map[graph.Node]int)
	bestCost[start] = 0

	iterations := 0

	for agenda.Len() > 0 {
		current := heap.Pop(agenda).(*SearchNode)

		// Jika node ini bukan versi termurah yang pernah ditemukan,
		// maka node ini diabaikan
		if best, valid := bestCost[current.Node]; valid && current.GCost > best {
			continue
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
		var traceStep *searchtrace.Step
		if trace != nil {
			traceStep = &searchtrace.Step{
				Expanded: searchtrace.Node{
					State:    current.Node,
					GCost:    current.GCost,
					HCost:    current.HCost,
					Priority: current.Priority,
				},
			}
		}

		for _, edge := range edges {
			nextNode := edge.State
			newGCost := current.GCost + edge.Cost
			hCost := 0
			if algorithm == AlgorithmGBFS || algorithm == AlgorithmAStar {
				hCost = Heuristic(board, nextNode, heuristicMode)
			}
			priority := calculatePriority(algorithm, newGCost, hCost)

			oldGCost, visited := bestCost[nextNode]
			if visited && newGCost >= oldGCost {
				appendTraceEdge(traceStep, searchtrace.Edge{
					From:      current.Node,
					To:        nextNode,
					Direction: edge.Direction,
					StepCost:  edge.Cost,
					GCost:     newGCost,
					HCost:     hCost,
					Priority:  priority,
					Accepted:  false,
				})
				continue
			}

			bestCost[nextNode] = newGCost

			searchNode := &SearchNode{
				Node:     nextNode,
				Parent:   current,
				Move:     edge.Direction,
				GCost:    newGCost,
				HCost:    hCost,
				Priority: priority,
				StepPath: edge.Path,
			}

			appendTraceEdge(traceStep, searchtrace.Edge{
				From:      current.Node,
				To:        nextNode,
				Direction: edge.Direction,
				StepCost:  edge.Cost,
				GCost:     newGCost,
				HCost:     hCost,
				Priority:  priority,
				Accepted:  true,
			})
			heap.Push(agenda, searchNode)
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

func appendTraceEdge(step *searchtrace.Step, edge searchtrace.Edge) {
	if step != nil {
		step.Edges = append(step.Edges, edge)
	}
}

func calculatePriority(algorithm Algorithm, gCost, hCost int) int {
	switch algorithm {
	case AlgorithmUCS:
		return gCost
	case AlgorithmGBFS:
		return hCost
	case AlgorithmAStar:
		return gCost + hCost
	default:
		return gCost
	}
}

// Helper untuk rekonstruksi solusi dari Goal untuk mencetak solusi lengkapnya
func reconstructResult(goal *SearchNode) Result {
	var moves []graph.Dir
	var pathNodes []graph.Node
	var stepPaths [][]graph.Pos

	current := goal

	for current != nil {
		pathNodes = append(pathNodes, current.Node)

		if current.Parent != nil {
			moves = append(moves, current.Move)
			stepPaths = append(stepPaths, current.StepPath)
		}

		current = current.Parent
	}

	reverseNodes(pathNodes)
	reverseDirs(moves)
	reverseStepPaths(stepPaths)

	return Result{
		Found:     true,
		Moves:     moves,
		PathNodes: pathNodes,
		StepPaths: stepPaths,
	}
}

func reverseNodes(nodes []graph.Node) {
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
}

func reverseDirs(dirs []graph.Dir) {
	for i, j := 0, len(dirs)-1; i < j; i, j = i+1, j-1 {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	}
}

func reverseStepPaths(paths [][]graph.Pos) {
	for i, j := 0, len(paths)-1; i < j; i, j = i+1, j-1 {
		paths[i], paths[j] = paths[j], paths[i]
	}
}
