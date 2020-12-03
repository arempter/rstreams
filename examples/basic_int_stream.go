package main

import (
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
)

func main() {

	ints := []string{"1", "2", "three", "4", "5", "six", "7", "8", "9", "10", "eleven", "12", "13", "11", "14", "20", "21", "22", "23"}

	addOne := func(e int) int { return e + 1 }

	even := func(e int) bool {
		if e%2 == 0 {
			return true
		}
		return false
	}

	stream := stream.FromSource(source.Slice(ints))

	stream.
		Map(processor.ToInt).
		Map(addOne).
		Filter(even).
		To(sink.Foreach()).
		Run()

}
