package test

import (
	"context"
	"testing"

	"github.com/arthureichelberger/goboiler/pkg/psql"
	"github.com/stretchr/testify/require"
)

func Txn(t *testing.T, ctx context.Context, db psql.DBGetter) (psql.DBGetter, func()) {
	t.Helper()

	txn, err := db(ctx).BeginTxx(ctx, nil)
	require.NoError(t, err)

	txQueryable := &psql.TxQueryable{Queryable: txn}

	return func(ctx context.Context) psql.Queryable {
			if subTxn := psql.TxFromContext(ctx); subTxn != nil {
				return subTxn
			}

			return txQueryable
		}, func() {
			// nolint:errcheck
			txn.Rollback()
		}
}
