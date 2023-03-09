package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/logging"
	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/options"

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

// RunMetricsServer provides an endpoint, exports prometheus metrics using prometheus client golang
func RunMetricsServer() error {
	router := mux.NewRouter()
	metricServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", opts.MetricsPort),
		WriteTimeout: time.Duration(int32(opts.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(opts.ReadTimeoutSeconds)) * time.Second,
	}

	router.Handle("/metrics", promhttp.Handler())
	logger.Info("starting metrics server", zap.Int("port", opts.MetricsPort))
	return metricServer.ListenAndServe()
}
