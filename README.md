## Streams in Go

tests on streams processing based on Go channels  

### usage

```
stream.
       Filter(processors.Filter, containsStringFunc).
       Via(processors.ToUpper).
       To(sinks.ForeachSink).
    Run()
```

* Sample basic stream
[sample stream code](./examples/sample_stream.go)

* Kafka consumer 
[sample kafka](./examples/kafka_stream.go)

* Merge streams
[sample merge](./examples/sample_merge.go)
