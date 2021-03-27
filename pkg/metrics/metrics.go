package metrics

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"time"
	"varnish-cache-invalidator/pkg/config"
	"varnish-cache-invalidator/pkg/logging"
)

var (
	logger *zap.Logger
	metricsPort, writeTimeoutSeconds, readTimeoutSeconds int
)

func init() {
	logger = logging.GetLogger()
	metricsPort = config.GetIntEnv("METRICS_PORT", 3001)
	writeTimeoutSeconds = config.GetIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	readTimeoutSeconds = config.GetIntEnv("READ_TIMEOUT_SECONDS", 10)
}

// TODO: Generate custom metrics, check below:
// https://prometheus.io/docs/guides/go-application/
// https://www.robustperception.io/prometheus-middleware-for-gorilla-mux

func RunMetricsServer(router *mux.Router) {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	metricServer := &http.Server{
		Handler: router,
		Addr: fmt.Sprintf(":%d", metricsPort),
		WriteTimeout: time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(readTimeoutSeconds)) * time.Second,
	}
	router.Handle("/metrics", promhttp.Handler())
	logger.Info("metric server is up and running", zap.Int("port", metricsPort))
	panic(metricServer.ListenAndServe())
}