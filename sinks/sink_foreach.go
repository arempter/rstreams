package sinks

import (
	"fmt"
)

func ForeachSink(in <-chan string) {
	for m := range in {
		fmt.Println(m)
	}
}
