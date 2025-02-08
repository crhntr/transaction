package transaction_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/crhntr/transaction"
	"github.com/crhntr/transaction/internal/fake"
)

func TestManager(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			conn = new(fake.Beginner)
			tx   = new(fake.Tx)
			f    = new(fake.Func)
		)
		conn.BeginTxReturns(tx, nil)

		m := transaction.NewManager(conn)

		err := m.Call(ctx, pgx.TxOptions{}, f.Spy)

		require.NoError(t, err)
		require.Equal(t, 1, tx.CommitCallCount())
		require.Zero(t, tx.RollbackCallCount())
		require.Equal(t, 1, f.CallCount())
	})

	t.Run("BeginTx fails", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			conn = new(fake.Beginner)
			tx   = new(fake.Tx)
			f    = new(fake.Func)
		)
		conn.BeginTxReturns(tx, fmt.Errorf("banana"))

		m := transaction.NewManager(conn)
		err := m.Call(ctx, pgx.TxOptions{}, f.Spy)

		require.ErrorContains(t, err, "banana")
		require.Zero(t, tx.CommitCallCount())
		require.Zero(t, tx.RollbackCallCount())
		require.Zero(t, tx.ExecCallCount())
		require.Zero(t, tx.BeginCallCount())
	})
	t.Run("f fails", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			conn = new(fake.Beginner)
			tx   = new(fake.Tx)
			f    = new(fake.Func)
		)
		conn.BeginTxReturns(tx, nil)
		f.Returns(fmt.Errorf("banana"))

		m := transaction.NewManager(conn)
		err := m.Call(ctx, pgx.TxOptions{}, f.Spy)

		require.ErrorContains(t, err, "banana")
		require.Equal(t, 1, tx.RollbackCallCount())
		require.Zero(t, tx.CommitCallCount())
	})
	t.Run("f calls panic with non error", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			conn = new(fake.Beginner)
			tx   = new(fake.Tx)
			f    = new(fake.Func)
		)
		conn.BeginTxReturns(tx, nil)
		f.Calls(func(context.Context, pgx.Tx) error { panic("lemon") })

		m := transaction.NewManager(conn)
		err := m.Call(ctx, pgx.TxOptions{}, f.Spy)

		require.ErrorContains(t, err, ": lemon")
		require.Equal(t, 1, tx.RollbackCallCount())
	})
	t.Run("f calls panic with error", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			conn = new(fake.Beginner)
			tx   = new(fake.Tx)
			f    = new(fake.Func)
		)
		conn.BeginTxReturns(tx, nil)
		f.Calls(func(ctx context.Context, tx pgx.Tx) error {
			panic("lemon")
		})

		m := transaction.NewManager(conn)
		err := m.Call(ctx, pgx.TxOptions{}, f.Spy)

		require.ErrorContains(t, err, "lemon")
		require.Equal(t, 1, tx.RollbackCallCount())
	})
	t.Run("panic and rollback messages", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			conn = new(fake.Beginner)
			tx   = new(fake.Tx)
			f    = new(fake.Func)
		)

		conn.BeginTxReturns(tx, nil)
		tx.RollbackCalls(func(ctx context.Context) error {
			if tx.RollbackCallCount() == 1 {
				return errors.New("lemon")
			}
			return nil
		})
		f.Calls(func(ctx context.Context, tx pgx.Tx) error {
			panic("banana")
		})

		m := transaction.NewManager(conn)
		err := m.Call(ctx, pgx.TxOptions{}, f.Spy)

		require.ErrorContains(t, err, "banana")
		require.ErrorContains(t, err, "lemon")
		require.Equal(t, 1, tx.RollbackCallCount())
	})
	t.Run("rollback and error message", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			conn = new(fake.Beginner)
			tx   = new(fake.Tx)
			f    = new(fake.Func)
		)
		conn.BeginTxReturns(tx, nil)
		tx.RollbackCalls(func(ctx context.Context) error {
			if tx.RollbackCallCount() == 1 {
				return errors.New("lemon")
			}
			return nil
		})
		f.Calls(func(ctx context.Context, tx pgx.Tx) error { return errors.New("banana") })

		m := transaction.NewManager(conn)
		err := m.Call(ctx, pgx.TxOptions{}, f.Spy)

		require.ErrorContains(t, err, "banana")
		require.ErrorContains(t, err, "lemon")
		require.Equal(t, 1, tx.RollbackCallCount())
	})
}
