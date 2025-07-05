import React, { useCallback, useState } from "react";
import ReactFlow, {
  MiniMap, Controls, Background,
  useNodesState, useEdgesState
} from "reactflow";
import "reactflow/dist/style.css";
import {
  Box, Button, Paper, Stack, Typography, Dialog, DialogTitle, DialogContent,
  DialogActions, TextField, MenuItem
} from "@mui/material";

const operatorFields = {
  map: [
    { name: "col", label: "Column", type: "text" },
    { name: "val", label: "Value", type: "text" }
  ],
  filter: [
    { name: "field", label: "Field", type: "text" },
    { name: "eq", label: "Equals", type: "text" }
  ],
  reduce: [
    { name: "key", label: "Key", type: "text" },
    { name: "agg", label: "Aggregation", type: "select", options: ["count", "sum", "avg"] }
  ],
  // Add more as needed
};

const initialNodes = [
  {
    id: '1',
    type: 'input',
    position: { x: 0, y: 100 },
    data: { label: 'Source: File', type: "source", params: { type: "file", path: "input.csv" } },
  },
  {
    id: '2',
    position: { x: 250, y: 100 },
    data: { label: 'Map Operator', type: "map", params: { col: "processed", val: "yes" } },
  },
  {
    id: '3',
    type: 'output',
    position: { x: 500, y: 100 },
    data: { label: 'Sink: File', type: "sink", params: { type: "file", path: "output.csv" } },
  },
];

const initialEdges = [
  { id: 'e1-2', source: '1', target: '2', animated: true },
  { id: 'e2-3', source: '2', target: '3', animated: true },
];

export default function VisualDesigner() {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  // Modal state
  const [editNode, setEditNode] = useState(null);
  const [form, setForm] = useState({});

  // Add new operator node
  const addOperator = useCallback((type, label) => {
    const newId = (nodes.length + 1).toString();
    const lastX = nodes.length * 220;
    const params = {};
    (operatorFields[type] || []).forEach(f => { params[f.name] = ""; });
    setNodes((nds) => [
      ...nds,
      {
        id: newId,
        position: { x: lastX, y: 100 },
        data: { label, type, params },
      },
    ]);
    // Auto-link from previous node
    if (nodes.length > 0) {
      setEdges((eds) => [
        ...eds,
        { id: `e${nodes.length}-${newId}`, source: nodes[nodes.length - 1].id, target: newId, animated: true }
      ]);
    }
  }, [nodes, setNodes, setEdges]);

  // Node click handler
  const onNodeClick = useCallback((evt, node) => {
    if (node.data.type === "source" || node.data.type === "sink") return; // skip editing
    setEditNode(node);
    setForm({ ...node.data.params });
  }, []);

  // Handle form changes in modal
  const handleFormChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  // Save node edits
  const handleSave = () => {
    setNodes(nds => nds.map(n => n.id === editNode.id
      ? {
          ...n,
          data: {
            ...n.data,
            params: { ...form },
            label: getNodeLabel(n.data.type, form),
          }
        }
      : n));
    setEditNode(null);
  };

  const getNodeLabel = (type, params) => {
    if (type === "map") return `Map: ${params.col || ""} = ${params.val || ""}`;
    if (type === "filter") return `Filter: ${params.field || ""} == ${params.eq || ""}`;
    if (type === "reduce") return `Reduce by ${params.key || ""} (${params.agg || ""})`;
    return type.charAt(0).toUpperCase() + type.slice(1) + " Operator";
  };

  // Export current graph to JSON pipeline spec
  const exportPipeline = () => {
    // Build pipeline spec from nodes/params
    const ops = nodes
      .filter(n => n.data.type !== "source" && n.data.type !== "sink")
      .map(n => ({
        type: n.data.type,
        params: { ...n.data.params }
      }));
    const spec = {
      source: nodes.find(n => n.data.type === "source")?.data.params || { type: "file", path: "input.csv" },
      operators: ops,
      sink: nodes.find(n => n.data.type === "sink")?.data.params || { type: "file", path: "output.csv" }
    };
    alert(JSON.stringify(spec, null, 2));
  };

  return (
    <Box sx={{ height: 500, width: "100%" }}>
      <Paper sx={{ mb: 2, p: 2 }}>
        <Stack direction="row" spacing={2} alignItems="center">
          <Typography variant="h5" fontWeight={600}>Visual Pipeline Designer</Typography>
          <Button variant="outlined" onClick={() => addOperator("map", "Map Operator")}>Add Map</Button>
          <Button variant="outlined" onClick={() => addOperator("filter", "Filter Operator")}>Add Filter</Button>
          <Button variant="outlined" onClick={() => addOperator("reduce", "Reduce Operator")}>Add Reduce</Button>
          <Button variant="contained" color="success" onClick={exportPipeline}>Export Pipeline JSON</Button>
        </Stack>
      </Paper>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={onNodeClick}
        fitView
      >
        <MiniMap />
        <Controls />
        <Background color="#bbb" gap={18} />
      </ReactFlow>

      {/* --- Modal Dialog for Editing Node --- */}
      <Dialog open={!!editNode} onClose={() => setEditNode(null)}>
        <DialogTitle>Edit Operator</DialogTitle>
        <DialogContent>
          {editNode && operatorFields[editNode.data.type] ? (
            operatorFields[editNode.data.type].map(field => (
              field.type === "select" ? (
                <TextField
                  select
                  fullWidth
                  margin="dense"
                  key={field.name}
                  name={field.name}
                  label={field.label}
                  value={form[field.name] || ""}
                  onChange={handleFormChange}
                >
                  {field.options.map(opt =>
                    <MenuItem value={opt} key={opt}>{opt}</MenuItem>
                  )}
                </TextField>
              ) : (
                <TextField
                  key={field.name}
                  margin="dense"
                  name={field.name}
                  label={field.label}
                  type={field.type}
                  fullWidth
                  value={form[field.name] || ""}
                  onChange={handleFormChange}
                />
              )
            ))
          ) : (
            <Typography>No fields for this operator.</Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditNode(null)}>Cancel</Button>
          <Button onClick={handleSave} variant="contained">Save</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
