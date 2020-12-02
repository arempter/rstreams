package processor

import (
	"strconv"
	"strings"
)

var (
	ToString = func(i interface{}) string {
		var res string
		res, ok := i.(string)
		if !ok {
			panic("failed to parse element as string")
		}
		return res
	}
	ToInt = func(i interface{}) int {
		s := i.(string)
		res, err := strconv.Atoi(s)
		if err != nil {
			panic("failed to parse element as int")
		}
		return res
	}
	ToUpper = func(e string) string {
		return strings.ToUpper(e)
	}
	ToLower = func(e string) string {
		return strings.ToLower(e)
	}
)
