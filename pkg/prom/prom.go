package prom

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const defaultNamespace = "goboiler"

func CounterFactory(name string, help string) prometheus.Counter {
	collector := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: defaultNamespace,
		Name:      name,
		Help:      help,
	})

	if err := prometheus.Register(collector); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(prometheus.Counter)
		}

		log.Panic().Err(err).Str("name", name).Msg("failed to register prometheus collector")
	}

	return collector
}

func CounterVecFactory(name string, help string, labels []string) *prometheus.CounterVec {
	collector := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: defaultNamespace,
		Name:      name,
		Help:      help,
	}, labels)

	if err := prometheus.Register(collector); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(*prometheus.CounterVec)
		}

		log.Panic().Err(err).Str("name", name).Msg("failed to register prometheus collector")
	}
	return collector
}

func HistogramVecFactory(name string, help string, labels []string) *prometheus.HistogramVec {
	collector := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: defaultNamespace,
		Name:      name,
		Help:      help,
	}, labels)

	if err := prometheus.Register(collector); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(*prometheus.HistogramVec)
		}

		log.Panic().Err(err).Str("name", name).Msg("failed to register prometheus collector")
	}
	return collector
}

func Handler(ctx context.Context) error {
	router := gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	server := http.Server{
		Addr:              "0.0.0.0:2112",
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Ctx(ctx).Error().Err(err).Msg("failed to start prometheus server")
		}
	}()

	<-ctx.Done()
	// nolint:contextcheck
	if err := server.Shutdown(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Ctx(ctx).Error().Err(err).Msg("failed to shutdown prometheus server")
		return fmt.Errorf("failed to shutdown prometheus server: %w", err)
	}

	return nil
}
