package sink_test

//go:generate mockgen -source=sink_test.go -package=sink_test -destination=operation_mock_test.go

import (
	"errors"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/golang/mock/gomock"
	"github.com/ormanli/sink"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

type Operation[I, O any] interface {
	Op(I) (O, error)
}

type dummy struct {
	i int
}

func Test_100ItemsIn10Batches(t *testing.T) {
	defer leaktest.CheckTimeout(t, 5*time.Second)()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	op := NewMockOperation[[]dummy, []dummy](ctrl)
	op.EXPECT().
		Op(gomock.Len(10)).
		Times(10).
		DoAndReturn(func(i []dummy) ([]dummy, error) {
			return i, nil
		})

	s, err := sink.NewSink[dummy, dummy](sink.Config[dummy, dummy]{
		MaxItemsForBatching:   10,
		MaxTimeoutForBatching: 10 * time.Millisecond,
		AddPoolSize:           10,
		CallbackPoolSize:      10,
		ExpensivePoolSize:     10,
		ExpensiveOperation:    op.Op,
	})

	require.NoError(t, err)
	defer s.Close()

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

func Test_ErrorFromExpensiveOperation(t *testing.T) {
	defer leaktest.CheckTimeout(t, 5*time.Second)()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	op := NewMockOperation[[]dummy, []dummy](ctrl)
	op.EXPECT().
		Op(gomock.Len(10)).
		DoAndReturn(func(i []dummy) ([]dummy, error) {
			return i, errors.New("expensive operation failed")
		})

	s, err := sink.NewSink[dummy, dummy](sink.Config[dummy, dummy]{
		MaxItemsForBatching:   10,
		MaxTimeoutForBatching: time.Millisecond,
		AddPoolSize:           10,
		CallbackPoolSize:      10,
		ExpensivePoolSize:     10,
		ExpensiveOperation:    op.Op,
	})

	require.NoError(t, err)
	defer s.Close()

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
