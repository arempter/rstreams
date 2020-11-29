package util

import (
	"errors"
	"reflect"
)

func IsFunc(v reflect.Value) error {
	if v.Kind() != reflect.Func {
		return errors.New("argument must be function")
	}
	return nil
}

func IsSlice(v reflect.Value) error {
	if v.Kind() != reflect.Slice {
		return errors.New("argument must be slice type")
	}
	return nil
}

func IsMapFunc(p interface{}) error {
	mVal := reflect.ValueOf(p)
	if err := IsFunc(mVal); err != nil {
		return err
	}
	if mVal.Type().NumIn() != 1 || mVal.Type().NumOut() == 0 || mVal.Type().NumOut() > 1 {
		return errors.New("map function must have one in argument and one return type")
	}
	return nil
}

func IsFilterFunc(p interface{}) error {
	mVal := reflect.ValueOf(p)
	if err := IsFunc(mVal); err != nil {
		return err
	}
	if mVal.Type().NumIn() != 1 || mVal.Type().NumOut() == 0 || mVal.Type().NumOut() > 1 || mVal.Type().Out(0).Kind() != reflect.Bool {
		return errors.New("filter function must have one in argument and one return type (bool)")
	}
	return nil
}
