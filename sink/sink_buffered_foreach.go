package sink

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type foreachBP struct {
	bufSize int
	onNext  chan bool
}

func BufferedForeach() *foreachBP {
	return &foreachBP{
		bufSize: 256,
	}
}

func (f *foreachBP) SetOnNextCh(c chan bool) {
	f.onNext = c
}

func (f foreachBP) HasBackpressure() bool {
	return true
}

func (f *foreachBP) Receive(in <-chan interface{}) {
	var buffer bytes.Buffer

	consume := func(bs []byte) {
		buffer.Write(bs)
		if buffer.Len() < f.bufSize {
			f.onNext <- true
		}
	}

	// simulate slow processing
	process := func() {
		if buffer.Len() > int(float32(f.bufSize)*0.6) {
			fmt.Println("not empty, reading buffer")
			time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
			buffer.WriteTo(os.Stdout)
		}
	}

	f.onNext <- true
	for e := range in {
		switch e.(type) {
		case []byte:
			bs := e.([]byte)
			go consume(bs)
			process()
		}
	}
}
