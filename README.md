# sink

sink allows you to batch expensive operations using items from different flows.

![GitHub](https://img.shields.io/github/license/ormanli/sink)
![Codecov](https://img.shields.io/codecov/c/github/ormanli/sink)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ormanli/sink/test)

## Installation

Install using go modules

```bash 
  go get -u github.com/ormanli/sink
```

## Usage/Examples

### Configuration

```go
cfg := sink.Config{
    MaxItemsForBatching:   10, // Maximum number of items to batch inputs, mandatory, can't be less than 1.
    MaxTimeoutForBatching: 10 * time.Millisecond, // Maximum time to wait for inputs, mandatory, can't be less than 1 millisecond.
    AddPoolSize:           10, // Add operation goroutine pool size, mandatory, can't be less than 1.
    CallbackPoolSize:      10, // Callback operation goroutine pool size, mandatory, can't be less than 1.
    ExpensivePoolSize:     10, // Expensive operation goroutine pool size, mandatory, can't be less than 1.
    ExpensiveOperation: func(i []interface{}) ([]interface{}, error) {
        time.Sleep(time.Second)

        return i, nil
    }, // Actual function that is called with batched items, mandatory.
    Logger:                 customLogger, // Logger is optional, if not provided log package used.
}
```
Sink will either wait until MaxItemsForBatching of items to arrive or wait until MaxTimeoutForBatching to start processing batched items.


### Sink

Processes given inputs by provided configuration. When it is no longer required, stop sink by calling `Close` method.

```go
s, err := sink.NewSink(cfg)
defer s.Close()

_, err = s.Add(dummy{i: 10})
```

### Sink With Context

Processes given inputs by provided configuration. It will run until the given context is canceled.

```go
ctx, cncl := context.WithCancel(context.Background())
s, err := sink.NewSinkWithContext(ctx, cfg)
defer cncl()

_, err = s.Add(dummy{i: 10})
```

## Running Tests

To run tests, run the following command

```bash
  make test
```

## Authors

- [ormanli](https://www.github.com/ormanli)
