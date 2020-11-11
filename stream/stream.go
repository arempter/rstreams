package stream

import (
	"fmt"
	"rstreams/processors"
	"rstreams/sinks"
	"rstreams/source"
)

//todo: logging && error
type Stream interface {
	Via(f processors.ProcFunc) *stream
	Filter(f processors.FilterFunc, c func(string) bool) *stream
	To(f sinks.SinkFunc) *stream
	Run()
	Stop()
}

type stream struct {
	source source.Source
	steps  []interface{}
	inChan <-chan interface{}
}

func (s *stream) Filter(f processors.FilterFunc, predicate func(string) bool) *stream {
	s.steps = append(s.steps, processors.FilterFuncSpec{
		Body:      f,
		Predicate: predicate,
	})
	return s
}

func (s *stream) To(f sinks.SinkFunc) *stream {
	s.steps = append(s.steps, f)
	return s
}

func (s *stream) Via(f processors.ProcFunc) *stream {
	s.steps = append(s.steps, f)
	return s
}

func (s *stream) Run() {
	go s.buildDAG()
}

func (s *stream) buildDAG() {
	go s.source.Emit()
	s.inChan = s.source.GetOutput()
	for _, step := range s.steps {
		switch step.(type) {
		case processors.FilterFuncSpec:
			ff := step.(processors.FilterFuncSpec)
			fOut := ff.Body(s.inChan, ff.Predicate)
			s.inChan = fOut
		case processors.ProcFunc:
			pOut := step.(processors.ProcFunc)(s.inChan)
			s.inChan = pOut
		case sinks.SinkFunc:
			step.(sinks.SinkFunc)(s.inChan)
		default:
			panic("Unsupported step type")
		}
	}
}
func (s *stream) Stop() {
	fmt.Println("Sending stop source signal...")
	s.source.Stop()
}

func NewStream(source source.Source) *stream {
	return &stream{
		source: source,
	}
}
