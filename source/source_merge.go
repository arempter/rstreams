package source

type merge struct {
	out       chan Element
	sources   []Source
	done      chan bool
	error     chan error
	consumers []chan<- bool
	Verbose   bool
}

func (m *merge) ErrorCh() <-chan error {
	return m.error
}

func (m *merge) GetOutput() <-chan Element {
	return m.out
}

func (m *merge) VerboseON() {
	m.Verbose = true
}

func (m *merge) VerboseOFF() {
	m.Verbose = false
}

func (m *merge) Emit() {
	defer close(m.out)

	processSource := func(in <-chan Element) {
		for e := range in {
			m.out <- e
		}
	}

	for _, s := range m.sources {
		go s.Emit()
		processSource(s.GetOutput())
	}
}

func (m *merge) Stop() {
	go func() {
		defer close(m.error)
		m.done <- true
	}()
}

func (m *merge) Subscribe(consCh chan<- bool) {
	m.consumers = append(m.consumers, consCh)
}

func MergeSources(sources ...Source) *merge {
	return &merge{
		out:     make(chan Element),
		sources: sources,
		done:    make(chan bool),
		error:   make(chan error),
		Verbose: false,
	}
}
