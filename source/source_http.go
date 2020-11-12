package source

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type httpSource struct {
	urls     []string
	out      chan interface{}
	done     chan bool
	runEvery time.Duration
}

func FromHttp(urls []string) *httpSource {
	return &httpSource{
		urls:     urls,
		out:      make(chan interface{}, 5),
		done:     make(chan bool),
		runEvery: 1 * time.Second,
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

	tick := time.NewTicker(h.runEvery)
	run := true
	for run == true {
		select {
		case <-h.done:
			run = false
		case <-tick.C:
			for _, u := range h.urls {
				go request(u)
			}
		}
	}
}

func (h httpSource) Stop() {
	log.Println("sending stop signal to http source...")
	h.done <- true
}
