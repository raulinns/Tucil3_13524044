package solver

import (
	"math"

	"tucil3/graph"
	"tucil3/parser"
)

const (
	HeuristicManhattan     = "H1"
	HeuristicManhattanCost = "H2"
	HeuristicRemaining     = "H3"
)

// Heuristic menghitung nilai h(n) berdasarkan mode heuristic yang dipilih.
//
// H1 = Manhattan normal ke target berikutnya.
// H2 = Manhattan ke target berikutnya dikali minTileCost.
// H3 = Total Manhattan menuju seluruh checkpoint tersisa lalu goal.
func Heuristic(board *parser.BoardConfig, node graph.Node, mode string) int {
	switch mode {
	case HeuristicManhattan:
		return manhattanToNextTarget(board, node)

	case HeuristicManhattanCost:
		return manhattanToNextTarget(board, node) * minTraversableCost(board)

	case HeuristicRemaining:
		return remainingManhattan(board, node)

	default:
		// Default heuristic utama
		return manhattanToNextTarget(board, node) * minTraversableCost(board)
	}
}

// manhattanToNextTarget menghitung jarak Manhattan dari node saat ini
// ke target berikutnya, yaitu checkpoint berikutnya jika masih ada,
// atau goal jika semua checkpoint sudah dilewati.
func manhattanToNextTarget(board *parser.BoardConfig, node graph.Node) int {
	targetRow, targetCol := nextTargetPosition(board, node)

	return abs(node.Row-targetRow) + abs(node.Col-targetCol)
}

// remainingManhattan menghitung total jarak Manhattan dari posisi saat ini
// ke semua checkpoint yang tersisa secara berurutan, lalu ke goal.
func remainingManhattan(board *parser.BoardConfig, node graph.Node) int {
	total := 0

	currentRow := node.Row
	currentCol := node.Col

	for target := node.NextTarget; target < board.NumCheckpoints; target++ {
		pos := board.CheckpointPos[target]

		targetRow := pos[0]
		targetCol := pos[1]

		total += abs(currentRow-targetRow) + abs(currentCol-targetCol)

		currentRow = targetRow
		currentCol = targetCol
	}

	total += abs(currentRow-board.GoalRow) + abs(currentCol-board.GoalCol)

	return total
}

// nextTargetPosition mengembalikan posisi target berikutnya.
// Jika masih ada checkpoint yang belum dilewati, targetnya adalah checkpoint NextTarget.
// Jika semua checkpoint sudah dilewati, targetnya adalah goal.
func nextTargetPosition(board *parser.BoardConfig, node graph.Node) (int, int) {
	if node.NextTarget < board.NumCheckpoints {
		pos := board.CheckpointPos[node.NextTarget]
		return pos[0], pos[1]
	}

	return board.GoalRow, board.GoalCol
}

// minTraversableCost mencari cost terkecil dari seluruh tile yang secara tipe bisa dilewati.
func minTraversableCost(board *parser.BoardConfig) int {
	minCost := math.MaxInt

	for r := 0; r < board.N; r++ {
		for c := 0; c < board.M; c++ {
			tile := board.Grid[r][c]

			if tile == 'X' || tile == 'L' {
				continue
			}

			if board.Cost[r][c] < minCost {
				minCost = board.Cost[r][c]
			}
		}
	}

	if minCost == math.MaxInt {
		return 0
	}

	return minCost
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
