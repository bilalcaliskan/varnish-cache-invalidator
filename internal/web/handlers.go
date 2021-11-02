package web

import (
	"fmt"
	"net/http"
	"varnish-cache-invalidator/internal/options"

	"go.uber.org/zap"
)

func banHandler(w http.ResponseWriter, r *http.Request) {
	var successCount int
	var response string
	logger = logger.With(zap.String("requestMethod", "BAN"))
	banRegex := r.Header.Get("ban-regex")
	if banRegex == "" {
		logger.Error("Unable to make a request to Varnish targets, header ban-regex must be set!")
		http.Error(w, "Header ban-regex must be set!", http.StatusBadRequest)
		return
	}

	for _, v := range options.VarnishInstances {
		req, _ := http.NewRequest("BAN", *v, nil)
		req.Header.Set("ban-url", banRegex)
		logger.Info("Making BAN request", zap.String("targetHost", *v))
		res, err := client.Do(req)
		if err != nil {
			logger.Error("An error occurred while making BAN request", zap.String("targetHost", *v),
				zap.String("error", err.Error()))
		}

		if res != nil && res.StatusCode == http.StatusOK {
			successCount++
		}

	}

	if successCount == len(options.VarnishInstances) {
		logger.Info("All BAN requests succeeded on Varnish pods!", zap.Int("successCount", successCount))
		w.WriteHeader(http.StatusOK)
	} else {
		logger.Warn("One or more Varnish BAN requests failed", zap.Int("successCount", successCount),
			zap.Int("failureCount", len(options.VarnishInstances)-successCount))
		response = fmt.Sprintf("One or more Varnish BAN requests failed, check the logs!\nSucceeded request = %d\n"+
			"Failed request = %d", successCount, len(options.VarnishInstances)-successCount)
		w.WriteHeader(http.StatusBadRequest)
	}

	writeResponse(w, response)
}

func purgeHandler(w http.ResponseWriter, r *http.Request) {
	var successCount int
	var response string
	logger = logger.With(zap.String("requestMethod", "PURGE"))
	purgePath := r.Header.Get("purge-path")
	if purgePath == "" {
		logger.Error("Unable to make a PURGE request to Varnish targets, header purge-path must be set!",
			zap.String("requestMethod", "PURGE"))
		http.Error(w, "Header purge-path must be set!", http.StatusBadRequest)
		return
	}

	purgeDomain := r.Header.Get("purge-domain")
	if purgeDomain == "" {
		logger.Error("Unable to make a PURGE request to Varnish targets, header purge-domain must be set!")
		http.Error(w, "Header purge-domain must be set!", http.StatusBadRequest)
		return
	}

	for _, v := range options.VarnishInstances {
		fullUrl := fmt.Sprintf("%s%s", *v, purgePath)
		req, _ := http.NewRequest("PURGE", fullUrl, nil)
		req.Host = purgeDomain

		logger.Info("Making PURGE request", zap.String("targetHost", *v))
		res, err := client.Do(req)
		if err != nil {
			logger.Error("An error occurred while making PURGE request", zap.String("targetHost", *v),
				zap.String("error", err.Error()))
		}

		if res != nil && res.StatusCode == http.StatusOK {
			successCount++
		}
	}

	if successCount == len(options.VarnishInstances) {
		logger.Info("All PURGE requests succeeded on Varnish pods!", zap.Int("successCount", successCount))
		w.WriteHeader(http.StatusOK)
	} else {
		logger.Warn("One or more Varnish PURGE requests failed", zap.Int("successCount", successCount),
			zap.Int("failureCount", len(options.VarnishInstances)-successCount))
		response = fmt.Sprintf("One or more Varnish PURGE requests failed, check the logs!\nSucceeded request = %d\n"+
			"Failed request = %d", successCount, len(options.VarnishInstances)-successCount)
		w.WriteHeader(http.StatusBadRequest)
	}

	writeResponse(w, response)
}
