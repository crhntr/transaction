# A Transaction Manager Utility

This package lets you do stuff like this [(full example here)](https://github.com/crhntr/muxt-example-htmx-sortable/blob/main/internal/database/tx.go):

```go
package database

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/crhntr/transaction"
)


type Caller interface {
	Call(ctx context.Context, options pgx.TxOptions, p transaction.Func) error
}

type Transactions struct {
	manager Caller
}

func NewTransactions(conn transaction.Beginner) *Transactions {
	return NewTransactionsWithCaller(transaction.NewManager(conn))
}

func NewTransactionsWithCaller(m Caller) *Transactions {
	return &Transactions{manager: m}
}

func (t Transactions) ReadOnly(ctx context.Context, f ReadOnlyFunc) error {
	o := pgx.TxOptions{AccessMode: pgx.ReadOnly}
	return t.manager.Call(ctx, o, f.Func)
}

func (t Transactions) UpdatePriorityList(ctx context.Context, f TaskPriorityUpdateFunc) error {
	o := pgx.TxOptions{AccessMode: pgx.ReadWrite, DeferrableMode: pgx.Deferrable}
	return t.manager.Call(ctx, o, f.Func)
}

type ReadOnlyFunc func(ReadOnlyQuerier) error

func (f ReadOnlyFunc) Func(_ context.Context, tx pgx.Tx) error { return f(New(tx)) }

type TaskPriorityUpdateFunc func(TaskPriorityUpdater) error

func (f TaskPriorityUpdateFunc) Func(ctx context.Context, tx pgx.Tx) error {
	const statement = `SET CONSTRAINTS unique_list_priority DEFERRED;`
	if _, err := tx.Exec(ctx, statement); err != nil {
		return err
	}
	return f(New(tx))
}

type ReadOnlyQuerier interface {
	ListByID(ctx context.Context, id int32) (List, error)
	Lists(ctx context.Context) ([]List, error)
	TasksByListID(ctx context.Context, id int32) ([]TasksByListIDRow, error)
}

type TaskPriorityUpdater interface {
	ListByID(ctx context.Context, id int32) (List, error)
	SetTaskPriority(ctx context.Context, arg SetTaskPriorityParams) error
}
```
