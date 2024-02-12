package test

import (
	"context"
	"testing"

	"github.com/arthureichelberger/goboiler/pkg/psql"
	"github.com/stretchr/testify/require"
)

func Txn(t *testing.T, ctx context.Context, db psql.Queryable) (psql.Queryable, func()) {
	t.Helper()

	txn, err := db.BeginTxx(ctx, nil)
	require.NoError(t, err)

	rollback := func() {
		_ = txn.Rollback() // nolint:errcheck
	}

	return txn, rollback
}
