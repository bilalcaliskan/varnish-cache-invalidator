package metrics

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"testing"
	"time"
)

func TestRunMetricsServer(t *testing.T) {
	errChan := make(chan error, 1)

	go func() {
		router := mux.NewRouter()
		metricServer := &http.Server{
			Handler: router,
			Addr: fmt.Sprintf(":%d", metricsPort),
			WriteTimeout: time.Duration(int32(writeTimeoutSeconds)) * time.Second,
			ReadTimeout:  time.Duration(int32(readTimeoutSeconds)) * time.Second,
		}
		router.Handle("/metrics", promhttp.Handler())
		logger.Info("metric server is up and running", zap.Int("port", metricsPort))
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