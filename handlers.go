package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func banHandler(w http.ResponseWriter, r *http.Request) {
	var successCount int
	var response string
	banRegex := r.Header.Get("ban-regex")
	if banRegex == "" {
		log.Println("Unable to make a request to Varnish targets, header ban-regex must be set!")
		http.Error(w, "Header ban-regex must be set!", http.StatusBadRequest)
		return
	}

	for _, v := range varnishInstances {
		req, _ := http.NewRequest("BAN", *v, nil)
		req.Header.Set("ban-url", banRegex)
		log.Printf(  "Making BAN request on host %s\n", *v)
		res, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
		}

		if res != nil && res.StatusCode == http.StatusOK {
			successCount++
		}

	}

	if successCount == len(varnishInstances) {
		response = "All BAN requests succeeded on Varnish pods!"
		w.WriteHeader(http.StatusOK)
	} else {
		response = fmt.Sprintf("One or more Varnish BAN requests failed, check the logs!\nSucceeded request = %d\n" +
			"Failed request = %d", successCount, len(varnishInstances) - successCount)
		w.WriteHeader(http.StatusBadRequest)
	}

	writeResponse(w, response)
}

func purgeHandler(w http.ResponseWriter, r *http.Request) {
	var successCount int
	var response string
	purgePath := r.Header.Get("purge-path")
	if purgePath == "" {
		log.Println("Unable to make a request to Varnish targets, header purge-path must be set!")
		http.Error(w, "Header purge-path must be set!", http.StatusBadRequest)
		return
	}

	for _, v := range varnishInstances {
		fullUrl := fmt.Sprintf("%s%s", *v, purgePath)
		req, err := http.NewRequest("PURGE", fullUrl, nil)
		if err != nil {
			log.Println(err.Error())
		}

		if req != nil {
			req.Host = purgeDomain
		}

		log.Printf("Making PURGE request on host %s\n", fullUrl)
		res, err := client.Do(req)
		if err != nil {
			log.Printf("An error occured while making PURGE request to %s!\n%v\n", fullUrl, err.Error())
			log.Println(err.Error())
		}

		if res != nil && res.StatusCode == http.StatusOK {
			successCount++
		}
	}

	if successCount == len(varnishInstances) {
		response = "All PURGE requests succeeded on Varnish pods!"
		w.WriteHeader(http.StatusOK)
	} else {
		response = fmt.Sprintf("One or more Varnish PURGE requests failed, check the logs!\nSucceeded request = %d\n" +
			"Failed request = %d", successCount, len(varnishInstances) - successCount)
		w.WriteHeader(http.StatusBadRequest)
	}

	writeResponse(w, response)
}

func registerHandlers(router *mux.Router) {
	router.HandleFunc("/ban", banHandler).Methods("BAN").Schemes("http").Name("ban")
	router.HandleFunc("/purge", purgeHandler).Methods("PURGE").Schemes("http").Name("purge")
}