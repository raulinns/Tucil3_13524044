package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tucil3/graph"
	"tucil3/parser"
	"tucil3/solver"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Ice Sliding Puzzle Solver ===")
	fmt.Print("Masukkan file input: ")

	inputPath, err := readLine(reader)
	if err != nil {
		fmt.Println("Gagal membaca input:", err)
		return
	}

	board, err := parser.ParseFile(inputPath)
	if err != nil {
		fmt.Println("Input tidak valid:", err)
		return
	}

	fmt.Println()
	fmt.Println("Board berhasil dibaca.")
	board.PrintBoard()

	fmt.Println()
	fmt.Println("Pilih algoritma:")
	fmt.Println("1. UCS")
	fmt.Println("2. GBFS")
	fmt.Println("3. A*")
	fmt.Print("Pilihan: ")

	algorithmChoice, err := readLine(reader)
	if err != nil {
		fmt.Println("Gagal membaca pilihan algoritma:", err)
		return
	}

	algorithmChoice = strings.ToUpper(strings.TrimSpace(algorithmChoice))

	heuristicChoice := solver.HeuristicManhattanCost

	if algorithmChoice == "2" || algorithmChoice == "GBFS" ||
		algorithmChoice == "3" || algorithmChoice == "A*" || algorithmChoice == "ASTAR" {

		fmt.Println()
		fmt.Println("Pilih heuristic:")
		fmt.Println("1. H1 - Manhattan Normal")
		fmt.Println("2. H2 - Manhattan x Min Tile Cost")
		fmt.Println("3. H3 - Manhattan Remaining")
		fmt.Print("Pilihan heuristic: ")

		heuristicInput, err := readLine(reader)
		if err != nil {
			fmt.Println("Gagal membaca pilihan heuristic:", err)
			return
		}

		heuristicChoice = normalizeHeuristicChoice(heuristicInput)
	}

	var result solver.Result
	var algorithmName string

	switch algorithmChoice {
	case "1", "UCS":
		algorithmName = "UCS"
		result = solver.UCS(board)

	case "2", "GBFS":
		algorithmName = "GBFS"
		result = solver.GBFS(board, heuristicChoice)

	case "3", "A*", "ASTAR":
		algorithmName = "A*"
		result = solver.AStar(board, heuristicChoice)

	default:
		fmt.Println("Pilihan algoritma tidak valid.")
		return
	}

	fmt.Println()
	fmt.Println("=== Hasil Pencarian ===")
	fmt.Println("Algoritma:", algorithmName)

	if algorithmName == "GBFS" || algorithmName == "A*" {
		fmt.Println("Heuristic:", heuristicChoice)
	}

	if !result.Found {
		fmt.Println("Solusi tidak ditemukan.")
		fmt.Println("Banyak iterasi:", result.Iterations)
		fmt.Printf("Waktu eksekusi: %d ms\n", result.ExecutionMs)
		return
	}

	fmt.Println("Solusi ditemukan.")
	fmt.Println("Solusi gerakan:", movesToString(result.Moves))
	fmt.Println("Cost solusi:", result.TotalCost)
	fmt.Println("Banyak iterasi:", result.Iterations)
	fmt.Printf("Waktu eksekusi: %d ms\n", result.ExecutionMs)

	fmt.Println()
	fmt.Println("=== Visualisasi Solusi ===")
	printSolution(board, result)
}

func readLine(reader *bufio.Reader) (string, error) {
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

func normalizeHeuristicChoice(input string) string {
	choice := strings.ToUpper(strings.TrimSpace(input))

	switch choice {
	case "1", "H1":
		return solver.HeuristicManhattan

	case "2", "H2":
		return solver.HeuristicManhattanCost

	case "3", "H3":
		return solver.HeuristicRemaining

	default:
		fmt.Println("Pilihan heuristic tidak valid. Menggunakan default H2.")
		return solver.HeuristicManhattanCost
	}
}

func movesToString(moves []graph.Dir) string {
	var builder strings.Builder

	for _, move := range moves {
		builder.WriteString(graph.DirName(move))
	}

	return builder.String()
}

func printSolution(board *parser.BoardConfig, result solver.Result) {
	for i, node := range result.PathNodes {
		if i == 0 {
			fmt.Println("Initial")
		} else {
			fmt.Printf("Step %d: %s\n", i, graph.DirName(result.Moves[i-1]))
		}

		grid := graph.ApplyStateToGrid(board, node)
		printGrid(grid)

		fmt.Println()
	}
}

func printGrid(grid [][]rune) {
	for _, row := range grid {
		fmt.Println(string(row))
	}
}
