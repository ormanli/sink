package sink

import "errors"

var (
	// ErrInvalidMaxItemsForBatching is an error.
	ErrInvalidMaxItemsForBatching = errors.New("max items for batching must be more than zero")
	// ErrNilExpensiveOperation is an error.
	ErrNilExpensiveOperation = errors.New("there is not expensive operation")
	// ErrInvalidMaxTimeoutForBatching is an error.
	ErrInvalidMaxTimeoutForBatching = errors.New("max timeout for matching must be more than zero")
	// ErrInvalidAddPoolSize is an error.
	ErrInvalidAddPoolSize = errors.New("add pool size must be more than zero")
	// ErrInvalidCallbackPoolSize is an error.
	ErrInvalidCallbackPoolSize = errors.New("callback pool size must be more than zero")
	// ErrInvalidExpensivePoolSize is an error.
	ErrInvalidExpensivePoolSize = errors.New("expensive pool size must be more than zero")
)
