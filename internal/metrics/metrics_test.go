package metrics

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func TestRunMetricsServer(t *testing.T) {
	errChan := make(chan error, 1)

	go func() {
		router := mux.NewRouter()
		metricServer := &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf(":%d", opts.MetricsPort),
			WriteTimeout: time.Duration(int32(opts.WriteTimeoutSeconds)) * time.Second,
			ReadTimeout:  time.Duration(int32(opts.ReadTimeoutSeconds)) * time.Second,
		}
		router.Handle("/metrics", promhttp.Handler())
		logger.Info("metric server is up and running", zap.Int("port", opts.MetricsPort))
		errChan <- metricServer.ListenAndServe()
	}()

	for {
		select {
		case c := <-errChan:
			t.Error(c)
		case <-time.After(15 * time.Second):
			t.Log("success")
			return
		}
	}
}
