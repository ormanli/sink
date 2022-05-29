package sink

import (
	"time"

	"github.com/panjf2000/ants/v2"
)

// Sink is a struct to process different request simultaneously.
type Sink[I, O any] struct {
	input         chan request[I, O]
	addPool       *ants.PoolWithFunc
	callbackPool  *ants.PoolWithFunc
	expensivePool *ants.PoolWithFunc
	logger        Logger
}

// NewSink initializes a sink with the provided config.
func NewSink[I, O any](config Config[I, O]) (*Sink[I, O], error) {
	err := validateConfig(config)
	if err != nil {
		return nil, err
	}

	if config.Logger == nil {
		config.Logger = standardLogger{}
	}

	s := &Sink[I, O]{logger: config.Logger}
	s.input = make(chan request[I, O])

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

	go func(batches chan []request[I, O]) {
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
func (s *Sink[I, O]) Add(value I) (O, error) {
	rqs := newItem[I, O](value)

	err := s.addPool.Invoke(rqs)
	if err != nil {
		var result O
		return result, err
	}

	rsp := <-rqs.callback

	return rsp.value, rsp.err
}

// Close closes sink to stop processing.
func (s *Sink[I, O]) Close() {
	close(s.input)

	s.addPool.Release()
	s.expensivePool.Release()
	s.callbackPool.Release()
}

func (s *Sink[I, O]) addFunc(i interface{}) {
	rq := i.(request[I, O])

	s.input <- rq
}

func (s *Sink[I, O]) callbackFunc(i interface{}) {
	drsp := i.(response[O])

	drsp.callback <- drsp

	close(drsp.callback)
}

func (s *Sink[I, O]) expensiveWrapper(f expensiveOperation[I, O]) func(i interface{}) {
	return func(i interface{}) {
		batch := i.([]request[I, O])

		values := make([]I, len(batch))
		for k := range values {
			values[k] = batch[k].value
		}

		responses, err := f(values)

		for k := range responses {
			r := response[O]{
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

func batch[I, O any](values <-chan request[I, O], maxItems int, maxTimeout time.Duration) chan []request[I, O] {
	batches := make(chan []request[I, O])

	go func(batches chan []request[I, O]) {
		defer close(batches)

		for keepGoing := true; keepGoing; {
			var batch []request[I, O]
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
