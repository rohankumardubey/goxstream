package source

import (
    "context"
    "encoding/json"
    "goxstream/internal/model"
    "github.com/segmentio/kafka-go"
    "time"
)

type KafkaSourceConfig struct {
    Brokers []string
    Topic   string
    GroupID string
}

func KafkaSource(cfg KafkaSourceConfig, out chan<- model.Event) error {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  cfg.Brokers,
        GroupID:  cfg.GroupID,
        Topic:    cfg.Topic,
        MinBytes: 1e3,
        MaxBytes: 1e6,
    })
    defer r.Close()
    ctx := context.Background()

    for {
        m, err := r.ReadMessage(ctx)
        if err != nil {
            break
        }
        var data map[string]interface{}
        if err := json.Unmarshal(m.Value, &data); err != nil {
            continue
        }
        out <- model.Event{
            Data:      data,
            Timestamp: time.Now(), // or extract from data["timestamp"]
        }
    }
    close(out)
    return nil
}
