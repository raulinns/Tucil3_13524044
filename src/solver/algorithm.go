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

// // Algoritma Bonus
// func BFS(board *parser.BoardConfig, heuristicMode string) Result {
// 	return graphSearch(board, AlgorithmBFS, heuristicMode)
// }
