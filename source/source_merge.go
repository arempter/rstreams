package source

type merge struct {
	out     chan interface{}
	sources []Source
}

func (m merge) GetOutput() <-chan interface{} {
	return m.out
}

func (m merge) GetErrorCh() <-chan error {
	panic("implement me")
}

func (m merge) Emit() {
	//todo: add sync.WaitGroup
	output := func(in <-chan interface{}) {
		for e := range in {
			m.out <- e
		}
	}

	for _, s := range m.sources {
		go s.Emit()
		go output(s.GetOutput())
	}
}

func (merge) Stop() {
	panic("todo")
}

func MergeSources(sources ...Source) *merge {
	return &merge{
		out:     make(chan interface{}),
		sources: sources,
	}
}
