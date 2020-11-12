package sink

import (
	"fmt"
	"reflect"
)

func Foreach(in <-chan interface{}) {
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
