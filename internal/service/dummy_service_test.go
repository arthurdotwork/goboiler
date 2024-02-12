package service_test

import (
	"context"
	"testing"

	"github.com/arthureichelberger/goboiler/internal/service"
	"github.com/arthureichelberger/goboiler/internal/store"
	"github.com/arthureichelberger/goboiler/pkg/psql"
	"github.com/arthureichelberger/goboiler/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestDummy(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
	require.NoError(t, err)

	txn, rollback := test.Txn(t, ctx, db)
	t.Cleanup(rollback)

	dummyStore := store.NewDummyStore(txn)
	dummyService := service.NewDummyService(dummyStore)

	t.Run("it should return the correct value", func(t *testing.T) {
		count, err := dummyService.Dummy(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(2), count)
	})
}
