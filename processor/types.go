package processor

type StepFuncPredicate func(<-chan interface{}, interface{}) <-chan interface{}
type StepFuncWithPredicate struct {
	Body      StepFuncPredicate
	Predicate interface{}
}

type StepFunc func(<-chan interface{}) <-chan interface{}
type StepFuncSpec struct {
	Body StepFunc
	ArgF interface{}
}
