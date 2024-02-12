package service

import (
	"context"
	"fmt"

	"github.com/arthureichelberger/goboiler/internal/store"
	"github.com/rs/zerolog/log"
)

type DummyService struct {
	dummyStore store.DummyStore
}

func NewDummyService(dummyStore store.DummyStore) DummyService {
	return DummyService{
		dummyStore: dummyStore,
	}
}

func (s DummyService) Dummy(ctx context.Context) (int64, error) {
	count, err := s.dummyStore.Dummy(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to dummy")
		return 0, fmt.Errorf("failed to dummy: %w", err)
	}

	return count, nil
}
