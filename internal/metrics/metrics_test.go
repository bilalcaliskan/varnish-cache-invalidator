package metrics

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunMetricsServer(t *testing.T) {
	errChan := make(chan error, 1)
	connChan := make(chan bool, 1)

	go func() {
		errChan <- RunMetricsServer()
	}()

	go func() {
		for i := 0; i <= 5; i++ {
			if i == 5 {
				t.Errorf("connection to port %d could not succeeded, not retrying!\n", opts.MetricsPort)
				return
			}

			_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", opts.MetricsPort))
			if err != nil {
				t.Logf("connection to port %d could not succeeded, retrying...\n", opts.MetricsPort)
				time.Sleep(1 * time.Second)
				continue
			}

			connChan <- true
			return
		}
	}()

	for {
		select {
		case c := <-errChan:
			t.Fatal(c)
		case <-connChan:
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics", opts.MetricsPort))
			assert.Nil(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.NotEmpty(t, string(body))
			return
		case <-time.After(20 * time.Second):
			t.Fatal("could not completed in 20 seconds, failing")
		}
	}
}
