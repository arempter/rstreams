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
	done         chan bool
	error        chan error
	consumers    []chan<- bool
	drainTimeout time.Duration
	Verbose      bool
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
		out:          make(chan Element, 1),
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
	defer close(s.out)

	for e := 0; e < iVal.Len(); e++ {
		nextE := Element{
			Payload:   iVal.Index(e).Interface(),
			Timestamp: time.Now(),
		}
		select {
		case <-s.done:
			s.notifyConsumers()
			return
		case s.out <- nextE:
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
