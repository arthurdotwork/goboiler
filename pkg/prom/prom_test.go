package prom_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/arthureichelberger/goboiler/pkg/prom"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestCounterFactory(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		counter := prom.CounterFactory("metrics_count", "help")
		sameCounter := prom.CounterFactory("metrics_count", "help")
		require.NotNil(t, counter)
		require.Equal(t, counter, sameCounter)

		counter.Inc()

		require.Equal(t, 1, testutil.CollectAndCount(counter))
	})
}

func TestCounterVecFactory(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		counterVec := prom.CounterVecFactory("metrics_counter_vec", "help", []string{"label"})
		sameCounterVec := prom.CounterVecFactory("metrics_counter_vec", "help", []string{"label"})
		require.NotNil(t, counterVec)
		require.Equal(t, counterVec, sameCounterVec)

		counterVec.WithLabelValues("value").Inc()

		require.Equal(t, 1, testutil.CollectAndCount(counterVec))
		require.InEpsilon(t, 1.0, testutil.ToFloat64(counterVec.WithLabelValues("value")), 0.01)
	})
}

func TestHistogramVecFactory(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		histogramVec := prom.HistogramVecFactory("metrics_histogram", "help", []string{"label"})
		sameHistogramVec := prom.HistogramVecFactory("metrics_histogram", "help", []string{"label"})
		require.NotNil(t, histogramVec)
		require.Equal(t, histogramVec, sameHistogramVec)

		histogramVec.WithLabelValues("value").Observe(1)

		require.Equal(t, 1, testutil.CollectAndCount(histogramVec))
	})
}

func TestHandler(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return prom.Handler(ctx)
	})

	done := make(chan struct{}, 1)
	require.Eventually(t, func() bool {
		defer func() { done <- struct{}{} }()

		resp, err := http.Get("http://localhost:2112/metrics") // nolint:noctx
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == http.StatusOK
	}, time.Second, time.Millisecond*100)

	<-done
	cancel()
	require.NoError(t, errGroup.Wait())

	_, err := http.Get("http://localhost:2112/metrics") // nolint:noctx,bodyclose
	require.Error(t, err)
}
