package source

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type httpSource struct {
	urls     []string
	out      chan interface{}
	onNext   chan bool
	done     chan bool
	runEvery time.Duration
	error    chan error
}

func (h httpSource) ErrorCh() <-chan error {
	return h.error
}

func (h httpSource) OnNextCh() chan bool {
	return h.onNext
}

func Http(urls []string) *httpSource {
	return &httpSource{
		urls:     urls,
		out:      make(chan interface{}, 20),
		onNext:   make(chan bool),
		done:     make(chan bool),
		runEvery: 1 * time.Second,
		error:    make(chan error),
	}
}

func (h httpSource) GetOutput() <-chan interface{} {
	return h.out
}

func (h httpSource) Emit() {
	request := func(url string) error {
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Request GET failed", err)
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read response")
			return err
		}
		h.out <- body
		return nil
	}

	run := true
	for run == true {
		select {
		case <-h.done:
			run = false
		// only emit on signal from consumer
		case <-h.onNext:
			fmt.Println("demand signal")
			for _, u := range h.urls {
				request(u)
			}
		}
	}
}

func (h httpSource) Stop() {
	log.Println("sending stop signal to http source...")
	h.done <- true
}
