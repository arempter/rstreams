package sink

import (
	"fmt"
	"reflect"
)

type foreach struct {
}

func Foreach() *foreach {
	return &foreach{}
}

func (f foreach) SetOnNextCh(in chan bool) {
	panic("not supported")
}

func (f foreach) HasBackpressure() bool {
	return false
}

func (f foreach) Receive(in <-chan interface{}) {
	for e := range in {
		switch e.(type) {
		case string:
			fmt.Println(e)
		case []byte:
			fmt.Print(string(e.([]byte)))
		default:
			fmt.Println("got other type", reflect.TypeOf(e), e)
		}
	}
}
