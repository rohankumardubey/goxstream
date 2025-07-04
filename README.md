# GoXStream

**GoXStream** is a modern, extensible real-time streaming engine built in Go, inspired by Apache Flink.
It enables processing and transforming streaming data with powerful, configurable pipelines using operators like map, filter, reduce, and both tumbling and sliding windows.
GoXStream supports dynamic job definitions via REST API and is ready for integration with a React UI dashboard.

---

## ✨ Features

## Features

- **Dynamic pipeline definition via REST API**
- **Event-time tumbling & sliding windows**
- **Watermark/late event support** (true stream semantics)
- **Chaining map, filter, reduce operators**
- **CSV file source/sink (DB/Kafka coming soon)**
- **Ready for React UI dashboard integration**
- **Easy extension: Add custom operators and connectors**

---

## 🚀 Quick Start

### 1. **Clone and Build**

```bash
git clone https://github.com/YOUR_GITHUB_USERNAME/goxstream.git
cd goxstream
go mod tidy
```
### 2. Prepare Input Data
***Place an example input.csv in the project root:***

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

### 4. Submit a Pipeline Job
***Use curl or Postman to submit a dynamic pipeline (example: sliding window reduce):***
***a. Regular time window:***


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
***Add New Operators:
Implement the Operator interface and add a factory to the operator registry.***

***Support New Sources/Sinks:
Implement a Source or Sink interface in internal/source or internal/sink.***

***React UI Integration:
Planned for interactive pipeline creation and monitoring.***

---

### 🔜 Roadmap
 [X] Count-based tumbling/sliding windows

 [X] Time-based tumbling/sliding windows

 [X] Watermark & late event support

 [ ] DB & Kafka sources/sinks

 [ ] React UI dashboard

 [ ] More aggregations: sum, avg, min, max

 [ ] State, session windows, custom UDFs

---


 ### 🙌 Contributing
***PRs, issues, and ideas are welcome!
Fork and submit improvements or new features—let’s build a great Go stream engine together!***


### GoXStream — Streaming, the Go way! 🚀