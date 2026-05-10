const tileNames = {
  X: "wall",
  "*": "ice",
  L: "lava",
  Z: "player",
  O: "goal",
};

function Board({ grid }) {
  if (!grid.length) {
    return (
      <div className="empty-state">
        Upload a testcase and run a solver.
      </div>
    );
  }

  const columns = grid[0]?.length || 0;

  return (
    <div
      className="board"
      style={{ "--columns": columns }}
      aria-label="Puzzle board"
    >
      {grid.flatMap((row, rowIndex) =>
        [...row].map((tile, colIndex) => {
          const kind = tileNames[tile] || (/\d/.test(tile) ? "checkpoint" : "ice");
          return (
            <div
              className={`tile tile-${kind}`}
              key={`${rowIndex}-${colIndex}`}
              title={`(${rowIndex}, ${colIndex}) ${tile}`}
            >
              {tile}
            </div>
          );
        }),
      )}
    </div>
  );
}

export default Board;
