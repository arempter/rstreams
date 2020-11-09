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

[sample stream code](./examples/sample_stream.go)
