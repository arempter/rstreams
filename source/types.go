package source

type Source interface {
	GetOutput() <-chan string
	Emit()
}
