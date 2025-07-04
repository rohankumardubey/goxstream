package sink

import (
    "encoding/csv"
    "os"
    "goxstream/internal/model"
    "sort"
	"fmt"
)

func FileSink(path string, in <-chan model.Event) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    writer := csv.NewWriter(f)
    var headers []string
    var headersWritten bool
    for event := range in {
        if !headersWritten {
            // Consistent header order (sort for now; you can use custom order if needed)
            for k := range event.Data {
                headers = append(headers, k)
            }
            sort.Strings(headers) // optional, or use custom order
            writer.Write(headers)
            headersWritten = true
        }
        row := make([]string, len(headers))
        for i, k := range headers {
            if val, ok := event.Data[k]; ok {
                row[i] = toString(val)
            }
        }
        writer.Write(row)
    }
    writer.Flush()
    return writer.Error()
}

// Helper to stringify interface{} to string
func toString(val interface{}) string {
    if s, ok := val.(string); ok {
        return s
    }
    return fmt.Sprintf("%v", val)
}
