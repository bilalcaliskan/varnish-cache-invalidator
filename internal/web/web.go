package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/logging"
	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/options"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	client *http.Client
	logger *zap.Logger
	opts   *options.VarnishCacheInvalidatorOptions
)

func init() {
	logger = logging.GetLogger()
	client = &http.Client{}
	opts = options.GetVarnishCacheInvalidatorOptions()
}

// RunWebServer runs the web server which multiplexes client requests
func RunWebServer() error {
	router := mux.NewRouter()
	webServer := initServer(router, fmt.Sprintf(":%d", opts.ServerPort),
		time.Duration(int32(opts.WriteTimeoutSeconds))*time.Second,
		time.Duration(int32(opts.ReadTimeoutSeconds))*time.Second, logger)

	logger.Info("starting web server", zap.Int("port", opts.ServerPort))
	return webServer.ListenAndServe()
}
