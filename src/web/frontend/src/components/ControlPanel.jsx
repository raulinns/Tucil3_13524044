import { useState } from "react";

const algorithms = [
  { value: "UCS", label: "UCS" },
  { value: "GBFS", label: "GBFS" },
  { value: "ASTAR", label: "A*" },
  { value: "BFS", label: "BFS", disabled: true },
  { value: "DFS", label: "DFS", disabled: true },
];

const heuristics = [
  { value: "H1", label: "H1 - Manhattan" },
  { value: "H2", label: "H2 - Manhattan x min cost" },
  { value: "H3", label: "H3 - Remaining checkpoints" },
];

function ControlPanel({ onSolve, loading }) {
  const [file, setFile] = useState(null);
  const [algorithm, setAlgorithm] = useState("UCS");
  const [heuristic, setHeuristic] = useState("H2");
  const needsHeuristic = algorithm === "GBFS" || algorithm === "ASTAR";

  function submit(event) {
    event.preventDefault();
    if (!file || loading) {
      return;
    }
    onSolve({ file, algorithm, heuristic });
  }

  return (
    <form className="tool-panel" onSubmit={submit}>
      <div className="field">
        <label htmlFor="file">Testcase</label>
        <input
          id="file"
          type="file"
          accept=".txt"
          onChange={(event) => setFile(event.target.files?.[0] || null)}
        />
      </div>

      <div className="field">
        <label htmlFor="algorithm">Algorithm</label>
        <select
          id="algorithm"
          value={algorithm}
          onChange={(event) => setAlgorithm(event.target.value)}
        >
          {algorithms.map((option) => (
            <option key={option.value} value={option.value} disabled={option.disabled}>
              {option.label}
            </option>
          ))}
        </select>
      </div>

      <div className="field">
        <label htmlFor="heuristic">Heuristic</label>
        <select
          id="heuristic"
          value={heuristic}
          disabled={!needsHeuristic}
          onChange={(event) => setHeuristic(event.target.value)}
        >
          {heuristics.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </div>

      <button className="primary-button" type="submit" disabled={!file || loading}>
        {loading ? "Solving..." : "Solve"}
      </button>
    </form>
  );
}

export default ControlPanel;
