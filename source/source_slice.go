package source

import (
	"errors"
	"fmt"
	"reflect"
	"rstreams/util"
	"time"
)

type sliceSource struct {
	in           interface{}
	out          chan Element
	onNext       chan bool
	done         chan bool
	error        chan error
	consumers    []chan<- bool
	drainTimeout time.Duration
	Verbose      bool
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

func (s *sliceSource) VerboseON() {
	s.Verbose = true
}

func (s *sliceSource) VerboseOFF() {
	s.Verbose = false
}

// Slice Source accepts any slice type. Panics on any other type
func Slice(i interface{}) *sliceSource {
	iVal := reflect.ValueOf(i)
	if err := util.IsSlice(iVal); err != nil {
		panic(err.Error())
	}
	return &sliceSource{
		in:           i,
		out:          make(chan Element),
		onNext:       make(chan bool),
		done:         make(chan bool),
		error:        make(chan error),
		drainTimeout: 30 * time.Millisecond,
		Verbose:      false,
	}
}

func (s *sliceSource) GetOutput() <-chan Element {
	return s.out
}

func (s *sliceSource) Emit() {
	iVal := reflect.ValueOf(s.in)
	run := true
	defer close(s.out)
	for run == true {
		select {
		case <-s.done:
			s.notifyConsumers()
			run = false
		case <-s.onNext:
			s.sendToErr("source => got demand signal")
			if iVal.Len() > 0 {
				s.out <- Element{
					Payload:   iVal.Index(0).Interface(),
					Timestamp: time.Now(),
				}
				iVal = iVal.Slice(1, iVal.Len())
			}
			if iVal.Len() == 0 {
				s.sendToErr("source => no more elements")
				s.notifyConsumers()
				run = false
			}
		}
	}
}

func (s *sliceSource) sendToErr(e string) {
	if s.Verbose {
		s.error <- errors.New(fmt.Sprintf(e))
	}
}

func (s *sliceSource) notifyConsumers() {
	if len(s.consumers) > 0 {
		// some time for consumer to get last element
		time.Sleep(s.drainTimeout)
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
