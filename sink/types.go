package sink

type Collector interface {
	SetOnNextCh(chan bool)
	HasBackpressure() bool
	Receive(in <-chan interface{})
}
