function ResultPanel({ result, currentStep }) {
  return (
    <section className="result-panel">
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
    </section>
  );
}

export default ResultPanel;
