package processor

import (
	"reflect"
	"rstreams/source"
	"rstreams/util"
)

func Map(in <-chan source.Element, predicate interface{}, out chan source.Element) {
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
			if e.Payload != nil {
				eVal := reflect.ValueOf(e.Payload)
				out <- source.Element{
					Payload:   reflect.ValueOf(predicate).Call([]reflect.Value{eVal})[0].Interface(),
					Timestamp: e.Timestamp,
				}
			}
		}
		//close(out)
	}()
}
