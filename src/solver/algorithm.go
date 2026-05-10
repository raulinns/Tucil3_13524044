package solver

import (
	"tucil3/parser"
	"tucil3/searchtrace"
)

// Algoritma Wajib
func UCS(board *parser.BoardConfig) Result {
	return graphSearch(board, AlgorithmUCS, "")
}

func GBFS(board *parser.BoardConfig, heuristicMode string) Result {
	return graphSearch(board, AlgorithmGBFS, heuristicMode)
}

func AStar(board *parser.BoardConfig, heuristicMode string) Result {
	return graphSearch(board, AlgorithmAStar, heuristicMode)
}

func BFS(board *parser.BoardConfig) Result {
	return unweightedGraphSearch(board, AlgorithmBFS, nil)
}

func DFS(board *parser.BoardConfig) Result {
	return unweightedGraphSearch(board, AlgorithmDFS, nil)
}

func UCSWithTrace(board *parser.BoardConfig) (Result, []searchtrace.Step) {
	var trace []searchtrace.Step
	result := graphSearchWithTrace(board, AlgorithmUCS, "", &trace)
	return result, trace
}

func GBFSWithTrace(board *parser.BoardConfig, heuristicMode string) (Result, []searchtrace.Step) {
	var trace []searchtrace.Step
	result := graphSearchWithTrace(board, AlgorithmGBFS, heuristicMode, &trace)
	return result, trace
}

func AStarWithTrace(board *parser.BoardConfig, heuristicMode string) (Result, []searchtrace.Step) {
	var trace []searchtrace.Step
	result := graphSearchWithTrace(board, AlgorithmAStar, heuristicMode, &trace)
	return result, trace
}

func BFSWithTrace(board *parser.BoardConfig) (Result, []searchtrace.Step) {
	var trace []searchtrace.Step
	result := unweightedGraphSearch(board, AlgorithmBFS, &trace)
	return result, trace
}

func DFSWithTrace(board *parser.BoardConfig) (Result, []searchtrace.Step) {
	var trace []searchtrace.Step
	result := unweightedGraphSearch(board, AlgorithmDFS, &trace)
	return result, trace
}
