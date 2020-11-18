package processor

type MapF func(interface{}) interface{}
type MapFunc func(<-chan interface{}, MapF) <-chan interface{}
type MapFuncSpec struct {
	Body MapFunc
	ArgF MapF
}

type ProcFunc func(<-chan interface{}) <-chan interface{}
type ProcFuncSpec struct {
	Body ProcFunc
	ArgF interface{}
}

type Cond func(interface{}) bool
type FilterFunc func(<-chan interface{}, Cond) <-chan interface{}
type FilterFuncSpec struct {
	Body      FilterFunc
	Predicate Cond
}
