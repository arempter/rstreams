package source

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type httpSource struct {
	urls       []string
	out        chan Element
	done       chan bool
	error      chan error
	consumers  []chan<- bool
	Verbose    bool
	rate       time.Duration
	burstLimit int
}

func (h *httpSource) ErrorCh() <-chan error {
	return h.error
}

func (h *httpSource) VerboseON() {
	h.Verbose = true
}

func (h *httpSource) VerboseOFF() {
	h.Verbose = false
}

func Http(urls []string, rate time.Duration) *httpSource {
	return &httpSource{
		urls:       urls,
		out:        make(chan Element),
		done:       make(chan bool),
		error:      make(chan error),
		Verbose:    false,
		rate:       rate, // req per second
		burstLimit: 10,
	}
}

func (h *httpSource) GetOutput() <-chan Element {
	return h.out
}

func (h *httpSource) Emit() {
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
	throttle := make(chan time.Time, h.burstLimit)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(h.rate)
		defer ticker.Stop()

		for t := range ticker.C {
			select {
			case throttle <- t:
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-h.done:
			return
		default:
			for _, u := range h.urls {
				<-throttle
				go request(u)
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

func (h *httpSource) sendToErr(e string) {
	if h.Verbose {
		h.error <- errors.New(fmt.Sprintf(e))
	}
}
