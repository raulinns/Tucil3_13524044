import { useMemo, useState } from "react";

const nodeWidth = 106;
const nodeHeight = 54;
const columnGap = 42;
const rowGap = 30;

function GraphView({ trace, activeTraceIndex }) {
  const [showAll, setShowAll] = useState(false);
  const graph = useMemo(() => buildLayout(trace, activeTraceIndex, showAll), [
    activeTraceIndex,
    showAll,
    trace,
  ]);

  if (!trace?.nodes?.length) {
    return (
      <section className="graph-panel">
        <div className="section-heading">
          <h2>Search Graph</h2>
          <span>0 nodes</span>
        </div>
        <div className="empty-state compact">Graph trace will appear after solving.</div>
      </section>
    );
  }

  return (
    <section className="graph-panel">
      <div className="section-heading">
        <h2>Search Graph</h2>
        <div className="graph-actions">
          <span>
            {graph.nodes.length} / {trace.nodes.length} nodes
          </span>
          <button type="button" onClick={() => setShowAll((value) => !value)}>
            {showAll ? "Limit" : "Show all"}
          </button>
        </div>
      </div>

      <div className="graph-canvas" style={{ height: graph.height }}>
        <svg
          className="edge-layer"
          width={graph.width}
          height={graph.height}
          viewBox={`0 0 ${graph.width} ${graph.height}`}
          role="img"
          aria-label="Search trace graph"
        >
          <defs>
            <marker id="arrow" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
              <path d="M0,0 L0,6 L9,3 z" />
            </marker>
          </defs>
          {graph.edges.map((edge) => {
            const from = graph.positions.get(edge.from);
            const to = graph.positions.get(edge.to);
            if (!from || !to) {
              return null;
            }
            const x1 = from.x + nodeWidth;
            const y1 = from.y + nodeHeight / 2;
            const x2 = to.x;
            const y2 = to.y + nodeHeight / 2;
            const midX = (x1 + x2) / 2;
            const d = `M ${x1} ${y1} C ${midX} ${y1}, ${midX} ${y2}, ${x2} ${y2}`;

            return (
              <g key={edge.id} className={edgeClass(edge)}>
                <path d={d} markerEnd="url(#arrow)" />
                <text x={midX} y={(y1 + y2) / 2 - 6}>
                  {edge.direction}:{edge.stepCost}
                </text>
              </g>
            );
          })}
        </svg>

        {graph.nodes.map((node) => (
          <div
            key={node.id}
            className={nodeClass(node, activeTraceIndex)}
            style={{
              left: graph.positions.get(node.id).x,
              top: graph.positions.get(node.id).y,
            }}
            title={`g=${node.gCost}, h=${node.hCost}, f=${node.priority}`}
          >
            <strong>{node.label}</strong>
            <span>
              g {node.gCost} | h {node.hCost} | p {node.priority}
            </span>
          </div>
        ))}
      </div>

      <div className="graph-legend">
        <span><i className="legend-dot expanded" /> expanded</span>
        <span><i className="legend-dot solution" /> solution</span>
        <span><i className="legend-line accepted" /> accepted edge</span>
        <span><i className="legend-line rejected" /> rejected edge</span>
      </div>
    </section>
  );
}

function buildLayout(trace, activeTraceIndex, showAll) {
  const allNodes = trace?.nodes || [];
  const allEdges = trace?.edges || [];
  const visibleNodeIDs = new Set();
  const maxNodes = showAll ? allNodes.length : 80;

  const sortedNodes = [...allNodes].sort((a, b) => {
    const ax = a.expandedAt < 0 ? Number.MAX_SAFE_INTEGER : a.expandedAt;
    const bx = b.expandedAt < 0 ? Number.MAX_SAFE_INTEGER : b.expandedAt;
    return ax - bx || a.id.localeCompare(b.id);
  });

  for (const node of sortedNodes) {
    const isActive = activeTraceIndex < 0 || node.expandedAt <= activeTraceIndex || node.solution;
    if ((showAll || isActive) && visibleNodeIDs.size < maxNodes) {
      visibleNodeIDs.add(node.id);
    }
  }

  if (!visibleNodeIDs.size) {
    for (const node of sortedNodes.slice(0, maxNodes)) {
      visibleNodeIDs.add(node.id);
    }
  }

  for (const edge of allEdges) {
    if (visibleNodeIDs.has(edge.from) && visibleNodeIDs.size < maxNodes) {
      visibleNodeIDs.add(edge.to);
    }
    if (visibleNodeIDs.has(edge.to) && visibleNodeIDs.size < maxNodes) {
      visibleNodeIDs.add(edge.from);
    }
  }

  const nodes = sortedNodes.filter((node) => visibleNodeIDs.has(node.id));
  const edges = allEdges.filter((edge) => visibleNodeIDs.has(edge.from) && visibleNodeIDs.has(edge.to));

  const columns = Math.min(6, Math.max(1, Math.ceil(Math.sqrt(nodes.length))));
  const positions = new Map();
  nodes.forEach((node, index) => {
    const col = index % columns;
    const row = Math.floor(index / columns);
    positions.set(node.id, {
      x: 24 + col * (nodeWidth + columnGap),
      y: 24 + row * (nodeHeight + rowGap),
    });
  });

  const rows = Math.max(1, Math.ceil(nodes.length / columns));
  const width = 48 + columns * nodeWidth + (columns - 1) * columnGap;
  const height = 48 + rows * nodeHeight + (rows - 1) * rowGap;

  return { edges, height, nodes, positions, width };
}

function nodeClass(node, activeTraceIndex) {
  const classes = ["graph-node"];
  if (node.solution) {
    classes.push("solution");
  }
  if (node.expandedAt >= 0) {
    classes.push("expanded");
  }
  if (node.expandedAt === activeTraceIndex) {
    classes.push("active");
  }
  return classes.join(" ");
}

function edgeClass(edge) {
  const classes = ["graph-edge"];
  if (edge.accepted) {
    classes.push("accepted");
  } else {
    classes.push("rejected");
  }
  if (edge.solution) {
    classes.push("solution");
  }
  return classes.join(" ");
}

export default GraphView;
