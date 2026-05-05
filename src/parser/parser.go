package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type BoardConfig struct {
	N, M           int      
	Grid           [][]rune 
	Cost           [][]int  
	StartRow       int      
	StartCol       int
	GoalRow        int 
	GoalCol        int
	NumCheckpoints int            
	CheckpointPos  map[int][2]int 
}

func ParseFile(filename string) (*BoardConfig, error) {
	if !strings.HasSuffix(filename, ".txt") {
		return nil, fmt.Errorf("file harus berekstensi .txt, diberikan: %q", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file: %w", err)
	}

	rawLines := strings.Split(string(data), "\n")
	var lines []string
	for _, l := range rawLines {
		lines = append(lines, strings.TrimRight(l, "\r"))
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("file kosong")
	}

	n, m, err := parseNM(lines[0])
	if err != nil {
		return nil, err
	}

	if len(lines) < n+1 {
		return nil, fmt.Errorf(
			"invalid: jumlah baris tidak cukup; butuh %d baris peta, hanya ada %d baris total",
			n, len(lines)-1,
		)
	}

	cfg := &BoardConfig{
		N:             n,
		M:             m,
		CheckpointPos: make(map[int][2]int),
	}

	grid, err := parseGrid(lines[1:n+1], cfg)
	if err != nil {
		return nil, err
	}
	cfg.Grid = grid

	costTokens := collectTokens(lines[n+1:])
	if len(costTokens) != n*m {
		return nil, fmt.Errorf(
			"invalid: jumlah nilai cost = %d, seharusnya %d (= %d × %d)",
			len(costTokens), n*m, n, m,
		)
	}

	costMatrix, err := parseCostMatrix(costTokens, n, m)
	if err != nil {
		return nil, err
	}
	cfg.Cost = costMatrix

	return cfg, nil
}

func parseNM(line string) (n, m int, err error) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf(
			"invalid: baris pertama harus berisi tepat 2 angka (N M), diberikan: %q", line,
		)
	}
	n, err = strconv.Atoi(parts[0])
	if err != nil || n <= 0 {
		return 0, 0, fmt.Errorf("invalid: N harus bilangan bulat positif, diberikan: %q", parts[0])
	}
	m, err = strconv.Atoi(parts[1])
	if err != nil || m <= 0 {
		return 0, 0, fmt.Errorf("invalid: M harus bilangan bulat positif, diberikan: %q", parts[1])
	}
	return n, m, nil
}

func parseGrid(gridLines []string, cfg *BoardConfig) ([][]rune, error) {
	startCount := 0
	goalCount := 0
	checkpointSet := make(map[int]bool)

	grid := make([][]rune, cfg.N)
	for r, line := range gridLines {
		runes := []rune(line)
		if len(runes) != cfg.M {
			return nil, fmt.Errorf(
				"invalid: baris peta ke-%d memiliki %d karakter, seharusnya %d",
				r+1, len(runes), cfg.M,
			)
		}
		grid[r] = make([]rune, cfg.M)
		for c, ch := range runes {
			switch {
			case ch == 'X' || ch == '*' || ch == 'L':
			case ch == 'Z':
				startCount++
				cfg.StartRow, cfg.StartCol = r, c
			case ch == 'O':
				goalCount++
				cfg.GoalRow, cfg.GoalCol = r, c
			case unicode.IsDigit(ch):
				digit := int(ch - '0')
				if checkpointSet[digit] {
					return nil, fmt.Errorf(
						"invalid: angka checkpoint %d muncul lebih dari satu kali di papan", digit,
					)
				}
				checkpointSet[digit] = true
				cfg.CheckpointPos[digit] = [2]int{r, c}
			default:
				return nil, fmt.Errorf(
					"invalid: karakter tidak dikenal %q pada posisi (%d,%d)", ch, r+1, c+1,
				)
			}
			grid[r][c] = ch
		}
	}

	// Validasi jumlah Z dan O
	if startCount != 1 {
		return nil, fmt.Errorf(
			"invalid: harus ada tepat 1 titik awal (Z), ditemukan %d", startCount,
		)
	}
	if goalCount != 1 {
		return nil, fmt.Errorf(
			"invalid: harus ada tepat 1 titik tujuan (O), ditemukan %d", goalCount,
		)
	}

	// Validasi checkpoint kontinu mulai dari 0
	numCP := len(checkpointSet)
	for i := range numCP {
		if !checkpointSet[i] {
			return nil, fmt.Errorf(
				"invalid: checkpoint harus berurutan mulai dari 0; angka %d tidak ditemukan di papan", i,
			)
		}
	}
	cfg.NumCheckpoints = numCP

	return grid, nil
}

func collectTokens(lines []string) []string {
	var tokens []string
	for _, line := range lines {
		fields := strings.Fields(line)
		tokens = append(tokens, fields...)
	}
	return tokens
}

func parseCostMatrix(tokens []string, n, m int) ([][]int, error) {
	cost := make([][]int, n)
	for r := range n {
		cost[r] = make([]int, m)
		for c := range m {
			idx := r*m + c
			val, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf(
					"invalid: nilai cost token ke-%d (%q) bukan bilangan bulat", idx+1, tokens[idx],
				)
			}
			if val < 0 {
				return nil, fmt.Errorf(
					"invalid: nilai cost tidak boleh negatif, ditemukan %d pada posisi (%d,%d)",
					val, r+1, c+1,
				)
			}
			cost[r][c] = val
		}
	}
	return cost, nil
}

func (b *BoardConfig) PrintBoard() {
	fmt.Printf("Papan %d×%d | Start: (%d,%d) | Goal: (%d,%d) | Checkpoints: %d\n",
		b.N, b.M, b.StartRow, b.StartCol, b.GoalRow, b.GoalCol, b.NumCheckpoints)
	fmt.Println("Grid:")
	for _, row := range b.Grid {
		fmt.Printf("  %s\n", string(row))
	}
	fmt.Println("Cost:")
	for _, row := range b.Cost {
		fmt.Print("  ")
		for j, v := range row {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("%4d", v)
		}
		fmt.Println()
	}
}
