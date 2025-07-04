package sink

import (
    "context"
    "database/sql"
    "fmt"
    "goxstream/internal/model"
    _ "github.com/lib/pq"
)

type DBSinkConfig struct {
    DSN   string
    Table string
}

func DBSink(cfg DBSinkConfig, in <-chan model.Event) error {
    db, err := sql.Open("postgres", cfg.DSN)
    if err != nil { return fmt.Errorf("open db: %w", err) }
    defer db.Close()

    for event := range in {
        // For simplicity, write only JSON-encoded Data
        _, err := db.ExecContext(context.Background(),
            fmt.Sprintf("INSERT INTO %s (data) VALUES ($1)", cfg.Table),
            event.Data)
        if err != nil {
            fmt.Println("db sink error:", err)
            continue
        }
    }
    return nil
}
