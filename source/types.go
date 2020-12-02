package source

import (
	"time"
)

type Source interface {
	GetOutput() <-chan Element
	OnNextCh() chan bool
	Emit()
	Stop()
	ErrorCh() <-chan error
	Subscribe(consCh chan<- bool)
}

type Element struct {
	Payload   interface{}
	Timestamp time.Time
}
