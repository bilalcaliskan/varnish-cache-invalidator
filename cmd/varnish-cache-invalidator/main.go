package main

import (
	"io/ioutil"
	"os"
	"strings"
	"varnish-cache-invalidator/internal/k8s"
	"varnish-cache-invalidator/internal/logging"
	"varnish-cache-invalidator/internal/metrics"
	"varnish-cache-invalidator/internal/options"
	"varnish-cache-invalidator/internal/web"

	"github.com/dimiro1/banner"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	router *mux.Router
	logger *zap.Logger
	opts   *options.VarnishCacheInvalidatorOptions
)

func init() {
	logger = logging.GetLogger()
	opts = options.GetVarnishCacheInvalidatorOptions()
	router = mux.NewRouter()
	bannerBytes, _ := ioutil.ReadFile("banner.txt")
	banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
}

func main() {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	// below check ensures that if our Varnish instances inside kubernetes or not
	if opts.InCluster {
		go k8s.RunPodInformer()
	} else {
		splitted := strings.Split(opts.TargetHosts, ",")
		for _, v := range splitted {
			k8s.VarnishInstances = append(k8s.VarnishInstances, &v)
		}
	}

	go metrics.RunMetricsServer(router)
	web.RunWebServer(router)
}
