package main

import (
	"bytes"
	"fmt"
	"os"
	"rstreams/source"
	"rstreams/stream"
)

func main() {
	sampleUrls := []string{"http://localhost:1234/device/1", "http://localhost:1234/device/2", "http://localhost:1234/device/3", "http://localhost:1234/device/4"}

	stream := stream.FromSource(source.Http(sampleUrls))

	go stream.WireTap()

	fastConsumer := func(in <-chan interface{}) {
		var msgSize = 71
		var buffer bytes.Buffer
		for e := range in {
			switch e.(type) {
			case []byte:
				if buffer.Len() < 10*msgSize {
					bs := e.([]byte)
					fmt.Sprintf("not enough elements to flush buffer, %d->%d", buffer.Len(), buffer.Cap())
					if len(bs) > 5 {
						buffer.Write(bs)
					}
				} else {
					fmt.Println("flushing buffer ")
					buffer.WriteTo(os.Stdout)
					buffer.Reset()
				}
			}
		}
	}

	stream.
		To(fastConsumer).
		Run()
}
