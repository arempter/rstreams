package processors

import (
	"strings"
)

func ToUpper(in <-chan interface{}) <-chan interface{} {
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

func ToLower(in <-chan interface{}) <-chan interface{} {
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
