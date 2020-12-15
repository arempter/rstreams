package processor

import (
	"rstreams/source"
)

type Step interface {
	Apply()
}

type StepFunc func(<-chan source.Element, interface{}) chan source.Element
type StepFuncSpec struct {
	Body StepFunc
	ArgF interface{}
}

func (s StepFuncSpec) Apply(in <-chan source.Element) chan source.Element {
	return s.Body(in, s.ArgF)
}
