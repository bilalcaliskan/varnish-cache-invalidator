package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
	"varnish-cache-invalidator/pkg/config"
	"varnish-cache-invalidator/pkg/logging"
)

var (
	client                                              *http.Client
	purgeDomain                                         string
	logger                                              *zap.Logger
	serverPort, writeTimeoutSeconds, readTimeoutSeconds int
)

func init() {
	logger = logging.GetLogger()
	client = &http.Client{}
	// purgeDomain will set Host header on purge requests. It must be changed to work properly on different environments.
	// A purge request hit the Varnish must match the host of the cache object.
	purgeDomain = config.GetStringEnv("PURGE_DOMAIN", "foo.example.com")
	serverPort = config.GetIntEnv("SERVER_PORT", 3000)
	writeTimeoutSeconds = config.GetIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	readTimeoutSeconds = config.GetIntEnv("READ_TIMEOUT_SECONDS", 10)
}

// RunWebServer runs the web server which multiplexes client requests
func RunWebServer(router *mux.Router) {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	webServer := initServer(router, fmt.Sprintf(":%d", serverPort), time.Duration(int32(writeTimeoutSeconds))*time.Second,
		time.Duration(int32(readTimeoutSeconds))*time.Second, logger)

	logger.Info("web server is up and running", zap.Int("port", serverPort))
	panic(webServer.ListenAndServe())
}
