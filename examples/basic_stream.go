package main

import (
	"fmt"
	_ "net/http/pprof"
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"runtime"
	"strings"
	"time"
)

func main() {
	words := strings.Split("Now you've got everything running, you may wonder: what now? This section we'll describe a basic flow on how to use Rokku to perform operations in S3. You may refer to the What is Rokku? document before diving in here. That will introduce you to the various components used.", " ")
	stream := stream.FromSource(source.Slice(words))

	containsStringFunc := func(i string) bool {
		if !strings.Contains(strings.ToLower(i), "a") {
			return false
		}
		return true
	}

	//go stream.WireTap()

	stream.
		Filter(containsStringFunc).
		Map(processor.ToUpper).
		CountBy(time.Second / 4).
		To(sink.Foreach()).
		Run()

	fmt.Println("took", time.Since(start))
}
