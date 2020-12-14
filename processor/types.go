package processor

import (
	"rstreams/source"
)

type StepFuncPredicate func(in <-chan source.Element, f interface{}, par int) chan source.Element
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
