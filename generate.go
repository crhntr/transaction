package transaction

//go:generate rm -rf internal/fake
//go:generate mkdir -p internal/fake
//go:generate counterfeiter -generate

//counterfeiter:generate -o internal/fake/tx.go --fake-name Tx github.com/jackc/pgx/v5.Tx

//counterfeiter:generate -o internal/fake/transaction_func.go --fake-name Func       . Func
//counterfeiter:generate -o internal/fake/beginner.go         --fake-name Beginner . Beginner
