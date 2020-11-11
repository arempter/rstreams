package sink

type SinkFunc func(<-chan interface{})
