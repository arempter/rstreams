package sink

import (
	"fmt"
)

func ForeachSink(in <-chan interface{}) {
	for e := range in {
		fmt.Println(e)
	}
}
