package api

import (
    "encoding/json"
    "net/http"
    "io"
    "goxstream/internal/model"
    "goxstream/internal/engine"
)

func StartAPIServer(addr string) error {
    http.HandleFunc("/jobs", jobHandler)
    return http.ListenAndServe(addr, nil)
}

func jobHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Read the entire body ONCE
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "could not read body", http.StatusBadRequest)
        return
    }

    // Unmarshal into struct
    var spec model.PipelineSpec
    if err := json.Unmarshal(bodyBytes, &spec); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    // Unmarshal into map for Raw fields
    var pipelineMap map[string]interface{}
    if err := json.Unmarshal(bodyBytes, &pipelineMap); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    if src, ok := pipelineMap["source"].(map[string]interface{}); ok {
        spec.Source.Raw = src
    } else {
        http.Error(w, "missing source in pipeline", http.StatusBadRequest)
        return
    }
    if sink, ok := pipelineMap["sink"].(map[string]interface{}); ok {
        spec.Sink.Raw = sink
    } else {
        http.Error(w, "missing sink in pipeline", http.StatusBadRequest)
        return
    }

    go func() {
        engine.BuildAndRunPipeline(spec)
    }()

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(`{"status":"started"}`))
}
