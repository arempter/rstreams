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
	CountByWindow(windowSize time.Duration) *stream
	Run()
	Stop()
	WireTap()
}

type stream struct {
	source source.Source
	steps  []interface{}
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

func (s *stream) CountByWindow(windowSize time.Duration) *stream {
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
	var elements <-chan source.Element
	go s.source.Emit()
	elements = s.source.GetOutput()
	s.errors = append(s.errors, s.source.ErrorCh())
	for _, step := range s.steps {
		switch step.(type) {
		case processor.StepFuncWithPredicate:
			stepFunc := step.(processor.StepFuncWithPredicate)
			stepOut := stepFunc.Body(elements, stepFunc.Predicate, 0)
			elements = stepOut
		case processor.StepFuncParWithPredicate:
			stepFunc := step.(processor.StepFuncParWithPredicate)
			stepOut := stepFunc.Body(elements, stepFunc.Predicate, stepFunc.Parallel)
			elements = stepOut
		case processor.StepFuncSpec:
			stepFunc := step.(processor.StepFuncSpec)
			stepOut := stepFunc.Body(elements, stepFunc.ArgF)
			elements = stepOut
		case sink.Collector:
			stepFunc := step.(sink.Collector)
			s.source.Subscribe(stepFunc.DoneCh())
			s.errors = append(s.errors, stepFunc.ErrorCh())
			stepFunc.Receive(elements)
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
