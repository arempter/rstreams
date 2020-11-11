package source

type merge struct {
	out     chan interface{}
	sources []Source
}

func (m merge) GetOutput() <-chan interface{} {
	return m.out
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

func FromMergeSource(sources ...Source) *merge {
	return &merge{
		out:     make(chan interface{}),
		sources: sources,
	}
}