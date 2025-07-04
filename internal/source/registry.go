package source

import (
	"fmt"
	"goxstream/internal/model"
)

// -------- Source Registry --------

type SourceFactory func(params map[string]interface{}, out chan<- model.Event) error

var registry = map[string]SourceFactory{
	"file":  fileSourceFactory,
	"db":    dbSourceFactory,
	"kafka": kafkaSourceFactory,
}

// BuildSource dynamically constructs the source based on JSON spec
func BuildSource(srcSpec map[string]interface{}, out chan<- model.Event) error {
	srcType, ok := srcSpec["type"].(string)
	if !ok {
		return fmt.Errorf("source missing 'type'")
	}
	factory, ok := registry[srcType]
	if !ok {
		return fmt.Errorf("unknown source type: %s", srcType)
	}
	return factory(srcSpec, out)
}

// -------- Adapters for each source type --------

// File source expects: { "type": "file", "path": "input.csv" }
func fileSourceFactory(params map[string]interface{}, out chan<- model.Event) error {
	path, ok := params["path"].(string)
	if !ok {
		return fmt.Errorf("file source expects 'path'")
	}
	return FileSource(path, out)
}

// DB source expects: { "type": "db", "dsn": "...", "query": "..." }
func dbSourceFactory(params map[string]interface{}, out chan<- model.Event) error {
	dsn, ok1 := params["dsn"].(string)
	query, ok2 := params["query"].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("db source expects 'dsn' and 'query'")
	}
	return DBSource(DBSourceConfig{DSN: dsn, Query: query}, out)
}

// Kafka source expects: { "type": "kafka", "brokers": [...], "topic": "...", "group_id": "..." }
func kafkaSourceFactory(params map[string]interface{}, out chan<- model.Event) error {
	brokersIface, ok1 := params["brokers"].([]interface{})
	topic, ok2 := params["topic"].(string)
	groupID, ok3 := params["group_id"].(string)
	if !ok1 || !ok2 || !ok3 {
		return fmt.Errorf("kafka source expects 'brokers', 'topic', 'group_id'")
	}
	brokers := make([]string, len(brokersIface))
	for i, v := range brokersIface {
		str, ok := v.(string)
		if !ok {
			return fmt.Errorf("kafka source: all brokers must be strings")
		}
		brokers[i] = str
	}
	return KafkaSource(KafkaSourceConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	}, out)
}
