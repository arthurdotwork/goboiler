package psql_test

import (
	"context"
	"testing"

	"github.com/arthureichelberger/goboiler/pkg/psql"
	"github.com/arthureichelberger/goboiler/pkg/test"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestWithinTransaction(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	defer rollback()

	_, err = txn(ctx).ExecContext(ctx, `CREATE TABLE transactor_transactions (id SERIAL PRIMARY KEY, name TEXT NOT NULL);`)
	require.NoError(t, err)

	transactor := psql.NewTransactor(txn(ctx))

	t.Run("it should rollback the transaction", func(t *testing.T) {
		err = transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			_, err := txn(ctx).ExecContext(ctx, `INSERT INTO transactor_transactions DEFAULT VALUES;`)
			require.Error(t, err)

			return err
		})
		require.Error(t, err)

		var count int
		err = txn(ctx).GetContext(ctx, &count, `SELECT COUNT(*) FROM transactor_transactions`)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("it should commit the transaction", func(t *testing.T) {
		err = transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			_, err = txn(ctx).ExecContext(ctx, `INSERT INTO transactor_transactions (name) VALUES ('test');`)
			require.NoError(t, err)

			return nil
		})
		require.NoError(t, err)

		var count int
		err = txn(ctx).GetContext(ctx, &count, `SELECT COUNT(*) FROM transactor_transactions`)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("it should work with concurrent access", func(t *testing.T) {
		err := transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			g, gCtx := errgroup.WithContext(ctx)

			for i := 0; i < 10; i++ {
				g.Go(func() error {
					var c int
					return txn(gCtx).GetContext(gCtx, &c, `SELECT COUNT(*) FROM transactor_transactions`)
				})
			}

			return g.Wait()
		})
		require.NoError(t, err)
	})
}
