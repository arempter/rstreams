package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"rstreams/processor"
	"rstreams/sink"
	"rstreams/source"
	"rstreams/stream"
	"runtime"
	"strings"
	"time"
)

func main() {
	//words := strings.Split("Now you've got everything running, you may wonder: what now? This section we'll describe a basic flow on how to use Rokku to perform operations in S3. You may refer to the What is Rokku? document before diving in here. That will introduce you to the various components used.", " ")

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	//runtime.SetBlockProfileRate(1)
	wordsFile, err := os.Open("little_princ.txt")
	if err != nil {
		fmt.Println("Cannot open file", err)
	}
	data := make([]byte, 100000)
	dataCount, err := wordsFile.Read(data)
	if err != nil {
		fmt.Println("Cannot read file", err)
	}
	fmt.Println("read", dataCount)
	words := strings.Split(string(data), " ")
	var wordsBig []string
	for a := 0; a < 100; a++ {
		wordsBig = append(wordsBig, words...)
	}

	fmt.Println(len(wordsBig))
	//go func() {
	//	for {
	//		fmt.Println("no of go routines", runtime.NumGoroutine())
	//		time.Sleep(1*time.Second)
	//	}
	//}()

	stream := stream.FromSource(source.Slice(wordsBig))

	containsStringFunc := func(i string) bool {
		if !strings.Contains(strings.ToLower(i), "a") {
			return false
		}
		return true
	}

	runtime.GOMAXPROCS(4)
	start := time.Now()
	stream.
		Filter(containsStringFunc).
		Map(processor.ToUpper).
		CountByWindow(250 * time.Millisecond).
		To(sink.Ignore()).
		Run()

	fmt.Println("took", time.Since(start))
}
