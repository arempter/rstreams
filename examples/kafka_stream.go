package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hamba/avro"
	"log"
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"strings"
	"time"
)

const schema = `{
			"namespace": "devices.Generator",
			"type": "record",
			"name": "GeneratorMessage",
			"fields": [
				{"name": "id", "type": "string"},
				{"name": "state", "type": "string"},
				{"name": "value", "type": "string"}
				]
			}`

func main() {
	kc := &kafka.ConfigMap{
		"bootstrap.servers":     "localhost:9092",
		"broker.address.family": "v4",
		"group.id":              "rstreams_test",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	}
	hasValueFunc := func(s string) bool {
		if !strings.Contains(strings.ToLower(s), "1.") {
			return false
		}
		return true
	}

	stream := stream.NewStream(source.FromKafkaAvro(kc, []string{"reactiveLab"}, schema))
	stream.
		Filter(processor.Filter, hasValueFunc).
		Via(decodeToNative).
		To(sink.ForeachSink).
		Run()

	// test stream stop
	time.Sleep(5 * time.Second)
	stream.Stop()
}

type DeviceData struct {
	Id    string `avro:"id"`
	State string `avro:"state"`
	Value string `avro:"value"`
}

func decodeToNative(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	deviceData := DeviceData{}
	schema, err := avro.Parse(schema)
	if err != nil {
		log.Fatal("failed to parse schema")
	}
	go func() {
		for d := range in {
			switch d.(type) {
			case []byte:
				err = avro.Unmarshal(schema, d.([]byte)[5:], &deviceData)
				if err != nil {
					log.Println("failed to unmarshall data")
				}
				out <- deviceData
			default:
			}
		}
	}()
	return out
}
