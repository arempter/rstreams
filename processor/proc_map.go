package processor

import (
	"reflect"
	"rstreams/util"
)

func Map(in <-chan interface{}, predicate interface{}, out chan interface{}) {
	if err := util.IsMapFunc(predicate); err != nil {
		panic(err.Error())
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				Map(in, predicate, out)
			}
		}()
		for e := range in {
			if e != nil {
				eVal := reflect.ValueOf(e)
				out <- reflect.ValueOf(predicate).Call([]reflect.Value{eVal})[0].Interface()
			}
		}
		//close(out)
	}()
}
