package main

import (
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"time"
)

func main() {
	sampleUrls := []string{"http://localhost:1234/device/1", "http://localhost:1234/device/2", "http://localhost:1234/device/3", "http://localhost:1234/device/4"}

	stream := stream.FromSource(source.Http(sampleUrls, time.Second/10))

	//go func() {
	//	time.Sleep(5 * time.Second)
	//	stream.Stop()
	//	os.Exit(0)
	//}()

	onlyValidJson := func(s interface{}) bool {
		if len(s.([]byte)) > 5 {
			return true
		}
		return false
	}

	go stream.WireTap()

	stream.
		Filter(onlyValidJson).
		To(sink.Foreach()).
		Run()
}
