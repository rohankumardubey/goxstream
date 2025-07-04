package source

import (
    "bufio"
    "encoding/csv"
    "os"
    "strings"
    "time"
    "goxstream/internal/model"
)

func FileSource(path string, out chan<- model.Event) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()

    reader := csv.NewReader(bufio.NewReader(f))
    headers, err := reader.Read()
    if err != nil {
        return err
    }

    // Find index of "timestamp" column, case-insensitive
    timestampIdx := -1
    for i, h := range headers {
        if strings.EqualFold(h, "timestamp") {
            timestampIdx = i
            break
        }
    }

    for {
        record, err := reader.Read()
        if err != nil {
            break
        }
        data := map[string]interface{}{}
        for i, h := range headers {
            data[h] = record[i]
        }

        var evtTime time.Time
        if timestampIdx != -1 {
            // Try to parse as RFC3339, else fallback to time.Now()
            if t, err := time.Parse(time.RFC3339, record[timestampIdx]); err == nil {
                evtTime = t
            } else {
                evtTime = time.Now()
            }
        } else {
            evtTime = time.Now()
        }

        out <- model.Event{
            Data:      data,
            Timestamp: evtTime,
        }
    }
    close(out)
    return nil
}
