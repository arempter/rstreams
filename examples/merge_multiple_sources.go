package main

import (
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"strings"
)

func main() {
	words := strings.Split("merge starts one more goroutine to close the outbound channel after all sends on that channel are done super", " ")
	words1 := strings.Split("The merge function converts a list of channels to a single channel by starting a goroutine for each inbound", " ")
	words2 := strings.Split("This pattern allows each receiving stage to be written as a range loop and ensures that all goroutines exit", " ")
	wordsSource := source.Slice(words)
	wordsSource1 := source.Slice(words1)
	wordsSource2 := source.Slice(words2)

	stream := stream.FromSource(source.MergeSources(wordsSource, wordsSource1, wordsSource2))

	containsStringFunc := func(i interface{}) bool {
		switch i.(type) {
		case string:
			if !strings.Contains(strings.ToLower(i.(string)), "s") {
				return false
			}
			return true
		default:
			return false
		}
	}

	stream.
		Filter(containsStringFunc).
		Via(processor.ToUpper()).
		To(sink.Foreach()).
		Run()
}
