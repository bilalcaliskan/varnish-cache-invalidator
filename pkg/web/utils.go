package web

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func registerHandlers(router *mux.Router) {
	router.HandleFunc("/ban", banHandler).Methods("BAN").Schemes("http").Name("ban")
	router.HandleFunc("/purge", purgeHandler).Methods("PURGE").Schemes("http").Name("purge")
}

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

func writeResponse(w http.ResponseWriter, response string) {
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Fatalln(err)
	}
}
