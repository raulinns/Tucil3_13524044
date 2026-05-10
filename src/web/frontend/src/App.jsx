import { useEffect, useMemo, useState } from "react";
import Board from "./components/Board.jsx";
import ControlPanel from "./components/ControlPanel.jsx";
import GraphView from "./components/GraphView.jsx";
import PlaybackControls from "./components/PlaybackControls.jsx";
import ResultPanel from "./components/ResultPanel.jsx";

const initialResult = {
  algorithm: "",
  heuristic: "",
  found: false,
  moves: "",
  totalCost: 0,
  iterations: 0,
  executionMs: 0,
  steps: [],
  trace: { nodes: [], edges: [], expansions: [] },
};

function App() {
  const [result, setResult] = useState(initialResult);
  const [currentStep, setCurrentStep] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [speed, setSpeed] = useState(500);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const steps = result.steps || [];
  const lastStepIndex = Math.max(steps.length - 1, 0);

  useEffect(() => {
    if (!isPlaying || steps.length <= 1) {
      return undefined;
    }

    const timer = window.setInterval(() => {
      setCurrentStep((step) => {
        if (step >= lastStepIndex) {
          setIsPlaying(false);
          return step;
        }
        return step + 1;
      });
    }, speed);

    return () => window.clearInterval(timer);
  }, [isPlaying, lastStepIndex, speed, steps.length]);

  const currentVisualStep = steps[currentStep] || null;

  const activeTraceIndex = useMemo(() => {
    if (!result.trace?.expansions?.length) {
      return -1;
    }
    return Math.min(currentStep, result.trace.expansions.length - 1);
  }, [currentStep, result.trace]);

  async function handleSolve({ file, algorithm, heuristic }) {
    setLoading(true);
    setError("");
    setIsPlaying(false);
    setCurrentStep(0);

    const body = new FormData();
    body.append("file", file);
    body.append("algorithm", algorithm);
    body.append("heuristic", heuristic);

    try {
      const response = await fetch("/solve", {
        method: "POST",
        body,
      });
      const payload = await response.json();

      if (!response.ok || payload.error) {
        throw new Error(payload.error || "Solve request failed");
      }

      setResult(payload);
    } catch (err) {
      setResult(initialResult);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="app-shell">
      <section className="topbar">
        <div>
          <p className="eyebrow">IF2211 Strategi Algoritma</p>
          <h1>Ice Sliding Puzzle Solver</h1>
        </div>
        <div className="status-pill">
          {loading ? "Solving" : result.steps?.length ? "Ready" : "Idle"}
        </div>
      </section>

      <section className="workspace">
        <aside className="side-panel">
          <ControlPanel onSolve={handleSolve} loading={loading} />
          {error && <div className="error-box">{error}</div>}
          <ResultPanel result={result} currentStep={currentVisualStep} />
          <PlaybackControls
            currentStep={currentStep}
            isPlaying={isPlaying}
            maxStep={lastStepIndex}
            speed={speed}
            disabled={!steps.length}
            onStepChange={setCurrentStep}
            onPlayingChange={setIsPlaying}
            onSpeedChange={setSpeed}
          />
        </aside>

        <section className="main-panel">
          <div className="board-region">
            <Board grid={currentVisualStep?.grid || []} />
          </div>
          <GraphView trace={result.trace} activeTraceIndex={activeTraceIndex} />
        </section>
      </section>
    </main>
  );
}

export default App;
