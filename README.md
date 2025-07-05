# GoXStream

GoXStream is an open-source, Flink-inspired **stream processing engine written in Go**—with a beautiful React dashboard for visual pipeline design, job submission, and history.
Build and run real-time data pipelines with file, database, or Kafka sources and sinks.

![GoXStream Dashboard Screenshot](./docs/screenshot-dashboard.png)
![Visual Designer Screenshot](./docs/screenshot-designer.png)

---

## 🚀 Features

- **Modular pipeline engine (in Go):** Compose pipelines from map, filter, reduce, window, time-window, and more
- **Dynamic REST API:** Submit pipelines and configure sources, sinks, operators via JSON
- **Pluggable sources/sinks:** File, Postgres, Kafka (more coming)
- **Windowing:** Tumbling, sliding, time-based, with watermark and late event support
- **Stateful operators and checkpointing** (coming soon!)
- **React dashboard:** Visual DAG pipeline builder (drag/drop), job submission, job history, JSON preview
- **Persistent job history (localStorage and soon, backend)**

---

## 🚀 Quick Start

### 1. **Clone and Build**

```bash
git clone https://github.com/YOUR_GITHUB_USERNAME/goxstream.git
cd goxstream
go mod tidy
```

### 2. Prepare Input Data

***Place an example input.csv in the project root for the :***

```bash
id,name,city
1,Alice,London
2,Bob,Berlin
3,Charlie,Paris
4,David,Berlin
5,Eve,Paris
6,Frank,Paris
7,Grace,Berlin
```

### 3. Run the API Server

```bash
go run ./cmd/goxstream/main.go
```

### 4. Open the React dashboard

```bash
cd goxstream-dashboard
npm install
npm start
```

### 5. 🧩 Example: Submit a Pipeline via UI or REST

_**Simple Map:**_

```bash
{
  "source": { "type": "file", "path": "input.csv" },
  "operators": [
    { "type": "map", "params": { "col": "processed", "val": "yes" } }
  ],
  "sink": { "type": "file", "path": "output.csv" }
}
```

_**Windowed Reduce:**_

```bash
{
  "source": { "type": "file", "path": "input.csv" },
  "operators": [
    {
      "type": "time_window",
      "params": {
        "duration": "10s",
        "inner": {
          "type": "reduce",
          "params": { "key": "city", "agg": "count" }
        }
      }
    }
  ],
  "sink": { "type": "file", "path": "output.csv" }
}
```

_**With Watermark/Late Event Support:**_

```bash
{
  "source": { "type": "file", "path": "input.csv" },
  "operators": [
    {
      "type": "time_window_watermark",
      "params": {
        "duration": "10s",
        "allowed_lateness": "5s",
        "inner": {
          "type": "reduce",
          "params": { "key": "city", "agg": "count" }
        }
      }
    }
  ],
  "sink": { "type": "file", "path": "output.csv" }
}
```
---

### 🖥️ Visual Pipeline Designer (UI)
- GoXStream’s UI lets you visually build pipelines (drag/drop), edit operator parameters, and export as JSON to run jobs.

- Submitted jobs and their configs are saved in a beautiful job history.

---

### 6. Submit a Pipeline Job (file → file example) via CURL

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "source": { "type": "file", "path": "input.csv" },
    "operators": [{ "type": "map", "params": { "col": "processed", "val": "yes" } }],
    "sink": { "type": "file", "path": "output.csv" }
  }'

