package stream

import (
	"log"
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"time"
)

//todo:
// logging && error (err type verb, debug etc)
// parallel - grouped
// restart source on err
// ToMat

type Stream interface {
	Filter(c interface{}) *stream
	Map(p interface{}) *stream
	To(f sink.Collector) *stream
	CountBy(windowSize time.Duration) *stream
	Run()
	Stop()
	WireTap()
}

type stream struct {
	source source.Source
	steps  []interface{}
	inChan <-chan source.Element
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

func (s *stream) MapPar(par int, predicate interface{}) *stream {
	s.steps = append(s.steps, processor.StepFuncParWithPredicate{
		Body:      processor.MapPar,
		Predicate: predicate,
		Parallel:  par,
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

func (s *stream) CountBy(windowSize time.Duration) *stream {
	s.steps = append(s.steps, processor.StepFuncSpec{
		Body: processor.Counter,
		ArgF: windowSize,
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
			out := make(chan source.Element)
			fSpec := step.(processor.StepFuncWithPredicate)
			fSpec.Body(s.inChan, fSpec.Predicate, out, 0)
			s.inChan = out
		case processor.StepFuncParWithPredicate:
			out := make(chan source.Element)
			fSpec := step.(processor.StepFuncParWithPredicate)
			fSpec.Body(s.inChan, fSpec.Predicate, out, fSpec.Parallel)
			s.inChan = out
		case processor.StepFuncSpec:
			sf := step.(processor.StepFuncSpec)
			pOut := sf.Body(s.inChan, sf.ArgF)
			s.inChan = pOut
		case sink.Collector:
			c := step.(sink.Collector)
			s.source.Subscribe(c.DoneCh())
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
	s.source.VerboseON()
	readErr := func(errIn <-chan error) {
		for e := range errIn {
			log.Println("DEBUG", e)
		}
	}
	for _, c := range s.errors {
		go readErr(c)
	}
}
