package stream

import (
	"log"
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
)

//todo:
// logging && error (err type verb, debug etc)
// parallel
// source freq
// restart source on err
// align buffer size for all components
// rework wireTap

type Stream interface {
	Filter(c interface{}) *stream
	Map(p interface{}) *stream
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

func (s *stream) Map(predicate interface{}) *stream {
	s.steps = append(s.steps, processor.StepFuncWithPredicate{
		Body:      processor.Map,
		Predicate: predicate,
	})
	return s
}

func (s *stream) Filter(predicate interface{}) *stream {
	s.steps = append(s.steps, processor.StepFuncWithPredicate{
		Body:      processor.Filter,
		Predicate: predicate,
	})
	return s
}

func (s *stream) To(f sink.Collector) *stream {
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
		case processor.StepFuncWithPredicate:
			fSpec := step.(processor.StepFuncWithPredicate)
			mOut := fSpec.Body(s.inChan, fSpec.Predicate)
			s.inChan = mOut
		case processor.StepFuncSpec:
			pf := step.(processor.StepFuncSpec).Body
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
