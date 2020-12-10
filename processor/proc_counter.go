package processor

import (
	"rstreams/source"
	"time"
)

func Counter(in <-chan source.Element) <-chan source.Element {
	c := 0
	out := make(chan source.Element)
	go func() {
		defer close(out)
		winStart := time.Now().Add(10 * time.Second)
		for e := range in {
			if e.Payload != nil {
				out <- e
				if !e.Timestamp.After(winStart) {
					c += 1
				} else {
					winStart = time.Now().Add(10 * time.Second)
					c = 0
				}
			}
		}
		defer println("elements count ===>", c)
	}()
	//close(out)
	return out
}
