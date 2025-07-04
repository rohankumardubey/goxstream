package sink

import (
    "context"
    "encoding/json"
    "goxstream/internal/model"
    "github.com/segmentio/kafka-go"
	"fmt"
)

type KafkaSinkConfig struct {
    Brokers []string
    Topic   string
}

func KafkaSink(cfg KafkaSinkConfig, in <-chan model.Event) error {
    w := kafka.NewWriter(kafka.WriterConfig{
        Brokers: cfg.Brokers,
        Topic:   cfg.Topic,
    })
    defer w.Close()

    ctx := context.Background()
    for event := range in {
        data, err := json.Marshal(event.Data)
        if err != nil {
            continue
        }
        msg := kafka.Message{Value: data}
        if err := w.WriteMessages(ctx, msg); err != nil {
            fmt.Println("kafka write error:", err)
            continue
        }
    }
    return nil
}
