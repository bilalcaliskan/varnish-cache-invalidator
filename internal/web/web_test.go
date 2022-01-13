package web

import (
	"testing"
	"time"
)

func TestRunWebServer(t *testing.T) {
	errChan := make(chan error, 1)

	go func() {
		errChan <- RunWebServer()
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
