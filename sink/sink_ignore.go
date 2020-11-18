package sink

type ignore struct {
	onNext chan bool
	error  chan error
}

func (i *ignore) SetOnNextCh(c chan bool) {
	i.onNext = c
}

func (i ignore) Receive(in <-chan interface{}) {
	i.onNext <- true
	for range in {
		go func() { i.onNext <- true }()
	}
}

func (i *ignore) ErrorCh() <-chan error {
	return i.error
}

func Ignore() *ignore {
	return &ignore{
		error: make(chan error),
	}
}
