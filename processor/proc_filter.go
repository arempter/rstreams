package processor

import (
	"reflect"
	"rstreams/util"
)

func Filter(in <-chan interface{}, predicate interface{}, out chan interface{}) {
	if err := util.IsFilterFunc(predicate); err != nil {
		panic(err.Error())
	}
	go func() {
		for e := range in {
			if e != nil {
				eVal := reflect.ValueOf(e)
				if reflect.ValueOf(predicate).Call([]reflect.Value{eVal})[0].Bool() {
					out <- e
				}
			}
		}
	}()
}
