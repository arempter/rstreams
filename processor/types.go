package processor

import (
	"rstreams/source"
)

type StepFuncPredicate func(in <-chan source.Element, f interface{}, out chan source.Element)
type StepFuncWithPredicate struct {
	Body      StepFuncPredicate
	Predicate interface{}
}

type StepFunc func(<-chan source.Element) <-chan source.Element
type StepFuncSpec struct {
	Body StepFunc
	ArgF interface{}
}
