package processor

func Map(in <-chan interface{}, predicate MapF) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for e := range in {
			out <- predicate(e)
		}
		close(out)
	}()
	return out
}
