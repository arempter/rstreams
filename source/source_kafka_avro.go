package source

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
)

type kafkaAvroSource struct {
	kafkaConf *kafka.ConfigMap
	topics    []string
	schema    string
	out       chan interface{}
	done      chan bool
}

func (k kafkaAvroSource) GetOutput() <-chan interface{} {
	return k.out
}
func (k kafkaAvroSource) Stop() {
	log.Println("sending stop signal to kafka consumer...")
	k.done <- true
}

func (k kafkaAvroSource) GetErrorCh() <-chan string {
	panic("implement me")
}

func (k *kafkaAvroSource) Emit() {
	c, err := kafka.NewConsumer(k.kafkaConf)
	if err != nil {
		log.Fatal("failed to create kafka consumer")
		os.Exit(1)
	}
	defer c.Close()

	err = c.SubscribeTopics(k.topics, nil)

	run := true
	for run == true {
		select {
		case <-k.done:
			log.Println("stopping kafka source")
			run = false
		default:
			ev := c.Poll(3000)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				if err != nil {
					log.Println("failed to deserialize from avro", err)
				}
				go func() {
					k.out <- e.Value
				}()
			case kafka.Error:
				if e.Code() == kafka.ErrAllBrokersDown {
					log.Println("failed", e.Error())
					run = false
				}
			default:
				log.Println("not sure what to do, not kafka.Message")
			}
		}
	}
}

func KafkaAvro(kafkaConf *kafka.ConfigMap, topics []string, schema string) *kafkaAvroSource {
	return &kafkaAvroSource{
		kafkaConf: kafkaConf,
		topics:    topics,
		schema:    schema,
		out:       make(chan interface{}, 5),
		done:      make(chan bool, 1),
	}
}
