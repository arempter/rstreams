package processors

func Filter(in <-chan interface{}, predicate Predicate) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for e := range in {
			switch e.(type) {
			case string:
				if predicate(e.(string)) {
					out <- e.(string)
				}
			default:
			}
		}
		close(out)
	}()
	return out
}
