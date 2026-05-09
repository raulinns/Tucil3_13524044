package game

import "tucil3/parser"

type Dir int

const (
	Up    Dir = iota // atas
	Down             // bawah
	Left             // kiri
	Right            // kanan
)

func DirName(d Dir) string {
	return [...]string{"U", "D", "L", "R"}[d]
}

var AllDirs = []Dir{Up, Down, Left, Right}

// State setiap kemungkinan node, dalam implementasinya node yang tidak valid (menyebabkan game over) tidak dicatat
type Node struct {
	Row, Col   int
	NextTarget int // checkpoint berikutnya; jika == board.NumCheckpoints, semua checkpoint sudah dilewati
}

// Gerakan slide yang valid dari current Node ke Node lainnya
type Edge struct {
	State     Node
	Cost      int // total cost tile yang dilalui dalam gerakan ini
	Direction Dir
	Path      []Pos // tile-tile yang dilalui
}

type Pos struct {
	Row, Col int
}

func IsGoal(board *parser.BoardConfig, s Node) bool {
	return s.Row == board.GoalRow &&
		s.Col == board.GoalCol &&
		s.NextTarget == board.NumCheckpoints
}

func Slide(board *parser.BoardConfig, s Node, dir Dir) (newState Node, cost int, path []Pos, valid bool) {
	// Arah untuk perubahan posisi di grid (array dua dimensi)
	dr, dc := slideDir(dir)

	curRow, curCol := s.Row, s.Col
	nextTarget := s.NextTarget
	totalCost := 0
	var tilePath []Pos

	for {
		nextRow := curRow + dr
		nextCol := curCol + dc

		// Validasi keluar papan, return false (game over)
		if nextRow < 0 || nextRow >= board.N || nextCol < 0 || nextCol >= board.M {
			return Node{}, 0, nil, false
		}

		// Pengecekan jenis tileTarget setelah bergerak
		tileTarget := board.Grid[nextRow][nextCol]

		// Berhenti ketika tile di depan merupakan tembok (X)
		if tileTarget == 'X' {
			break
		}

		// Melewati L, gameover
		if tileTarget == 'L' {
			return Node{}, 0, nil, false
		}

		// Melewati checkpoint
		if tileTarget >= '0' && tileTarget <= '9' {
			digit := int(tileTarget - '0')
			// Melewati checkpoint yang belum dapat dilewati, game over
			if digit > nextTarget {
				return Node{}, 0, nil, false
			}
			// Melewati checkpoint yang benar, target diinkremen
			if digit == nextTarget {
				nextTarget++
			}
			// Melewati checkpoint yang sudah dilewati, tidak perlu melakukan apa-apa, dianggap tile normal (*)
		}

		// Gerak ke tile berikutnya, tambahkan cost dan catat path
		curRow, curCol = nextRow, nextCol
		totalCost += board.Cost[curRow][curCol]
		tilePath = append(tilePath, Pos{curRow, curCol})
	}

	// Jika karakter tidak bergerak, maka state tidak valid (Misal, karena langsung bertemu dengan tembok)
	if curRow == s.Row && curCol == s.Col {
		return Node{}, 0, nil, false
	}

	newState = Node{Row: curRow, Col: curCol, NextTarget: nextTarget}
	return newState, totalCost, tilePath, true
}

func slideDir(dir Dir) (dr, dc int) {
	switch dir {
	case Up:
		return -1, 0
	case Down:
		return 1, 0
	case Left:
		return 0, -1
	case Right:
		return 0, 1
	}
	return 0, 0
}

// Membangun graf pencarian secara implisit, setiap node yang sedang dikunjungi akan dicek node selanjutnya yang valid
func GetSuccessors(board *parser.BoardConfig, s Node) []Edge {
	var successors []Edge
	// Mengecek ke semua arah untuk melihat gerakan yang valid
	// Kemungkinan next node yang dicatat hanya yang dapat dilewati dan tidak menyebabkan game over
	for _, dir := range AllDirs {
		newState, cost, path, valid := Slide(board, s, dir)
		if valid {
			successors = append(successors, Edge{
				State:     newState,
				Cost:      cost,
				Direction: dir,
				Path:      path,
			})
		}
	}
	return successors
}

// Deep copy current Grid
func CopyGrid(grid [][]rune) [][]rune {
	cp := make([][]rune, len(grid))
	for i, row := range grid {
		cp[i] = make([]rune, len(row))
		copy(cp[i], row)
	}
	return cp
}

// ApplyStateToGrid mengembalikan salinan grid yang sudah disesuaikan
// dengan posisi aktor dan checkpoint yang sudah dilewati pada node s.
func ApplyStateToGrid(board *parser.BoardConfig, s Node) [][]rune {
	g := CopyGrid(board.Grid)

	// Tandai checkpoint yang sudah dilewati sebagai tile biasa
	for i := range s.NextTarget {
		pos, valid := board.CheckpointPos[i]
		if valid {
			g[pos[0]][pos[1]] = '*'
		}
	}

	// Tandai petak yang sedang diinjak sebagai Z
	g[s.Row][s.Col] = 'Z'

	return g
}
