package web

import (
	"fmt"
	"net/http"
	"varnish-cache-invalidator/internal/k8s"

	"go.uber.org/zap"
)

func banHandler(w http.ResponseWriter, r *http.Request) {
	var successCount int
	var response string
	banRegex := r.Header.Get("ban-regex")
	if banRegex == "" {
		logger.Error("Unable to make a request to Varnish targets, header ban-regex must be set!",
			zap.String("requestMethod", "BAN"))
		http.Error(w, "Header ban-regex must be set!", http.StatusBadRequest)
		return
	}

	for _, v := range k8s.VarnishInstances {
		req, _ := http.NewRequest("BAN", *v, nil)
		req.Header.Set("ban-url", banRegex)
		logger.Info("Making BAN request", zap.String("requestMethod", "BAN"), zap.String("targetHost", *v))
		res, err := client.Do(req)
		if err != nil {
			logger.Error("An error occurred while making BAN request", zap.String("requestMethod", "BAN"),
				zap.String("error", err.Error()))
		}

		if res != nil && res.StatusCode == http.StatusOK {
			successCount++
		}

	}

	if successCount == len(k8s.VarnishInstances) {
		logger.Info("All BAN requests succeeded on Varnish pods!", zap.String("requestMethod", "BAN"),
			zap.Int("successCount", successCount))
		w.WriteHeader(http.StatusOK)
	} else {
		logger.Warn("One or more Varnish BAN requests failed", zap.String("requestMethod", "BAN"),
			zap.Int("successCount", successCount), zap.Int("failureCount", len(k8s.VarnishInstances)-successCount))
		response = fmt.Sprintf("One or more Varnish BAN requests failed, check the logs!\nSucceeded request = %d\n"+
			"Failed request = %d", successCount, len(k8s.VarnishInstances)-successCount)
		w.WriteHeader(http.StatusBadRequest)
	}

	writeResponse(w, response)
}

func purgeHandler(w http.ResponseWriter, r *http.Request) {
	var successCount int
	var response string
	purgePath := r.Header.Get("purge-path")
	if purgePath == "" {
		logger.Error("Unable to make a PURGE request to Varnish targets, header purge-path must be set!",
			zap.String("requestMethod", "PURGE"))
		http.Error(w, "Header purge-path must be set!", http.StatusBadRequest)
		return
	}

	for _, v := range k8s.VarnishInstances {
		fullUrl := fmt.Sprintf("%s%s", *v, purgePath)
		req, _ := http.NewRequest("PURGE", fullUrl, nil)
		req.Host = opts.PurgeDomain

		logger.Info("Making PURGE request", zap.String("requestMethod", "PURGE"), zap.String("targetHost", *v))
		res, err := client.Do(req)
		if err != nil {
			logger.Error("An error occurred while making PURGE request", zap.String("requestMethod", "PURGE"),
				zap.String("error", err.Error()))
		}

		if res != nil && res.StatusCode == http.StatusOK {
			successCount++
		}
	}

	if successCount == len(k8s.VarnishInstances) {
		logger.Info("All PURGE requests succeeded on Varnish pods!", zap.String("requestMethod", "PURGE"),
			zap.Int("successCount", successCount))
		w.WriteHeader(http.StatusOK)
	} else {
		logger.Warn("One or more Varnish PURGE requests failed", zap.String("requestMethod", "PURGE"),
			zap.Int("successCount", successCount), zap.Int("failureCount", len(k8s.VarnishInstances)-successCount))
		response = fmt.Sprintf("One or more Varnish PURGE requests failed, check the logs!\nSucceeded request = %d\n"+
			"Failed request = %d", successCount, len(k8s.VarnishInstances)-successCount)
		w.WriteHeader(http.StatusBadRequest)
	}

	writeResponse(w, response)
}
