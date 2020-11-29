package processor

import (
	"reflect"
	"rstreams/util"
)

func Map(in <-chan interface{}, predicate interface{}) <-chan interface{} {
	if err := util.IsMapFunc(predicate); err != nil {
		panic(err.Error())
	}
	out := make(chan interface{})
	go func() {
		for e := range in {
			if e != nil {
				eVal := reflect.ValueOf(e)
				out <- reflect.ValueOf(predicate).Call([]reflect.Value{eVal})[0].Interface()
			}
		}
		close(out)
	}()
	return out
}
