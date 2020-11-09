package source

type SliceSource interface {
	GetOutput() <-chan string
	Emit() *sliceSource
}

type KafkaAvroSource interface {
	GetOutput() <-chan string
	Emit() *kafkaAvroSource
}
