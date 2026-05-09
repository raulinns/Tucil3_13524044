package solver

import "tucil3/parser"

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

// // Algoritma Bonus
// func BFS(board *parser.BoardConfig, heuristicMode string) Result {
// 	return graphSearch(board, AlgorithmBFS, heuristicMode)
// }
