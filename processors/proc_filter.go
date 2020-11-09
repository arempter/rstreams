package processors

func Filter(in <-chan string, predicate Predicate) <-chan string {
	out := make(chan string)
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
