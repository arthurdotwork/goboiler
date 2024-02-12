package psql_test

import (
	"context"
	"testing"

	"github.com/arthureichelberger/goboiler/pkg/psql"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("it should return an error if credentials are wrong", func(t *testing.T) {
		db, err := psql.Connect(ctx, "postgres", "wrong", "localhost", "5432", "postgres")
		require.Error(t, err)
		require.Nil(t, db)
	})

	t.Run("it should connect to the database", func(t *testing.T) {
		db, err := psql.Connect(ctx, "postgres", "postgres", "localhost", "5432", "postgres")
		require.NoError(t, err)
		require.NotNil(t, db)
	})
}
