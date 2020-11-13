package source

type Source interface {
	GetOutput() <-chan interface{}
	Emit()
	Stop()
	GetErrorCh() <-chan error
}
