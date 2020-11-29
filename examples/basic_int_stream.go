package main

import (
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
)

func main() {

	ints := []int{1, 2, 3, 4, 5, 7, 9, 15, 20}

	addOne := func(e int) int { return e + 1 }

	even := func(e int) bool {
		if e%2 == 0 {
			return true
		}
		return false
	}

	stream.FromSource(source.Slice(ints)).
		Map(addOne).
		Filter(even).
		To(sink.Foreach()).
		Run()
}
