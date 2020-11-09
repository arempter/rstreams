package processors

import (
	"strings"
)

func ToUpper(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for e := range in {
			out <- strings.ToUpper(e)
		}
		close(out)
	}()
	return out
}

func ToLower(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for e := range in {
			out <- strings.ToLower(e)
		}
		close(out)
	}()
	return out
}
