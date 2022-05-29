package sink

import "time"

// Config holds configuration for the sink.
type Config[I, O any] struct {
	MaxItemsForBatching   int
	MaxTimeoutForBatching time.Duration

	Logger Logger

	AddPoolSize        int
	CallbackPoolSize   int
	ExpensivePoolSize  int
	ExpensiveOperation expensiveOperation[I, O]
}

func validateConfig[I, O any](config Config[I, O]) error {
	if config.MaxItemsForBatching < 1 {
		return ErrInvalidMaxItemsForBatching
	}

	if config.ExpensiveOperation == nil {
		return ErrNilExpensiveOperation
	}

	if config.MaxTimeoutForBatching < time.Millisecond {
		return ErrInvalidMaxTimeoutForBatching
	}

	if config.AddPoolSize < 1 {
		return ErrInvalidAddPoolSize
	}

	if config.CallbackPoolSize < 1 {
		return ErrInvalidCallbackPoolSize
	}

	if config.ExpensivePoolSize < 1 {
		return ErrInvalidExpensivePoolSize
	}

	return nil
}
