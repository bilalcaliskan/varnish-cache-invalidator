package web

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func TestRunWebServer(t *testing.T) {
	errChan := make(chan error, 1)

	go func() {
		router := mux.NewRouter()
		webServer := initServer(router, fmt.Sprintf(":%d", opts.ServerPort),
			time.Duration(int32(opts.WriteTimeoutSeconds))*time.Second,
			time.Duration(int32(opts.ReadTimeoutSeconds))*time.Second, logger)
		logger.Info("starting web server", zap.Int("port", opts.ServerPort))
		errChan <- webServer.ListenAndServe()
	}()

	for {
		select {
		case c := <-errChan:
			t.Error(c)
		case <-time.After(10 * time.Second):
			t.Log("success")
			return
		}
	}
}
