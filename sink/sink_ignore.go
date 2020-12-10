package sink

import (
	"rstreams/source"
)

type ignore struct {
	onNext chan bool
	error  chan error
	done   chan bool
}

func (i *ignore) DoneCh() chan<- bool {
	return i.done
}

func (i *ignore) SetOnNextCh(c chan bool) {
	i.onNext = c
}

func (i *ignore) Receive(in <-chan source.Element) {
	for range in {
	}
}

func (i *ignore) ErrorCh() <-chan error {
	return i.error
}

func Ignore() *ignore {
	return &ignore{
		error: make(chan error),
		done:  make(chan bool),
	}
}
