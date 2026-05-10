package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"tucil3/game"
	"tucil3/parser"
	"tucil3/searchtrace"
	"tucil3/solver"
)

const maxUploadSize = 2 << 20

type VisualStep struct {
	Label string   `json:"label"`
	Grid  []string `json:"grid"`
}

type SolveResponse struct {
	Found       bool         `json:"found"`
	Moves       string       `json:"moves"`
	TotalCost   int          `json:"totalCost"`
	Iterations  int          `json:"iterations"`
	ExecutionMs int64        `json:"executionMs"`
	Steps       []VisualStep `json:"steps"`
	Trace       GraphTrace   `json:"trace"`
	Error       string       `json:"error,omitempty"`
}

type GraphTrace struct {
	Nodes      []GraphNode `json:"nodes"`
	Edges      []GraphEdge `json:"edges"`
	Expansions []string    `json:"expansions"`
}

type GraphNode struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Row        int    `json:"row"`
	Col        int    `json:"col"`
	NextTarget int    `json:"nextTarget"`
	GCost      int    `json:"gCost"`
	HCost      int    `json:"hCost"`
	Priority   int    `json:"priority"`
	ExpandedAt int    `json:"expandedAt"`
	Solution   bool   `json:"solution"`
}

type GraphEdge struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Direction string `json:"direction"`
	StepCost  int    `json:"stepCost"`
	GCost     int    `json:"gCost"`
	HCost     int    `json:"hCost"`
	Priority  int    `json:"priority"`
	Accepted  bool   `json:"accepted"`
	Solution  bool   `json:"solution"`
	TraceStep int    `json:"traceStep"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/solve", handleSolve)
	if frontendDir := os.Getenv("FRONTEND_DIST"); frontendDir != "" {
		mux.Handle("/", spaFileServer(frontendDir))
	} else {
		mux.HandleFunc("/", handleIndex)
	}

	addr := ":8080"
	if envAddr := os.Getenv("PORT"); envAddr != "" {
		addr = ":" + envAddr
	}

	log.Printf("web backend listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"usage":  "POST multipart/form-data to /solve with file, algorithm, heuristic",
	})
}

func spaFileServer(root string) http.Handler {
	fileServer := http.FileServer(http.Dir(root))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		fullPath := root + string(os.PathSeparator) + path
		if _, err := os.Stat(fullPath); err != nil {
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	})
}

func handleSolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, SolveResponse{Error: "method harus POST"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeJSON(w, http.StatusBadRequest, SolveResponse{Error: "gagal membaca multipart form: " + err.Error()})
		return
	}

	algorithm := normalizeAlgorithm(r.FormValue("algorithm"))
	heuristic := normalizeHeuristic(r.FormValue("heuristic"))

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, SolveResponse{Error: "field file wajib diisi"})
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".txt") {
		writeJSON(w, http.StatusBadRequest, SolveResponse{Error: "file harus berekstensi .txt"})
		return
	}

	tmpFile, err := os.CreateTemp("", "ice-puzzle-*.txt")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, SolveResponse{Error: "gagal membuat file sementara"})
		return
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := io.Copy(tmpFile, file); err != nil {
		tmpFile.Close()
		writeJSON(w, http.StatusInternalServerError, SolveResponse{Error: "gagal menyimpan upload"})
		return
	}
	if err := tmpFile.Close(); err != nil {
		writeJSON(w, http.StatusInternalServerError, SolveResponse{Error: "gagal menutup file sementara"})
		return
	}

	board, err := parser.ParseFile(tmpPath)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, SolveResponse{Error: err.Error()})
		return
	}

	result, trace, err := runSolver(board, algorithm, heuristic)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, SolveResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, SolveResponse{
		Found:       result.Found,
		Moves:       movesToString(result.Moves),
		TotalCost:   result.TotalCost,
		Iterations:  result.Iterations,
		ExecutionMs: result.ExecutionMs,
		Steps:       BuildVisualSteps(board, result),
		Trace:       BuildGraphTrace(result, trace),
	})
}

func normalizeAlgorithm(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func normalizeHeuristic(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case solver.HeuristicManhattan:
		return solver.HeuristicManhattan
	case solver.HeuristicRemaining:
		return solver.HeuristicRemaining
	default:
		return solver.HeuristicManhattanCost
	}
}

func runSolver(board *parser.BoardConfig, algorithm, heuristic string) (solver.Result, []searchtrace.Step, error) {
	switch algorithm {
	case "UCS":
		result, trace := solver.UCSWithTrace(board)
		return result, trace, nil
	case "GBFS":
		result, trace := solver.GBFSWithTrace(board, heuristic)
		return result, trace, nil
	case "ASTAR", "A*":
		result, trace := solver.AStarWithTrace(board, heuristic)
		return result, trace, nil
	default:
		return solver.Result{}, nil, fmt.Errorf("algorithm tidak valid: %q", algorithm)
	}
}

func BuildVisualSteps(board *parser.BoardConfig, result solver.Result) []VisualStep {
	if len(result.PathNodes) == 0 {
		start := graph.Node{Row: board.StartRow, Col: board.StartCol, NextTarget: 0}
		return []VisualStep{{
			Label: "Initial",
			Grid:  gridToStrings(graph.ApplyStateToGrid(board, start)),
		}}
	}

	steps := make([]VisualStep, 0, len(result.PathNodes))
	for i, node := range result.PathNodes {
		label := "Initial"
		if i > 0 {
			label = fmt.Sprintf("Step %d: %s", i, graph.DirName(result.Moves[i-1]))
		}

		steps = append(steps, VisualStep{
			Label: label,
			Grid:  gridToStrings(graph.ApplyStateToGrid(board, node)),
		})
	}

	return steps
}

func BuildGraphTrace(result solver.Result, trace []searchtrace.Step) GraphTrace {
	nodesByID := make(map[string]GraphNode)
	var nodeOrder []string
	var edges []GraphEdge
	var expansions []string

	solutionNodes := make(map[string]bool)
	solutionEdges := make(map[string]bool)
	for i, node := range result.PathNodes {
		id := nodeID(node)
		solutionNodes[id] = true
		if i > 0 {
			prev := result.PathNodes[i-1]
			solutionEdges[edgeKey(prev, node, result.Moves[i-1])] = true
		}
	}

	upsertNode := func(node graph.Node, gCost, hCost, priority, expandedAt int) {
		id := nodeID(node)
		existing, exists := nodesByID[id]
		if !exists {
			nodeOrder = append(nodeOrder, id)
			existing = GraphNode{
				ID:         id,
				Label:      nodeLabel(node),
				Row:        node.Row,
				Col:        node.Col,
				NextTarget: node.NextTarget,
				ExpandedAt: -1,
			}
		}
		if !exists || expandedAt >= 0 || gCost < existing.GCost {
			existing.GCost = gCost
			existing.HCost = hCost
			existing.Priority = priority
		}
		if expandedAt >= 0 && (existing.ExpandedAt == -1 || expandedAt < existing.ExpandedAt) {
			existing.ExpandedAt = expandedAt
		}
		existing.Solution = solutionNodes[id]
		nodesByID[id] = existing
	}

	for stepIndex, step := range trace {
		upsertNode(step.Expanded.State, step.Expanded.GCost, step.Expanded.HCost, step.Expanded.Priority, stepIndex)
		expansions = append(expansions, nodeID(step.Expanded.State))

		for edgeIndex, edge := range step.Edges {
			upsertNode(edge.To, edge.GCost, edge.HCost, edge.Priority, -1)
			fromID := nodeID(edge.From)
			toID := nodeID(edge.To)

			edges = append(edges, GraphEdge{
				ID:        fmt.Sprintf("e-%d-%d-%s-%s", stepIndex, edgeIndex, fromID, toID),
				From:      fromID,
				To:        toID,
				Direction: graph.DirName(edge.Direction),
				StepCost:  edge.StepCost,
				GCost:     edge.GCost,
				HCost:     edge.HCost,
				Priority:  edge.Priority,
				Accepted:  edge.Accepted,
				Solution:  solutionEdges[edgeKey(edge.From, edge.To, edge.Direction)],
				TraceStep: stepIndex,
			})
		}
	}

	nodes := make([]GraphNode, 0, len(nodeOrder))
	for _, id := range nodeOrder {
		nodes = append(nodes, nodesByID[id])
	}

	return GraphTrace{
		Nodes:      nodes,
		Edges:      edges,
		Expansions: expansions,
	}
}

func gridToStrings(grid [][]rune) []string {
	rows := make([]string, 0, len(grid))
	for _, row := range grid {
		rows = append(rows, string(row))
	}
	return rows
}

func movesToString(moves []graph.Dir) string {
	var builder strings.Builder
	for _, move := range moves {
		builder.WriteString(graph.DirName(move))
	}
	return builder.String()
}

func nodeID(node graph.Node) string {
	return fmt.Sprintf("r%d-c%d-t%d", node.Row, node.Col, node.NextTarget)
}

func nodeLabel(node graph.Node) string {
	return fmt.Sprintf("(%d,%d|%d)", node.Row, node.Col, node.NextTarget)
}

func edgeKey(from, to graph.Node, dir graph.Dir) string {
	return fmt.Sprintf("%s>%s:%s", nodeID(from), nodeID(to), graph.DirName(dir))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
