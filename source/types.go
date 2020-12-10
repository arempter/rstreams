package source

import (
	"time"
)

type Source interface {
	GetOutput() <-chan Element
	Emit()
	Stop()
	ErrorCh() <-chan error
	Subscribe(consCh chan<- bool)
	VerboseON()
	VerboseOFF()
}

type Element struct {
	Payload   interface{}
	Timestamp time.Time
}
