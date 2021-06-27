package sink

import "time"

// Config holds configuration for the sink.
type Config struct {
	MaxItemsForBatching   int
	MaxTimeoutForBatching time.Duration

	Logger Logger

	AddPoolSize        int
	CallbackPoolSize   int
	ExpensivePoolSize  int
	ExpensiveOperation expensiveOperation
}

func validateConfig(config Config) error {
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
