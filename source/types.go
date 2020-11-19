package source

type Source interface {
	GetOutput() <-chan interface{}
	OnNextCh() chan bool
	Emit()
	Stop()
	ErrorCh() <-chan error
	Subscribe(consCh chan<- bool)
}
