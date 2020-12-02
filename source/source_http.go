package source

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type httpSource struct {
	urls      []string
	out       chan Element
	onNext    chan bool
	done      chan bool
	error     chan error
	consumers []chan<- bool
}

func (h httpSource) ErrorCh() <-chan error {
	return h.error
}

func (h httpSource) OnNextCh() chan bool {
	return h.onNext
}

func Http(urls []string) *httpSource {
	return &httpSource{
		urls:   urls,
		out:    make(chan Element),
		onNext: make(chan bool),
		done:   make(chan bool),
		error:  make(chan error),
	}
}

func (h httpSource) GetOutput() <-chan Element {
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
		h.out <- Element{
			Payload:   body,
			Timestamp: time.Now(),
		}
		return nil
	}

	run := true
	for run == true {
		select {
		case <-h.done:
			run = false
		// only emit on signal from consumer
		case <-h.onNext:
			h.error <- errors.New(fmt.Sprintf("source => got demand signal"))
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

func (h *httpSource) Subscribe(consCh chan<- bool) {
	h.consumers = append(h.consumers, consCh)
}
