package source

type sliceSource struct {
	in  []string
	out <-chan string
}

func FromSlice(in []string) *sliceSource {
	return &sliceSource{in: in}
}

func (s sliceSource) GetOutput() <-chan string {
	return s.out
}

func (s *sliceSource) Emit() *sliceSource {
	ch := make(chan string, 5)
	go func() {
		for _, e := range s.in {
			ch <- e
		}
		close(ch)
	}()
	s.out = ch
	return s
}
