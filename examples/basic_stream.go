package main

import (
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"strings"
)

func main() {
	words := strings.Split("Now you've got everything running, you may wonder: what now? This section we'll describe a basic flow on how to use Rokku to perform operations in S3. You may refer to the What is Rokku? document before diving in here. That will introduce you to the various components used.", " ")
	stream := stream.FromSource(source.Slice(words))

	containsStringFunc := func(i interface{}) bool {
		if !strings.Contains(strings.ToLower(i.(string)), "a") {
			return false
		}
		return true
	}

	addOK := func(i interface{}) interface{} {
		return i.(string) + "-ok"
	}

	//go stream.WireTap()

	stream.
		Filter(processor.Filter, containsStringFunc).
		Map(processor.Map, addOK).
		Via(processor.ToUpper()).
		To(sink.Foreach()).
		Run()

	stream.Stop()

}
