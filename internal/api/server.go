package api

import (
    "encoding/json"
    "net/http"
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

    var spec model.PipelineSpec
    if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    go func() {
        engine.BuildAndRunPipeline(spec)
    }()

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(`{"status":"started"}`))
}
