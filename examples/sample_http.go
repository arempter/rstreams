package main

import (
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"time"
)

func main() {

	sampleUrls := []string{"http://localhost:1234/device/1", "http://localhost:1234/device/2", "http://localhost:1234/device/3", "http://localhost:1234/device/4"}

	stream := stream.NewStream(source.FromHttp(sampleUrls))

	go func() {
		time.Sleep(3 * time.Second)
		stream.Stop()
	}()

	stream.
		To(sink.Foreach).
		Run()
}
