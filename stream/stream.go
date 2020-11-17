package stream

import (
	"fmt"
	"log"
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
)

//todo: logging && error
//todo: parallel

type Stream interface {
	Via(f processor.ProcFunc) *stream
	Filter(f processor.FilterFunc, c func(string) bool) *stream //todo: move to procFunc?
	To(f sink.Collector) *stream
	Run()
	Stop()
	WireTap()
}

type stream struct {
	source source.Source
	steps  []interface{}
	inChan <-chan interface{}
	done   chan bool
	error  chan error
}

func FromSource(source source.Source) *stream {
	return &stream{
		source: source,
		done:   make(chan bool),
		error:  make(chan error),
	}
}
func (s *stream) Filter(f processor.FilterFunc, predicate func(interface{}) bool) *stream {
	s.steps = append(s.steps, processor.FilterFuncSpec{
		Body:      f,
		Predicate: predicate,
	})
	return s
}

func (s *stream) To(f sink.Collector) *stream {
	s.steps = append(s.steps, f)
	return s
}

func (s *stream) Via(f processor.ProcFunc) *stream {
	s.steps = append(s.steps, f)
	return s
}

func (s *stream) Run() {
	s.runnableDAG()
}

func (s *stream) runnableDAG() {
	go s.source.Emit()
	s.inChan = s.source.GetOutput()
	for _, step := range s.steps {
		switch step.(type) {
		case processor.FilterFuncSpec:
			ff := step.(processor.FilterFuncSpec)
			fOut := ff.Body(s.inChan, ff.Predicate)
			s.inChan = fOut
		case processor.ProcFunc:
			pOut := step.(processor.ProcFunc)(s.inChan)
			s.inChan = pOut
		case sink.Collector:
			c := step.(sink.Collector)
			if c.HasBackpressure() {
				fmt.Println("got BPS connecting ch")
				c.SetOnNextCh(s.source.OnNextCh())
			}
			c.Receive(s.inChan)
		default:
			panic("Unsupported step type")
		}
	}
}

// stop source for now...
func (s *stream) Stop() {
	fmt.Println("Sending stop source signal...")
	s.source.Stop()
}

func (s *stream) WireTap() {
	go func() {
		for se := range s.source.ErrorCh() {
			s.error <- se
		}
	}()
	for e := range s.error {
		log.Println("DEBUG", e)
	}
}
