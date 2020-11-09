package processors

type Predicate func(string) bool
type ProcFunc func(<-chan string) <-chan string
type FilterFunc func(<-chan string, Predicate) <-chan string
type FilterFuncSpec struct {
	Body      FilterFunc
	Predicate Predicate
}
