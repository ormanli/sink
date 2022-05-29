package sink

type request[I, O any] struct {
	value    I
	callback chan response[O]
}

type response[O any] struct {
	callback chan response[O]
	value    O
	err      error
}

func newItem[I, O any](value I) request[I, O] {
	return request[I, O]{
		value:    value,
		callback: make(chan response[O]),
	}
}

type expensiveOperation[I, O any] func([]I) ([]O, error)
