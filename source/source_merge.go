package source

type merge struct {
	out       chan interface{}
	sources   []Source
	onNext    chan bool
	done      chan bool
	error     chan error
	consumers []chan<- bool
}

func (m *merge) OnNextCh() chan bool {
	return m.onNext
}

func (m *merge) ErrorCh() <-chan error {
	return m.error
}

func (m *merge) GetOutput() <-chan interface{} {
	return m.out
}

func (m *merge) Emit() {
	processSource := func(in <-chan interface{}, onNext chan bool) {
		run := true
		for run == true {
			select {
			case <-m.done:
				run = false
			case onNext <- true:
			case e, open := <-in:
				if !open {
					run = false
				}
				m.out <- e
			}
		}
	}

	defer close(m.out)
	for _, s := range m.sources {
		go s.Emit()
		processSource(s.GetOutput(), s.OnNextCh())
	}
}

func (merge) Stop() {
	panic("todo")
}

func (m *merge) Subscribe(consCh chan<- bool) {
	m.consumers = append(m.consumers, consCh)
}

func MergeSources(sources ...Source) *merge {
	return &merge{
		out:     make(chan interface{}),
		sources: sources,
		onNext:  make(chan bool),
		done:    make(chan bool),
		error:   make(chan error),
	}
}
