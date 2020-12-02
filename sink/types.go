package sink

import (
	"rstreams/source"
)

type Collector interface {
	SetOnNextCh(chan bool)
	Receive(in <-chan source.Element)
	ErrorCh() <-chan error
	DoneCh() chan<- bool
}
