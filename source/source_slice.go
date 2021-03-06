package source

import (
	"errors"
	"fmt"
)

type sliceSource struct {
	in        []string
	out       chan interface{}
	onNext    chan bool
	done      chan bool
	error     chan error
	consumers []chan<- bool
}

func (s *sliceSource) OnNextCh() chan bool {
	return s.onNext
}

func (s *sliceSource) ErrorCh() <-chan error {
	return s.error
}

func (s sliceSource) Stop() {
	go func() {
		defer close(s.error)
		s.done <- true
	}()

}

func Slice(in []string) *sliceSource {
	return &sliceSource{
		in:     in,
		out:    make(chan interface{}, 5),
		onNext: make(chan bool),
		done:   make(chan bool),
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
		case <-s.done:
			s.notifyConsumers()
			run = false
		case <-s.onNext:
			s.sendToErr("source => got demand signal")
			if len(s.in) > 0 {
				s.out <- s.in[0]
				s.in = s.in[1:]
			}
			if len(s.in) == 0 {
				s.sendToErr("source => no more elements")
				s.notifyConsumers()
				run = false
			}
		}
	}
}

func (s *sliceSource) sendToErr(e string) {
	go func() {
		s.error <- errors.New(fmt.Sprintf(e))
	}()
}

func (s *sliceSource) notifyConsumers() {
	if len(s.consumers) > 0 {
		for _, done := range s.consumers {
			go func() {
				done <- true
			}()
		}
	}
}

func (s *sliceSource) Subscribe(consCh chan<- bool) {
	s.consumers = append(s.consumers, consCh)
}
