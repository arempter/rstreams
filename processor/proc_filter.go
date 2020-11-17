package processor

func Filter(in <-chan interface{}, predicate Predicate) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for e := range in {
			if predicate(e) {
				out <- e
			}
		}
		close(out)
	}()
	return out
}
