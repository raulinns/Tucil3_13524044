function PlaybackControls({
  currentStep,
  disabled,
  isPlaying,
  maxStep,
  speed,
  onPlayingChange,
  onSpeedChange,
  onStepChange,
}) {
  function clampStep(value) {
    return Math.min(Math.max(value, 0), maxStep);
  }

  return (
    <section className="playback-panel">
      <div className="playback-buttons">
        <button
          type="button"
          disabled={disabled || currentStep <= 0}
          onClick={() => onStepChange(clampStep(currentStep - 1))}
        >
          Prev
        </button>
        <button
          type="button"
          disabled={disabled}
          onClick={() => onPlayingChange(!isPlaying)}
        >
          {isPlaying ? "Pause" : "Play"}
        </button>
        <button
          type="button"
          disabled={disabled || currentStep >= maxStep}
          onClick={() => onStepChange(clampStep(currentStep + 1))}
        >
          Next
        </button>
      </div>

      <label className="range-field">
        Step
        <input
          type="range"
          min="0"
          max={maxStep}
          value={currentStep}
          disabled={disabled}
          onChange={(event) => onStepChange(Number(event.target.value))}
        />
        <span>
          {currentStep} / {maxStep}
        </span>
      </label>

      <label className="range-field">
        Speed
        <input
          type="range"
          min="150"
          max="1500"
          step="50"
          value={speed}
          disabled={disabled}
          onChange={(event) => onSpeedChange(Number(event.target.value))}
        />
        <span>{speed} ms</span>
      </label>
    </section>
  );
}

export default PlaybackControls;
