package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arthureichelberger/goboiler/internal/handler"
	"github.com/arthureichelberger/goboiler/internal/middleware"
	"github.com/arthureichelberger/goboiler/internal/service"
	"github.com/arthureichelberger/goboiler/internal/store"
	"github.com/arthureichelberger/goboiler/pkg/prom"
	"github.com/arthureichelberger/goboiler/pkg/psql"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		cancel()
	}()

	if err := run(ctx); err != nil {
		log.Fatal().Err(err)
	}
}

func run(ctx context.Context) error {
	// packages
	db, err := psql.Connect(
		ctx,
		env("DATABASE_USERNAME", "postgres"),
		env("DATABASE_PASSWORD", "postgres"),
		env("DATABASE_HOST", "localhost"),
		env("DATABASE_PORT", "5432"),
		env("DATABASE_NAME", "postgres"),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// stores
	dummyStore := store.NewDummyStore(db)

	// services
	dummyService := service.NewDummyService(dummyStore)

	router := gin.New()
	router.Use(middleware.InstrumentedMiddleware())
	router.GET("/ping", handler.PingHandler())
	router.GET("/dummy", handler.DummyHandler(dummyService))

	httpServer := &http.Server{
		Addr:              env("HTTP_ADDR", "0.0.0.0:8080"),
		Handler:           router,
		ReadHeaderTimeout: time.Second * 2,
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to listen and serve: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		<-ctx.Done()
		log.Debug().Msg("shutting down application")
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := httpServer.Shutdown(timeout); err != nil {
			return fmt.Errorf("failed to shutdown http server: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		if err := prom.Handler(ctx); err != nil {
			return fmt.Errorf("failed to run prometheus handler: %w", err)
		}

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	log.Debug().Msg("application is shutting down")
	return nil
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
