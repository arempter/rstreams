package main

import (
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
)

func main() {
	sampleUrls := []string{"http://localhost:1234/device/1", "http://localhost:1234/device/2", "http://localhost:1234/device/3", "http://localhost:1234/device/4"}

	stream := stream.FromSource(source.Http(sampleUrls))

	//go func() {
	//	time.Sleep(3 * time.Second)
	//	stream.Stop()
	//	os.Exit(0)
	//}()

	go stream.WireTap()

	stream.
		To(sink.BufferedForeach()).
		Run()
}
