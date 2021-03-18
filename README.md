# sink

sink allows you to batch expensive operations using items from different flows.

```go
s, err := sink.NewSink(sink.Config{
    MaxItemsForBatching:   10,
    MaxTimeoutForBatching: 10 * time.Millisecond,
    AddPoolSize:           10,
    CallbackPoolSize:      10,
    ExpensivePoolSize:     10,
    ExpensiveOperation: func(i []interface{}) ([]interface{}, error) {
        time.Sleep(time.Second)

        return i, nil
    },
})
defer s.Close()

_, err = s.Add(dummy{i: 10})
```
