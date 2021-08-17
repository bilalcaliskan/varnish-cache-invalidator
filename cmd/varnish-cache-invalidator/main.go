package main

import (
	"github.com/dimiro1/banner"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
	"varnish-cache-invalidator/pkg/k8s"
	"varnish-cache-invalidator/pkg/logging"
	"varnish-cache-invalidator/pkg/metrics"
	"varnish-cache-invalidator/pkg/options"
	"varnish-cache-invalidator/pkg/web"
)

var (
	router    *mux.Router
	logger    *zap.Logger
	clientSet *kubernetes.Clientset
	vcio      *options.VarnishCacheInvalidatorOptions
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

	logger.Info("initializing kube client", zap.String("masterUrl", vcio.MasterUrl),
		zap.String("kubeConfigPath", vcio.KubeConfigPath), zap.Bool("inCluster", vcio.InCluster))
	restConfig, err := k8s.GetConfig(vcio.MasterUrl, vcio.KubeConfigPath, vcio.InCluster)
	if err != nil {
		logger.Fatal("fatal error occurred while initializing kube client", zap.String("error", err.Error()))
	}

	clientSet, err = k8s.GetClientSet(restConfig)
	if err != nil {
		logger.Fatal("fatal error occurred while getting client set", zap.String("error", err.Error()))
	}

	// below check ensures that we will use our Varnish instances as Kubernetes pods
	if vcio.TargetHosts == "" {
		go k8s.RunPodInformer(clientSet, vcio.VarnishLabel, vcio.VarnishNamespace)
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
