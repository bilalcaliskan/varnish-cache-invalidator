package web

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// registerHandlers registers the handlers of the web server
func registerHandlers(router *mux.Router) {
	router.HandleFunc("/purge", purgeHandler).Methods("PURGE").Schemes("http").Name("purge")
}

// initServer initializes the web server
func initServer(router *mux.Router, addr string, writeTimeout time.Duration, readTimeout time.Duration,
	lgr *zap.Logger) *http.Server {
	logger = lgr
	registerHandlers(router)
	return &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
}

// writeResponse helper function writes a http response body as string to http.ResponseWriter
func writeResponse(w http.ResponseWriter, response string) {
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Fatalln(err)
	}
}
