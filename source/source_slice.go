package source

import (
	"errors"
	"fmt"
)

type sliceSource struct {
	in     []string
	out    chan interface{}
	onNext chan bool
	done   chan bool
	error  chan error
}

func (s *sliceSource) OnNextCh() chan bool {
	return s.onNext
}

func (s *sliceSource) ErrorCh() <-chan error {
	return s.error
}

func (s sliceSource) Stop() {
	defer close(s.error)
}

func Slice(in []string) *sliceSource {
	return &sliceSource{
		in:     in,
		out:    make(chan interface{}, 5),
		onNext: make(chan bool),
		error:  make(chan error),
	}
}

func (s *sliceSource) GetOutput() <-chan interface{} {
	return s.out
}

func (s *sliceSource) Emit() {
	run := true
	defer close(s.out)
	for run == true {
		select {
		case <-s.onNext:
			go func() { s.error <- errors.New(fmt.Sprintf("source => got demand signal")) }()
			if len(s.in) > 0 {
				s.out <- s.in[0]
				s.in = s.in[1:]
			}
			if len(s.in) == 0 {
				go func() { s.error <- errors.New(fmt.Sprintf("source => no more elements")) }()
				run = false
			}
		}
	}
}
