package sink

import (
	"time"

	"github.com/panjf2000/ants/v2"
)

// Sink is a struct to process different request simultaneously.
type Sink struct {
	input         chan request
	addPool       *ants.PoolWithFunc
	callbackPool  *ants.PoolWithFunc
	expensivePool *ants.PoolWithFunc
	logger        Logger
}

// NewSink initializes a sink with the provided config.
func NewSink(config Config) (*Sink, error) {
	err := validateConfig(config)
	if err != nil {
		return nil, err
	}

	if config.Logger == nil {
		config.Logger = standardLogger{}
	}

	s := &Sink{logger: config.Logger}
	s.input = make(chan request)

	options := []ants.Option{ants.WithLogger(config.Logger)}

	s.addPool, err = ants.NewPoolWithFunc(config.AddPoolSize, s.addFunc, options...)
	if err != nil {
		return nil, err
	}

	s.callbackPool, err = ants.NewPoolWithFunc(config.CallbackPoolSize, s.callbackFunc, options...)
	if err != nil {
		return nil, err
	}

	s.expensivePool, err = ants.NewPoolWithFunc(config.ExpensivePoolSize, s.expensiveWrapper(config.ExpensiveOperation), options...)
	if err != nil {
		return nil, err
	}

	batches := batch(s.input, config.MaxItemsForBatching, config.MaxTimeoutForBatching)

	go func(batches chan []request) {
		for batch := range batches {
			err := s.expensivePool.Invoke(batch)
			if err != nil {
				s.logger.Printf("error: %v", err)
			}
		}
	}(batches)

	return s, nil
}

// Add adds a value to the sink and waits for result.
func (s *Sink) Add(value interface{}) (interface{}, error) {
	rqs := newItem(value)

	err := s.addPool.Invoke(rqs)
	if err != nil {
		return nil, err
	}

	rsp := <-rqs.callback

	return rsp.value, rsp.err
}

// Close closes sink to stop processing.
func (s *Sink) Close() {
	close(s.input)

	s.addPool.Release()
	s.expensivePool.Release()
	s.callbackPool.Release()
}

func (s *Sink) addFunc(i interface{}) {
	rq := i.(request)

	s.input <- rq
}

func (s *Sink) callbackFunc(i interface{}) {
	drsp := i.(response)

	drsp.callback <- drsp

	close(drsp.callback)
}

func (s *Sink) expensiveWrapper(f expensiveOperation) func(i interface{}) {
	return func(i interface{}) {
		batch := i.([]request)

		values := make([]interface{}, len(batch))
		for k := range values {
			values[k] = batch[k].value
		}

		responses, err := f(values)

		for k := range responses {
			r := response{
				callback: batch[k].callback,
			}

			if err != nil {
				r.err = err
			} else {
				r.value = responses[k]
			}

			err := s.callbackPool.Invoke(r)
			if err != nil {
				s.logger.Printf("error: %v", err)
			}
		}
	}
}

func batch(values <-chan request, maxItems int, maxTimeout time.Duration) chan []request {
	batches := make(chan []request)

	go func(batches chan []request) {
		defer close(batches)

		for keepGoing := true; keepGoing; {
			var batch []request
			expire := time.After(maxTimeout)
			for {
				select {
				case value, ok := <-values:
					if !ok {
						keepGoing = false
						goto done
					}

					batch = append(batch, value)
					if len(batch) == maxItems {
						goto done
					}

				case <-expire:
					goto done
				}
			}

		done:
			if len(batch) > 0 {
				batches <- batch
			}
		}
	}(batches)

	return batches
}
