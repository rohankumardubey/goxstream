package source

import (
    "context"
    "database/sql"
    "fmt"
    "goxstream/internal/model"
    _ "github.com/lib/pq" // or your preferred DB driver
    "time"
)

type DBSourceConfig struct {
    DSN   string
    Query string
}

func DBSource(cfg DBSourceConfig, out chan<- model.Event) error {
    db, err := sql.Open("postgres", cfg.DSN)
    if err != nil { return fmt.Errorf("open db: %w", err) }
    defer db.Close()

    rows, err := db.QueryContext(context.Background(), cfg.Query)
    if err != nil { return fmt.Errorf("query: %w", err) }
    defer rows.Close()

    cols, _ := rows.Columns()
    for rows.Next() {
        values := make([]interface{}, len(cols))
        ptrs := make([]interface{}, len(cols))
        for i := range values {
            ptrs[i] = &values[i]
        }
        if err := rows.Scan(ptrs...); err != nil {
            continue
        }
        data := map[string]interface{}{}
        for i, col := range cols {
            val := values[i]
            // Optionally: Parse time columns for Timestamp
            data[col] = val
        }
        // Use time.Now() or extract a timestamp column if available
        out <- model.Event{Data: data, Timestamp: time.Now()}
    }
    close(out)
    return nil
}
