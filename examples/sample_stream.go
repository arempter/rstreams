package main

import (
	"rstreams/processors"
	"rstreams/sinks"
	"rstreams/source"
	"rstreams/stream"
	"strings"
	"time"
)

func main() {
	words := strings.Split("Now you've got everything running, you may wonder: what now? This section we'll describe a basic flow on how to use Rokku to perform operations in S3. You may refer to the What is Rokku? document before diving in here. That will introduce you to the various components used.", " ")
	wordsSource := source.FromSlice(words)
	stream := stream.NewStream(wordsSource)

	containsStringFunc := func(s string) bool {
		if !strings.Contains(strings.ToLower(s), "a") {
			return false
		}
		return true
	}

	stream.
		Filter(processors.Filter, containsStringFunc).
		Via(processors.ToUpper).
		To(sinks.ForeachSink).
		Run()

	time.Sleep(1 * time.Second)
}
