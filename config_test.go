package sink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		c         Config[any, any]
		errString string
	}{
		{
			name:      "MaxItemsForBatching",
			c:         Config[any, any]{},
			errString: "max items for batching must be more than zero",
		},
		{
			name: "ExpensiveOperation",
			c: Config[any, any]{
				MaxItemsForBatching: 1,
			},
			errString: "there is not expensive operation",
		},
		{
			name: "MaxTimeoutForBatching",
			c: Config[any, any]{
				MaxItemsForBatching: 1,
				ExpensiveOperation: func(i []any) ([]any, error) {
					return nil, nil
				},
			},
			errString: "max timeout for matching must be more than 1 millisecond",
		},
		{
			name: "AddPoolSize",
			c: Config[any, any]{
				MaxItemsForBatching: 1,
				ExpensiveOperation: func(i []any) ([]any, error) {
					return nil, nil
				},
				MaxTimeoutForBatching: time.Millisecond,
			},
			errString: "add pool size must be more than zero",
		},
		{
			name: "CallbackPoolSize",
			c: Config[any, any]{
				MaxItemsForBatching: 1,
				ExpensiveOperation: func(i []any) ([]any, error) {
					return nil, nil
				},
				MaxTimeoutForBatching: time.Millisecond,
				AddPoolSize:           1,
			},
			errString: "callback pool size must be more than zero",
		},
		{
			name: "ExpensivePoolSize",
			c: Config[any, any]{
				MaxItemsForBatching: 1,
				ExpensiveOperation: func(i []any) ([]any, error) {
					return nil, nil
				},
				MaxTimeoutForBatching: time.Millisecond,
				AddPoolSize:           1,
				CallbackPoolSize:      1,
			},
			errString: "expensive pool size must be more than zero",
		},
		{
			name: "Valid",
			c: Config[any, any]{
				MaxItemsForBatching: 1,
				ExpensiveOperation: func(i []any) ([]any, error) {
					return nil, nil
				},
				MaxTimeoutForBatching: time.Millisecond,
				AddPoolSize:           1,
				CallbackPoolSize:      1,
				ExpensivePoolSize:     1,
			},
			errString: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateConfig(test.c)
			if test.errString == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.errString)
			}
		})
	}
}
