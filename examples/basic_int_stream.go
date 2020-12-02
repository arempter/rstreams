package main

import (
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
)

func main() {

	ints := []string{"1", "2", "three", "4", "5", "six", "7"}

	addOne := func(e int) int { return e + 1 }

	even := func(e int) bool {
		if e%2 == 0 {
			return true
		}
		return false
	}

	stream.FromSource(source.Slice(ints)).
		Map(processor.ToInt).
		Map(addOne).
		Filter(even).
		To(sink.Foreach()).
		Run()
}
