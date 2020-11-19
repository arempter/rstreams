package stream

import (
	"log"
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
)

//todo: logging && error
//todo: parallel

type Stream interface {
	Via(f processor.ProcFuncSpec) *stream
	Filter(f processor.FilterFunc, c func(string) bool) *stream //todo: move to procFunc?
	Map(f processor.MapFunc, p processor.MapF) *stream
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
	errors []<-chan error
}

func FromSource(source source.Source) *stream {
	return &stream{
		source: source,
		done:   make(chan bool),
	}
}

func (s *stream) Map(f processor.MapFunc, predicate processor.MapF) *stream {
	s.steps = append(s.steps, processor.MapFuncSpec{
		Body: f,
		ArgF: predicate,
	})
	return s
}

func (s *stream) Filter(f processor.FilterFunc, predicate processor.Cond) *stream {
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

func (s *stream) Via(f processor.ProcFuncSpec) *stream {
	s.steps = append(s.steps, f)
	return s
}

func (s *stream) Run() {
	s.runnableDAG()
}

func (s *stream) runnableDAG() {
	go s.source.Emit()
	s.inChan = s.source.GetOutput()
	s.errors = append(s.errors, s.source.ErrorCh())
	for _, step := range s.steps {
		switch step.(type) {
		case processor.MapFuncSpec:
			mf := step.(processor.MapFuncSpec)
			mOut := mf.Body(s.inChan, mf.ArgF)
			s.inChan = mOut
		case processor.FilterFuncSpec:
			ff := step.(processor.FilterFuncSpec)
			fOut := ff.Body(s.inChan, ff.Predicate)
			s.inChan = fOut
		case processor.ProcFuncSpec:
			pf := step.(processor.ProcFuncSpec).Body
			pOut := pf(s.inChan)
			s.inChan = pOut
		case sink.Collector:
			c := step.(sink.Collector)
			s.source.Subscribe(c.DoneCh())
			c.SetOnNextCh(s.source.OnNextCh())
			s.errors = append(s.errors, c.ErrorCh())
			c.Receive(s.inChan)
		default:
			panic("Unsupported step type")
		}
	}
}

func (s *stream) Stop() {
	s.source.Stop()
}

// todo: change to multiplex and accept f to process E
func (s *stream) WireTap() {
	readErr := func(errIn <-chan error) {
		for e := range errIn {
			log.Println("DEBUG", e)
		}
	}
	for _, c := range s.errors {
		go readErr(c)
	}
}
