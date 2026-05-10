package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	graph "tucil3/game"
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
	fmt.Println("4. BFS")
	fmt.Println("5. DFS")
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

	case "4", "BFS":
		algorithmName = "BFS"
		result = solver.BFS(board)

	case "5", "DFS":
		algorithmName = "DFS"
		result = solver.DFS(board)

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
		promptSaveSolution(reader, board, result, algorithmName, heuristicChoice)
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
	promptSaveSolution(reader, board, result, algorithmName, heuristicChoice)
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

func promptSaveSolution(
	reader *bufio.Reader,
	board *parser.BoardConfig,
	result solver.Result,
	algorithmName string,
	heuristicChoice string,
) {
	fmt.Println()
	fmt.Print("Simpan hasil solusi ke file .txt? (y/n): ")
	choice, err := readLine(reader)
	if err != nil {
		fmt.Println("Gagal membaca pilihan simpan:", err)
		return
	}

	choice = strings.ToLower(strings.TrimSpace(choice))
	if choice != "y" && choice != "ya" {
		return
	}

	fmt.Print("Masukkan path output (contoh: test/output.txt): ")
	outputPath, err := readLine(reader)
	if err != nil {
		fmt.Println("Gagal membaca path output:", err)
		return
	}

	outputPath = strings.TrimSpace(outputPath)
	if outputPath == "" {
		fmt.Println("Path output tidak boleh kosong.")
		return
	}

	if !strings.HasSuffix(strings.ToLower(outputPath), ".txt") {
		outputPath += ".txt"
	}

	content := formatSolutionText(board, result, algorithmName, heuristicChoice)
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		fmt.Println("Gagal menyimpan hasil:", err)
		return
	}

	fmt.Println("Hasil solusi disimpan ke:", outputPath)
}

func formatSolutionText(
	board *parser.BoardConfig,
	result solver.Result,
	algorithmName string,
	heuristicChoice string,
) string {
	var builder strings.Builder

	builder.WriteString("Ice Sliding Puzzle Solver\n")
	builder.WriteString("=========================\n\n")
	builder.WriteString("Algoritma: ")
	builder.WriteString(algorithmName)
	builder.WriteString("\n")

	if algorithmName == "GBFS" || algorithmName == "A*" {
		builder.WriteString("Heuristic: ")
		builder.WriteString(heuristicChoice)
		builder.WriteString("\n")
	}

	if result.Found {
		builder.WriteString("Status: Solusi ditemukan\n")
		builder.WriteString("Solusi gerakan: ")
		builder.WriteString(movesToString(result.Moves))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("Cost solusi: %d\n", result.TotalCost))
	} else {
		builder.WriteString("Status: Solusi tidak ditemukan\n")
	}

	builder.WriteString(fmt.Sprintf("Banyak iterasi: %d\n", result.Iterations))
	builder.WriteString(fmt.Sprintf("Waktu eksekusi: %d ms\n\n", result.ExecutionMs))
	builder.WriteString("Visualisasi Solusi\n")
	builder.WriteString("------------------\n")

	if len(result.PathNodes) == 0 {
		start := graph.Node{Row: board.StartRow, Col: board.StartCol, NextTarget: 0}
		builder.WriteString("Initial\n")
		writeGrid(&builder, graph.ApplyStateToGrid(board, start))
		return builder.String()
	}

	for i, node := range result.PathNodes {
		if i == 0 {
			builder.WriteString("Initial\n")
		} else {
			builder.WriteString(fmt.Sprintf("Step %d: %s\n", i, graph.DirName(result.Moves[i-1])))
		}

		writeGrid(&builder, graph.ApplyStateToGrid(board, node))
		builder.WriteString("\n")
	}

	return builder.String()
}

func writeGrid(builder *strings.Builder, grid [][]rune) {
	for _, row := range grid {
		builder.WriteString(string(row))
		builder.WriteString("\n")
	}
}
