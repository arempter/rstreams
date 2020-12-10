package processor

import (
	"rstreams/source"
)

type StepFuncPredicate func(in <-chan source.Element, f interface{}, out chan source.Element, par int)
type StepFuncWithPredicate struct {
	Body      StepFuncPredicate
	Predicate interface{}
}

type StepFuncParWithPredicate struct {
	Body      StepFuncPredicate
	Predicate interface{}
	Parallel  int
}

type StepFunc func(<-chan source.Element, interface{}) <-chan source.Element
type StepFuncSpec struct {
	Body StepFunc
	ArgF interface{}
}
