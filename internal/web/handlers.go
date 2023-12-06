package web

import (
	"fmt"
	"net/http"

	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/options"

	"go.uber.org/zap"
)

func purgeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		successCount, failureCount int
		httpResponse               string
	)

	purgePath := r.Header.Get("purge-path")
	if purgePath == "" {
		logger.Error("unable to make a PURGE request to Varnish targets, header purge-path must be set!")
		http.Error(w, "Header purge-path must be set!", http.StatusBadRequest)
		return
	}

	purgeDomain := r.Header.Get("purge-domain")
	if purgeDomain == "" {
		logger.Error("unable to make a PURGE request to Varnish targets, header purge-domain must be set!")
		http.Error(w, "Header purge-domain must be set!", http.StatusBadRequest)
		return
	}

	for _, v := range options.VarnishInstances {
		fullUrl := fmt.Sprintf("%s%s", v, purgePath)
		req, _ := http.NewRequest("PURGE", fullUrl, nil)
		// fullUrl := fmt.Sprintf("http://192.168.49.2:30654%s", purgePath)
		// req.Host = "nginx.default.svc"
		req.Host = purgeDomain

		logger.Info("making PURGE request", zap.String("url", fullUrl))
		res, err := client.Do(req)
		if err != nil {
			logger.Error("an error occurred while making PURGE request", zap.String("url", fullUrl),
				zap.String("error", err.Error()))
			failureCount++
		}

		if res.StatusCode == http.StatusOK {
			successCount++
		}
	}

	if successCount == len(options.VarnishInstances) {
		logger.Info("all PURGE requests succeeded on Varnish pods!", zap.Int("successCount", successCount),
			zap.Int("failureCount", failureCount))
		httpResponse = fmt.Sprintf("All PURGE requests succeeded on Varnish pods!\nSucceeded request = %d\n"+
			"Failed request = %d\n", successCount, failureCount)
		w.WriteHeader(http.StatusOK)
	} else {
		logger.Warn("one or more Varnish PURGE requests failed", zap.Int("successCount", successCount),
			zap.Int("failureCount", len(options.VarnishInstances)-successCount))
		httpResponse = fmt.Sprintf("One or more Varnish PURGE requests failed, check the logs!\nSucceeded request = %d\n"+
			"Failed request = %d\n", successCount, failureCount)
		w.WriteHeader(http.StatusInternalServerError)
	}

	writeResponse(w, httpResponse)
}
