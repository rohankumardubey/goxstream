package sink

import (
	"fmt"
	"goxstream/internal/model"
)

// -------- Sink Registry --------

type SinkFactory func(params map[string]interface{}, in <-chan model.Event) error

var registry = map[string]SinkFactory{
	"file":  fileSinkFactory,
	"db":    dbSinkFactory,
	"kafka": kafkaSinkFactory,
}

// BuildSink dynamically constructs the sink based on JSON spec
func BuildSink(sinkSpec map[string]interface{}, in <-chan model.Event) error {
	sinkType, ok := sinkSpec["type"].(string)
	if !ok {
		return fmt.Errorf("sink missing 'type'")
	}
	factory, ok := registry[sinkType]
	if !ok {
		return fmt.Errorf("unknown sink type: %s", sinkType)
	}
	return factory(sinkSpec, in)
}

// -------- Adapters for each sink type --------

// File sink expects: { "type": "file", "path": "output.csv" }
func fileSinkFactory(params map[string]interface{}, in <-chan model.Event) error {
	path, ok := params["path"].(string)
	if !ok {
		return fmt.Errorf("file sink expects 'path'")
	}
	return FileSink(path, in)
}

// DB sink expects: { "type": "db", "dsn": "...", "table": "..." }
func dbSinkFactory(params map[string]interface{}, in <-chan model.Event) error {
	dsn, ok1 := params["dsn"].(string)
	table, ok2 := params["table"].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("db sink expects 'dsn' and 'table'")
	}
	return DBSink(DBSinkConfig{DSN: dsn, Table: table}, in)
}

// Kafka sink expects: { "type": "kafka", "brokers": [...], "topic": "..." }
func kafkaSinkFactory(params map[string]interface{}, in <-chan model.Event) error {
	brokersIface, ok1 := params["brokers"].([]interface{})
	topic, ok2 := params["topic"].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("kafka sink expects 'brokers' and 'topic'")
	}
	brokers := make([]string, len(brokersIface))
	for i, v := range brokersIface {
		str, ok := v.(string)
		if !ok {
			return fmt.Errorf("kafka sink: all brokers must be strings")
		}
		brokers[i] = str
	}
	return KafkaSink(KafkaSinkConfig{
		Brokers: brokers,
		Topic:   topic,
	}, in)
}
