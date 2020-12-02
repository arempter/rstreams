package sink

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"rstreams/source"
	"time"
)

type foreachBP struct {
	bufSize int
	onNext  chan bool
	error   chan error
	done    chan bool
}

func (f *foreachBP) ErrorCh() <-chan error {
	return f.error
}

func (f *foreachBP) DoneCh() chan<- bool {
	return f.done
}

func BufferedForeach() *foreachBP {
	return &foreachBP{
		bufSize: 256,
		error:   make(chan error),
		done:    make(chan bool),
	}
}

func (f *foreachBP) SetOnNextCh(c chan bool) {
	f.onNext = c
}

func (f foreachBP) Receive(in <-chan source.Element) {
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
			f.error <- errors.New(fmt.Sprintf("sink => buffer not empty, reading..."))
			time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
			buffer.WriteTo(os.Stdout)
		}
	}

	f.onNext <- true
	for e := range in {
		switch e.Payload.(type) {
		case []byte:
			bs := e.Payload.([]byte)
			go consume(bs)
			process()
		}
	}
}
