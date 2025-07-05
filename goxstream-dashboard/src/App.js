import React, { useState } from "react";
import VisualDesigner from "./VisualDesigner";
import {
  Box,
  Container,
  CssBaseline,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Paper,
  Typography,
  Button,
  Tabs,
  Tab,
  TextField,
  Card,
  CardContent,
  Chip,
  Divider,
  Stack
} from "@mui/material";
import DashboardIcon from "@mui/icons-material/Dashboard";
import DesignServicesIcon from "@mui/icons-material/DesignServices";
import HistoryIcon from "@mui/icons-material/History";
import AddTaskIcon from "@mui/icons-material/AddTask";

const drawerWidth = 210;

// --- PipelineSubmit component ---
function PipelineSubmit({ onJobSubmit }) {
  const [json, setJson] = useState(`{
  "source": { "type": "file", "path": "input.csv" },
  "operators": [
    { "type": "map", "params": { "col": "processed", "val": "yes" } }
  ],
  "sink": { "type": "file", "path": "output.csv" }
}`);
  const [status, setStatus] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    setSubmitting(true);
    setStatus("Submitting...");
    try {
      const resp = await fetch("http://localhost:8080/jobs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: json,
      });
      if (resp.ok) {
        setStatus("✅ Job submitted!");
        onJobSubmit({
          ...JSON.parse(json),
          submitted: new Date().toISOString(),
        });
      } else {
        setStatus("❌ Error: " + (await resp.text()));
      }
    } catch (err) {
      setStatus("❌ Error: " + err.message);
    }
    setSubmitting(false);
  }

  return (
    <Paper sx={{ p: 3, mb: 4 }}>
      <Stack direction="row" alignItems="center" spacing={2} mb={2}>
        <AddTaskIcon color="primary" />
        <Typography variant="h5" fontWeight={600}>
          Submit Pipeline Job
        </Typography>
      </Stack>
      <form onSubmit={handleSubmit}>
        <TextField
          label="Pipeline JSON"
          multiline
          minRows={10}
          maxRows={20}
          fullWidth
          variant="outlined"
          value={json}
          onChange={e => setJson(e.target.value)}
          sx={{
            fontFamily: "monospace",
            fontSize: 14,
            background: "#fcfcfc",
            mb: 2,
          }}
        />
        <Stack direction="row" spacing={2} alignItems="center" mt={1}>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            disabled={submitting}
          >
            Submit
          </Button>
          <Typography variant="body1">{status}</Typography>
        </Stack>
      </form>
    </Paper>
  );
}

// --- JobHistory component (with pretty cards & chips) ---
function JobHistory({ jobs }) {
  if (jobs.length === 0)
    return <Typography color="text.secondary">No jobs submitted yet.</Typography>;

  return (
    <Box sx={{ mt: 2 }}>
      <Stack direction="row" alignItems="center" spacing={2} mb={2}>
        <HistoryIcon color="primary" />
        <Typography variant="h5" fontWeight={600}>
          Job History
        </Typography>
      </Stack>
      <Stack spacing={3}>
        {jobs.map((job, i) => (
          <Card key={i} sx={{ background: "#fafaff", borderLeft: "6px solid #2196f3" }}>
            <CardContent>
              <Stack
                direction="row"
                spacing={2}
                alignItems="center"
                justifyContent="space-between"
              >
                <Typography variant="subtitle1" fontWeight={600}>
                  Job #{jobs.length - i}
                </Typography>
                <Chip
                  label={job.source?.type?.toUpperCase() || "UNKNOWN"}
                  color="primary"
                  size="small"
                  sx={{ ml: 1 }}
                />
                <Chip
                  label={job.sink?.type?.toUpperCase() || "UNKNOWN"}
                  color="success"
                  size="small"
                />
                <Typography variant="caption" color="text.secondary">
                  {job.submitted
                    ? new Date(job.submitted).toLocaleString()
                    : ""}
                </Typography>
              </Stack>
              <Divider sx={{ my: 1.2 }} />
              <pre
                style={{
                  background: "#f6f8fa",
                  padding: 12,
                  borderRadius: 8,
                  fontFamily: "monospace",
                  fontSize: 14,
                  overflow: "auto",
                  margin: 0,
                }}
              >
                {JSON.stringify(job, null, 2)}
              </pre>
            </CardContent>
          </Card>
        ))}
      </Stack>
    </Box>
  );
}

// --- Main App component ---
function App() {
  const [page, setPage] = useState("dashboard");
  const [jobHistory, setJobHistory] = useState(() => {
    const saved = localStorage.getItem("goxstreamJobs");
    return saved ? JSON.parse(saved) : [];
  });

  function handleJobSubmit(job) {
    const newHistory = [job, ...jobHistory];
    setJobHistory(newHistory);
    localStorage.setItem("goxstreamJobs", JSON.stringify(newHistory));
  }

  // For sidebar
  const navItems = [
    {
      label: "Dashboard",
      icon: <DashboardIcon color={page === "dashboard" ? "primary" : "action"} />,
      value: "dashboard",
    },
    {
      label: "Visual Designer",
      icon: <DesignServicesIcon color={page === "designer" ? "primary" : "action"} />,
      value: "designer",
    },
    {
      label: "Job History",
      icon: <HistoryIcon color={page === "history" ? "primary" : "action"} />,
      value: "history",
    },
  ];

  return (
    <Box sx={{ display: "flex" }}>
      <CssBaseline />
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: {
            width: drawerWidth,
            boxSizing: "border-box",
            background: "#f5f7fb",
            borderRight: "1px solid #e0e7ef",
          },
        }}
      >
        <Box sx={{ display: "flex", alignItems: "center", px: 2, py: 2 }}>
          <img src="https://cdn.jsdelivr.net/gh/gilbarbara/logos@master/logos/go.svg" alt="Go" width={36} style={{ marginRight: 12 }} />
          <Typography variant="h6" fontWeight={700}>
            GoXStream
          </Typography>
        </Box>
        <Divider />
        <List>
          {navItems.map(item => (
            <ListItem key={item.value} disablePadding>
              <ListItemButton
                selected={page === item.value}
                onClick={() => setPage(item.value)}
              >
                <ListItemIcon>{item.icon}</ListItemIcon>
                <ListItemText primary={item.label} />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </Drawer>

      <Container
        maxWidth="md"
        sx={{
          flexGrow: 1,
          mt: 5,
          mb: 6,
          ml: `${drawerWidth}px`,
          minHeight: "100vh",
        }}
      >
        <Box sx={{ mb: 4 }}>
          {page === "dashboard" && (
            <>
              <PipelineSubmit onJobSubmit={handleJobSubmit} />
              <JobHistory jobs={jobHistory} />
            </>
          )}
          {page === "designer" && <VisualDesigner />}
          {page === "history" && <JobHistory jobs={jobHistory} />}
        </Box>
      </Container>
    </Box>
  );
}

export default App;
