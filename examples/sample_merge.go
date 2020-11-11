package main

import (
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"strings"
	"time"
)

func main() {
	words := strings.Split("merge starts one more goroutine to close the outbound channel after all sends on that channel are done super", " ")
	words1 := strings.Split("The merge function converts a list of channels to a single channel by starting a goroutine for each inbound", " ")
	words2 := strings.Split("This pattern allows each receiving stage to be written as a range loop and ensures that all goroutines exit", " ")
	wordsSource := source.FromSlice(words)
	wordsSource1 := source.FromSlice(words1)
	wordsSource2 := source.FromSlice(words2)

	stream := stream.NewStream(source.FromMergeSource(wordsSource, wordsSource1, wordsSource2))

	containsStringFunc := func(s string) bool {
		if !strings.Contains(strings.ToLower(s), "s") {
			return false
		}
		return true
	}

	stream.
		Filter(processor.Filter, containsStringFunc).
		Via(processor.ToUpper).
		To(sink.ForeachSink).
		Run()

	time.Sleep(1 * time.Second)
}
