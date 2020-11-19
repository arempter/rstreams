package sink

type Collector interface {
	SetOnNextCh(chan bool)
	Receive(in <-chan interface{})
	ErrorCh() <-chan error
	DoneCh() chan<- bool
}
