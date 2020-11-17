package processor

import (
	"strings"
)

func ToUpper() ProcFuncSpec {
	return ProcFuncSpec{Body: toUpper}
}

func ToLower() ProcFuncSpec {
	return ProcFuncSpec{Body: toLower}
}

func toUpper(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for e := range in {
			switch e.(type) {
			case string:
				out <- strings.ToUpper(e.(string))
			default:
			}
		}
		close(out)
	}()
	return out
}

func toLower(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for e := range in {
			switch e.(type) {
			case string:
				out <- strings.ToLower(e.(string))
			default:
			}
		}
		close(out)
	}()
	return out
}