```

---

### 7. Submit a Pipeline Job

***Use curl or Postman to submit a dynamic pipeline (example: sliding window reduce):***
***a. Regular time window:***

***Input.csv file for both below time_window and time_window_watermark***

```bash
id,name,city,score,timestamp
1,Alice,London,10,2024-07-04T15:00:00Z
2,Bob,Berlin,15,2024-07-04T15:00:05Z
3,Charlie,Paris,12,2024-07-04T15:00:08Z
4,David,Berlin,8,2024-07-04T15:00:12Z
5,Eve,Paris,14,2024-07-04T15:00:16Z
6,Frank,London,11,2024-07-04T15:00:22Z
7,Grace,Berlin,9,2024-07-04T15:00:29Z
8,Harry,Paris,13,2024-07-04T15:00:35Z
9,Ivy,London,17,2024-07-04T15:00:41Z
10,Jack,Berlin,16,2024-07-04T15:00:43Z
```

***a. Time window with reduce event support:***

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "source": { "type": "file", "path": "input.csv" },
    "operators": [
      {
        "type": "time_window",
        "params": {
          "duration": "10s",
          "inner": {
            "type": "reduce",
            "params": { "key": "city", "agg": "count" }
          }
        }
      }
    ],
    "sink": { "type": "file", "path": "output.csv" }
  }'
```

***b. Time window with watermark/late event support:***

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "source": { "type": "file", "path": "input.csv" },
    "operators": [
      {
        "type": "time_window_watermark",
        "params": {
          "duration": "10s",
          "allowed_lateness": "5s",
          "inner": {
            "type": "reduce",
            "params": { "key": "city", "agg": "count" }
          }
        }
      }
    ],
    "sink": { "type": "file", "path": "output.csv" }
  }'
```

***Your results will be in output.csv with a window_id column.***

```bash
city,count,window_end,window_id
London,1,2024-07-04T15:00:10Z,1
Berlin,1,2024-07-04T15:00:10Z,1
Paris,1,2024-07-04T15:00:10Z,1
...
```

---

### 🛠️ Architecture

```bash
[Source] --> [Map] --> [Filter] --> [Window/Reduce] --> [Sink]
   |            |         |           |                   |
 [File]   [Add/Transform] [Select] [Tumbling/Sliding]  [File]

```

***Operators: Implemented as Go interfaces, dynamically composed at runtime.***

***REST API: Accepts JSON job specs, launches pipeline as background Go routines.***

---

### 📝 JSON Job Spec

***A pipeline is defined by a simple JSON:***

```bash
{
  "source": {"type": "file", "path": "input.csv"},
  "operators": [
    {"type": "map", "params": {"col": "processed", "val": "yes"}},
    {"type": "filter", "params": {"field": "city", "eq": "Berlin"}},
    {
      "type": "sliding_window",
      "params": {
        "size": 3,
        "step": 1,
        "inner": {
          "type": "reduce",
          "params": {"key": "city", "agg": "count"}
        }
      }
    }
  ],
  "sink": {"type": "file", "path": "output.csv"}
}
```

---

### 📚 Operator Types

```bash
| Type             | Description              | Example Params                        |
| ---------------- | ------------------------ | ------------------------------------- |
| map              | Add or transform columns | `col`, `val`                          |
| filter           | Filter rows by condition | `field`, `eq`                         |
| reduce           | Aggregate/group by field | `key`, `agg` (`count`, future: `sum`) |
| tumbling\_window | Non-overlapping windows  | `size`, `inner`                       |
| sliding\_window  | Overlapping windows      | `size`, `step`, `inner`               |
```

---

### 🧑‍💻 Extending GoXStream

_Add New Operators:

Implement the Operator interface and add a factory to the operator registry._

_Support New Sources/Sinks:

Implement a Source or Sink interface in internal/source or internal/sink._

_React UI Integration:

Planned for interactive pipeline creation and monitoring._

---

### 🔜 Roadmap

- [x] Dynamic operator/source/sink registry

- [x] File/DB/Kafka connectors

- [x] REST API & React UI for pipeline design and monitoring

- [x] Visual DAG editor with drag/drop (reactflow)

- [x] Windowing and watermark support

- [x] Persistent job history (localStorage)

- [ ] Checkpoint and fault-tolerance (coming next)

- [ ] Backend job monitoring/status APIs

- [ ] Multi-job/cluster execution

- [ ] More analytics and ML operators

---

### 🙌 Contributing
***PRs, issues, and ideas are welcome!
Fork and submit improvements or new features—let’s build a great Go stream engine together!***

### GoXStream — Streaming, the Go way! 🚀