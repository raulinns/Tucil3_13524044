function ResultPanel({ result, currentStep }) {
  function handleSave() {
    if (!result.steps?.length) {
      return;
    }

    const blob = new Blob([formatSolutionText(result)], { type: "text/plain;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    const algorithm = (result.algorithm || "solver").replace(/[^a-z0-9]+/gi, "-").toLowerCase();

    link.href = url;
    link.download = `ice-sliding-${algorithm}-solution.txt`;
    document.body.appendChild(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
  }

  return (
    <section className="result-panel">
      <div className="metric-row">
        <span>Algorithm</span>
        <strong>{result.algorithm || "-"}</strong>
      </div>
      {result.heuristic && (
        <div className="metric-row">
          <span>Heuristic</span>
          <strong>{result.heuristic}</strong>
        </div>
      )}
      <div className="metric-row">
        <span>Status</span>
        <strong>{result.steps?.length ? (result.found ? "Found" : "Not found") : "-"}</strong>
      </div>
      <div className="metric-row">
        <span>Moves</span>
        <strong>{result.moves || "-"}</strong>
      </div>
      <div className="metric-row">
        <span>Total cost</span>
        <strong>{result.totalCost ?? 0}</strong>
      </div>
      <div className="metric-row">
        <span>Iterations</span>
        <strong>{result.iterations ?? 0}</strong>
      </div>
      <div className="metric-row">
        <span>Execution</span>
        <strong>{result.executionMs ?? 0} ms</strong>
      </div>
      <div className="step-label">{currentStep?.label || "No step selected"}</div>
      <button className="secondary-button" type="button" disabled={!result.steps?.length} onClick={handleSave}>
        Save .txt
      </button>
    </section>
  );
}

function formatSolutionText(result) {
  const lines = [
    "Ice Sliding Puzzle Solver",
    "=========================",
    "",
    `Algoritma: ${result.algorithm || "-"}`,
  ];

  if (result.heuristic) {
    lines.push(`Heuristic: ${result.heuristic}`);
  }

  if (result.found) {
    lines.push("Status: Solusi ditemukan");
    lines.push(`Solusi gerakan: ${result.moves || "-"}`);
    lines.push(`Cost solusi: ${result.totalCost ?? 0}`);
  } else {
    lines.push("Status: Solusi tidak ditemukan");
  }

  lines.push(`Banyak iterasi: ${result.iterations ?? 0}`);
  lines.push(`Waktu eksekusi: ${result.executionMs ?? 0} ms`);
  lines.push("");
  lines.push("Visualisasi Solusi");
  lines.push("------------------");

  for (const step of result.steps || []) {
    lines.push(step.label);
    lines.push(...step.grid);
    lines.push("");
  }

  return `${lines.join("\n")}\n`;
}

export default ResultPanel;
