package metrics

import (
	"fmt"
	"net/http"
	"time"
	"varnish-cache-invalidator/internal/logging"
	"varnish-cache-invalidator/internal/options"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	opts   *options.VarnishCacheInvalidatorOptions
)

func init() {
	logger = logging.GetLogger()
	opts = options.GetVarnishCacheInvalidatorOptions()
}

// TODO: Generate custom metrics, check below:
// https://prometheus.io/docs/guides/go-application/
// https://www.robustperception.io/prometheus-middleware-for-gorilla-mux

// RunMetricsServer provides an endpoint, exports prometheus metrics using prometheus client golang
func RunMetricsServer(router *mux.Router) {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	metricServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", opts.MetricsPort),
		WriteTimeout: time.Duration(int32(opts.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(opts.ReadTimeoutSeconds)) * time.Second,
	}
	router.Handle("/metrics", promhttp.Handler())
	logger.Info("starting metrics server", zap.Int("port", opts.MetricsPort))
	panic(metricServer.ListenAndServe())
}
