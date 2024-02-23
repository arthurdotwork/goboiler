package service

import (
	"context"
	"fmt"

	"github.com/arthureichelberger/goboiler/internal/store"
	"github.com/arthureichelberger/goboiler/pkg/psql"
	"github.com/rs/zerolog/log"
)

type DummyService struct {
	dummyStore store.DummyStore
	transactor psql.Transactor
}

func NewDummyService(dummyStore store.DummyStore, transactor psql.Transactor) DummyService {
	return DummyService{
		dummyStore: dummyStore,
		transactor: transactor,
	}
}

func (s DummyService) Dummy(ctx context.Context) (int64, error) {
	var count int64

	op := func(ctx context.Context) error {
		var err error
		count, err = s.dummyStore.Dummy(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to dummy")
			return fmt.Errorf("failed to dummy: %w", err)
		}

		return nil
	}

	if err := s.transactor.WithinTransaction(ctx, op); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to dummy")
		return 0, fmt.Errorf("failed to dummy: %w", err)
	}

	return count, nil
}
