package processor

import "strings"

var (
	//todo conv err
	ToString = func(i interface{}) string {
		return i.(string)
	}
	ToInt = func(i interface{}) int {
		return i.(int)
	}
	ToUpper = func(e string) string {
		return strings.ToUpper(e)
	}
	ToLower = func(e string) string {
		return strings.ToLower(e)
	}
)
