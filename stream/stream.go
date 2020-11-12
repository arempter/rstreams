package stream

import (
	"fmt"
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
)

//todo: logging && error
//todo: parallel
type Stream interface {
	Via(f processor.ProcFunc) *stream
	Filter(f processor.FilterFunc, c func(string) bool) *stream //todo: move to procFunc?
	To(f sink.SinkFunc) *stream
	Run()
	Stop()
}

type stream struct {
	source source.Source
	steps  []interface{}
	inChan <-chan interface{}
	done   chan bool
}

func (s *stream) Filter(f processor.FilterFunc, predicate func(string) bool) *stream {
	s.steps = append(s.steps, processor.FilterFuncSpec{
		Body:      f,
		Predicate: predicate,
	})
	return s
}

func (s *stream) To(f sink.SinkFunc) *stream {
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
		case sink.SinkFunc:
			step.(sink.SinkFunc)(s.inChan)
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

func NewStream(source source.Source) *stream {
	return &stream{
		source: source,
		done:   make(chan bool),
	}
}
