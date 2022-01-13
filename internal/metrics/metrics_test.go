package metrics

import (
	"testing"
	"time"
)

func TestRunMetricsServer(t *testing.T) {
	errChan := make(chan error, 1)

	go func() {
		errChan <- RunMetricsServer()
	}()

	for {
		select {
		case c := <-errChan:
			t.Error(c)
		case <-time.After(5 * time.Second):
			t.Log("success")
			return
		}
	}
}
