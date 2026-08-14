package querier

import (
	"context"
	"database/sql"
)

type Querier interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
