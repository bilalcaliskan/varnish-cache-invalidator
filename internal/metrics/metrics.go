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
	vcio   *options.VarnishCacheInvalidatorOptions
)

func init() {
	logger = logging.GetLogger()
	vcio = options.GetVarnishCacheInvalidatorOptions()
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
		Addr:         fmt.Sprintf(":%d", vcio.MetricsPort),
		WriteTimeout: time.Duration(int32(vcio.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(vcio.ReadTimeoutSeconds)) * time.Second,
	}
	router.Handle("/metrics", promhttp.Handler())
	logger.Info("metric server is up and running", zap.Int("port", vcio.MetricsPort))
	panic(metricServer.ListenAndServe())
}
