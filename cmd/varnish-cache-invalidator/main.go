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
	vcio   *options.VarnishCacheInvalidatorOptions
)

func init() {
	logger = logging.GetLogger()
	vcio = options.GetVarnishCacheInvalidatorOptions()
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

	// below check ensures that we will use our Varnish instances as Kubernetes pods
	if vcio.TargetHosts == "" {
		go k8s.RunPodInformer()
	} else {
		// TODO: Check the case that Varnish instances are not running inside Kubernetes. Check them after standalone installation
		splitted := strings.Split(vcio.TargetHosts, ",")
		for _, v := range splitted {
			k8s.VarnishInstances = append(k8s.VarnishInstances, &v)
		}
	}

	go metrics.RunMetricsServer(router)
	web.RunWebServer(router)
}
