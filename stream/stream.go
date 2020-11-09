package stream

import (
	"rstreams/processors"
	"rstreams/sinks"
	"rstreams/source"
)

//todo: drain / stop

type Stream interface {
	Via(f processors.ProcFunc) *stream
	Filter(f processors.FilterFunc, c func(string) bool) *stream
	To(f sinks.SinkFunc) *stream
	Run() *stream
}

type stream struct {
	source source.SliceSource
	steps  []interface{}
	inChan <-chan string
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

func (s *stream) Run() *stream {
	return s.buildDAG()
}

func (s *stream) buildDAG() *stream {
	s.inChan = s.source.Emit().GetOutput()
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
	return s
}

func NewStream(source source.SliceSource) *stream {
	return &stream{
		source: source,
	}
}
