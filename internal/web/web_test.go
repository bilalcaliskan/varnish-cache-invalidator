package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
	"varnish-cache-invalidator/internal/options"
)

func TestPurgeHandler(t *testing.T) {
	errChan := make(chan error, 1)
	defer close(errChan)
	connChan := make(chan bool, 1)
	defer close(connChan)
	mockChan := make(chan bool, 1)
	defer close(mockChan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			if _, err := fmt.Fprint(writer, "asdasfas"); err != nil {
				t.Errorf("a fatal error occured while writing response body: %s", err.Error())
				return
			}
		}))
		options.VarnishInstances = append(options.VarnishInstances, &mockServer.URL)
		wg.Done()
		<-mockChan
		defer mockServer.Close()
	}()
	wg.Wait()

	go func() {
		errChan <- RunWebServer()
	}()

	go func() {
		for i := 0; i <= 5; i++ {
			if i == 5 {
				t.Errorf("connection to port %d could not succeeded, not retrying!\n", opts.ServerPort)
				return
			}

			_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", opts.ServerPort))
			if err != nil {
				t.Logf("connection to port %d could not succeeded, retrying...\n", opts.ServerPort)
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
			var cases = []struct {
				name, purgePath, purgeDomain string
				expectedCode                 int
			}{
				{"case1", "/", "example.com", 200},
				{"case2", "", "example.com", 400},
				{"case3", "/", "", 400},
			}

			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					client := &http.Client{Timeout: 10 * time.Second}
					req, err := http.NewRequest("PURGE", fmt.Sprintf("http://127.0.0.1:%d/purge", opts.ServerPort), nil)
					assert.Nil(t, err)
					assert.NotNil(t, req)

					req.Header.Set("purge-path", tc.purgePath)
					req.Header.Set("purge-domain", tc.purgeDomain)

					resp, err := client.Do(req)
					assert.NotNil(t, resp)
					assert.Nil(t, err)
					assert.Equal(t, tc.expectedCode, resp.StatusCode)
				})
			}
			return
		}
	}
}
