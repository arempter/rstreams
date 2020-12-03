package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"runtime"
	"time"
)

func main() {
	sampleUrls := []string{"http://localhost:1234/device/1", "http://localhost:1234/device/2", "http://localhost:1234/device/3", "http://localhost:1234/device/4"}

	stream := stream.FromSource(source.Http(sampleUrls))

	//go func() {
	//	time.Sleep(10 * time.Second)
	//	stream.Stop()
	//	os.Exit(0)
	//}()

	onlyValidJson := func(s interface{}) bool {
		if len(s.([]byte)) > 5 {
			return true
		}
		return false
	}

	mapSlow := func(e interface{}) interface{} {
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		return e
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("go routines", runtime.NumGoroutine())
		}
	}()
	go stream.WireTap()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	stream.
		Map(mapSlow).
		Filter(onlyValidJson).
		To(sink.Foreach()).
		Run()
}
