package main

import (
	"fmt"
	"rstreams/source"
	"rstreams/stream"
	"time"
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

	slowConsumer := func(in <-chan interface{}) {
		var noOfElements = 0
		for e := range in {
			noOfElements += 1
			switch e.(type) {
			case []byte:
				time.Sleep(1200 * time.Millisecond)
				fmt.Print(string(e.([]byte)))
			}
		}
	}

	stream.
		To(slowConsumer).
		Run()
}
