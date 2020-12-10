package processor

import (
	"reflect"
	"rstreams/source"
	"rstreams/util"
)

func Filter(in <-chan source.Element, predicate interface{}, out chan source.Element, par int) {
	if err := util.IsFilterFunc(predicate); err != nil {
		panic(err.Error())
	}
	go func() {
		for e := range in {
			if e.Payload != nil {
				eVal := reflect.ValueOf(e.Payload)
				if reflect.ValueOf(predicate).Call([]reflect.Value{eVal})[0].Bool() {
					out <- e
				}
			}
		}
		close(out)
	}()
}
