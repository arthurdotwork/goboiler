package store

import (
	"context"
	"fmt"

	"github.com/arthureichelberger/goboiler/pkg/psql"
)

type DummyStore interface {
	Dummy(ctx context.Context) (int64, error)
}

type dummyStore struct {
	db psql.DBGetter
}

func NewDummyStore(db psql.DBGetter) DummyStore {
	return dummyStore{
		db: db,
	}
}

func (s dummyStore) Dummy(ctx context.Context) (int64, error) {
	query := `SELECT 1+1 as result`

	var count int64
	if err := s.db(ctx).GetContext(ctx, &count, query); err != nil {
		return 0, fmt.Errorf("failed to dummy: %w", err)
	}

	return count, nil
}
