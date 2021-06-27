package sink_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/golang/mock/gomock"
	"github.com/ormanli/sink"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func Test_100ItemsIn10BatchesWithContext(t *testing.T) {
	defer leaktest.CheckTimeout(t, 5*time.Second)()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	op := NewMockOperation(ctrl)
	op.EXPECT().
		Op(gomock.Len(10)).
		Times(10).
		DoAndReturn(func(i []interface{}) ([]interface{}, error) {
			return i, nil
		})

	ctx, cncl := context.WithCancel(context.Background())

	s, err := sink.NewSinkWithContext(ctx, sink.Config{
		MaxItemsForBatching:   10,
		MaxTimeoutForBatching: 10 * time.Millisecond,
		AddPoolSize:           10,
		CallbackPoolSize:      10,
		ExpensivePoolSize:     10,
		ExpensiveOperation:    op.Op,
	})

	require.NoError(t, err)
	defer cncl()

	var g errgroup.Group

	for i := 0; i < 100; i++ {
		i := i
		g.Go(func() error {
			_, err := s.Add(dummy{i: i})

			return err
		})
	}

	err = g.Wait()
	require.NoError(t, err)
}

func Test_ErrorFromExpensiveOperationWithContext(t *testing.T) {
	defer leaktest.CheckTimeout(t, 5*time.Second)()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	op := NewMockOperation(ctrl)
	op.EXPECT().
		Op(gomock.Len(10)).
		DoAndReturn(func(i []interface{}) ([]interface{}, error) {
			return i, errors.New("expensive operation failed")
		})

	ctx, cncl := context.WithCancel(context.Background())

	s, err := sink.NewSinkWithContext(ctx, sink.Config{
		MaxItemsForBatching:   10,
		MaxTimeoutForBatching: time.Second,
		AddPoolSize:           10,
		CallbackPoolSize:      10,
		ExpensivePoolSize:     10,
		ExpensiveOperation:    op.Op,
	})

	require.NoError(t, err)
	defer cncl()

	var g errgroup.Group

	for i := 0; i < 10; i++ {
		i := i
		g.Go(func() error {
			_, err := s.Add(dummy{i: i})

			return err
		})
	}

	err = g.Wait()
	require.EqualError(t, err, "expensive operation failed")
}
