package sinks

type SinkFunc func(<-chan interface{})
