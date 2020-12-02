package sink

import (
	"fmt"
	"rstreams/source"
)

type foreach struct {
	onNext chan bool
	error  chan error
	done   chan bool
}

func (f *foreach) ErrorCh() <-chan error {
	return f.error
}

func (f *foreach) DoneCh() chan<- bool {
	return f.done
}

func Foreach() *foreach {
	return &foreach{
		error: make(chan error),
		done:  make(chan bool),
	}
}

func (f *foreach) SetOnNextCh(c chan bool) {
	f.onNext = c
}

func (f *foreach) Receive(in <-chan source.Element) {
	run := true
	for run == true {
		select {
		case <-f.done:
			run = false
		case f.onNext <- true:
		case e, open := <-in:
			if !open {
				run = false
			}
			switch e.Payload.(type) {
			case string:
				fmt.Println(e.Payload)
			case int:
				fmt.Println(e.Payload)
			case []byte:
				fmt.Print(string(e.Payload.([]byte)))
			}
		}
	}
}
