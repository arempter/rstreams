package processor

type Predicate func(interface{}) bool
type ProcFunc func(<-chan interface{}) <-chan interface{}
type FilterFunc func(<-chan interface{}, Predicate) <-chan interface{}
type FilterFuncSpec struct {
	Body      FilterFunc
	Predicate Predicate
}
