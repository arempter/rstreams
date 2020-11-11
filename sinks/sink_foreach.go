package sinks

import (
	"fmt"
)

func ForeachSink(in <-chan interface{}) {
	for m := range in {
		fmt.Println(m)
	}
}
