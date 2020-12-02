## Streams in Go

tests on streams processing based on Go channels  

### usage

```
stream.FromSource(source.Slice(ints)).
		Map(processor.ToInt).
		Map(addOne).
		Filter(even).
		To(sink.Foreach()).
		Run()
```

* A basic stream howto
[sample stream code](./examples/basic_stream.go)

* Kafka Avro consumer 
[sample kafka](./examples/avro_kafka.go)

* Merge multiple sources
[sample merge](./examples/merge_multiple_sources.go)
