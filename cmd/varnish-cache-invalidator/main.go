package main

import (
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strings"
	"varnish-cache-invalidator/internal/k8s"
	"varnish-cache-invalidator/internal/logging"
	"varnish-cache-invalidator/internal/metrics"
	"varnish-cache-invalidator/internal/options"
	"varnish-cache-invalidator/internal/web"

	"github.com/dimiro1/banner"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	opts   *options.VarnishCacheInvalidatorOptions
)

func init() {
	logger = logging.GetLogger()
	opts = options.GetVarnishCacheInvalidatorOptions()
	bannerBytes, _ := ioutil.ReadFile("banner.txt")
	banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
}

func main() {
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()

	logger.Info("initializing kube client")

	var (
		restConfig *rest.Config
		clientSet  *kubernetes.Clientset
		err        error
	)

	if restConfig, err = k8s.GetConfig(); err != nil {
		logger.Fatal("fatal error occurred while initializing kube client", zap.String("error", err.Error()))
	}

	if clientSet, err = k8s.GetClientSet(restConfig); err != nil {
		logger.Fatal("fatal error occurred while getting client set", zap.String("error", err.Error()))
	}

	logger = logger.With(zap.Bool("inCluster", opts.InCluster), zap.String("masterIp", restConfig.Host),
		zap.String("varnishLabel", opts.VarnishLabel), zap.String("varnishNamespace", opts.VarnishNamespace))

	// below check ensures that if our Varnish instances inside kubernetes or not
	if opts.InCluster {
		logger.Info("will use kubernetes pod instances, running pod informer to fetch pods")
		go k8s.RunPodInformer(clientSet)
	} else {
		splitted := strings.Split(opts.TargetHosts, ",")
		logger.Info("will use standalone varnish instances", zap.Any("instances", splitted))
		for _, v := range splitted {
			options.VarnishInstances = append(options.VarnishInstances, &v)
		}
	}

	go func() {
		if err := metrics.RunMetricsServer(); err != nil {
			logger.Fatal("fatal error occured while spinning metrics server", zap.Error(err))
		}
	}()

	if err := web.RunWebServer(); err != nil {
		logger.Fatal("fatal error occured while spinning web server", zap.Error(err))
	}
}
