package source

type sliceSource struct {
	in   []string
	out  chan interface{}
	done chan bool
}

func (s sliceSource) OnNextCh() chan bool {
	panic("not supported")
}

func (s sliceSource) ErrorCh() <-chan error {
	panic("not supported")
}

func (s sliceSource) Stop() {
	panic("not supported")
}

func Slice(in []string) *sliceSource {
	return &sliceSource{
		in:  in,
		out: make(chan interface{}, 5),
	}
}

func (s sliceSource) GetOutput() <-chan interface{} {
	return s.out
}

func (s *sliceSource) Emit() {
	go func() {
		for _, e := range s.in {
			s.out <- e
		}
		close(s.out)
	}()
}
