package source

type sliceSource struct {
	in   []string
	out  chan string
	done chan bool
}

func (s sliceSource) Stop() {
	panic("not supported")
}

func FromSlice(in []string) *sliceSource {
	return &sliceSource{
		in:  in,
		out: make(chan string, 5),
	}
}

func (s sliceSource) GetOutput() <-chan string {
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
