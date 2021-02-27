package sink

type request struct {
	value    interface{}
	callback chan response
}

type response struct {
	callback chan response
	value    interface{}
	err      error
}

func newItem(value interface{}) request {
	return request{
		value:    value,
		callback: make(chan response),
	}
}

type expensiveOperation func([]interface{}) ([]interface{}, error)
