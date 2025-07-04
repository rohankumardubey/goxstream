package source

import (
    "bufio"
    "encoding/csv"
    "os"
    "time"
    "goxstream/internal/model"
)

func FileSource(path string, out chan<- model.Event) error {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close()
    reader := csv.NewReader(bufio.NewReader(f))
    headers, err := reader.Read()
    if err != nil { return err }
    for {
        record, err := reader.Read()
        if err != nil { break }
        data := map[string]interface{}{}
        for i, h := range headers {
            data[h] = record[i]
        }
        out <- model.Event{
            Data: data,
            Timestamp: time.Now(), // Or parse from record
        }
    }
    close(out)
    return nil
}
