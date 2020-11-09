package sinks

type SinkFunc func(<-chan string)
