package processor

import (
	"fmt"
	"rstreams/source"
	"time"
)

func Counter(in <-chan source.Element, windowSize interface{}) <-chan source.Element {
	c := 0
	w, ok := windowSize.(time.Duration)
	if !ok {
		panic("window size must be of type time.Duration")
	}
	out := make(chan source.Element)
	go func() {
		defer close(out)
		go func() {
			ticker := time.NewTicker(w)
			defer ticker.Stop()
			for range ticker.C {
				println(fmt.Sprintf("current elements count per %s window ===> %d", w, c))
			}
		}()
		winStart := time.Now().Add(w)
		for e := range in {
			if e.Payload != nil {
				out <- e
				if !e.Timestamp.After(winStart) {
					c += 1
				} else {
					winStart = time.Now().Add(w)
					c = 0
				}
			}
		}
	}()
	return out
}
