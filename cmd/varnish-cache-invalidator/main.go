package main

import (
	_ "github.com/dimiro1/banner/autoload"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"
	"strings"
	"varnish-cache-invalidator/pkg/config"
	"varnish-cache-invalidator/pkg/k8s"
	"varnish-cache-invalidator/pkg/logging"
	"varnish-cache-invalidator/pkg/metrics"
	"varnish-cache-invalidator/pkg/web"
)

var (
	router                                                                 *mux.Router
	logger                                                                 *zap.Logger
	clientSet                                                              *kubernetes.Clientset
	masterUrl, kubeConfigPath, varnishNamespace, varnishLabel, targetHosts string
	inCluster                                                              bool
)

func init() {
	logger = logging.GetLogger()

	masterUrl = config.GetStringEnv("MASTER_URL", "")
	kubeConfigPath = config.GetStringEnv("KUBE_CONFIG_PATH", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	varnishNamespace = config.GetStringEnv("VARNISH_NAMESPACE", "default")
	varnishLabel = config.GetStringEnv("VARNISH_LABEL", "app=varnish")
	inCluster = config.GetBoolEnv("IN_CLUSTER", false)
	// targetHosts used when our Varnish instances are not running in Kubernetes as a pod
	// use comma separated list of instances. ex: TARGET_HOSTS=http://172.17.0.7:6081,http://172.17.0.8:6081
	targetHosts = config.GetStringEnv("TARGET_HOSTS", "")

	router = mux.NewRouter()
}

func main() {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	logger.Info("initializing kube client", zap.String("masterUrl", masterUrl),
		zap.String("kubeConfigPath", kubeConfigPath), zap.Bool("inCluster", inCluster))
	restConfig, err := k8s.GetConfig(masterUrl, kubeConfigPath, inCluster)
	if err != nil {
		logger.Fatal("fatal error occurred while initializing kube client", zap.String("error", err.Error()))
	}

	clientSet, err = k8s.GetClientSet(restConfig)
	if err != nil {
		logger.Fatal("fatal error occurred while getting client set", zap.String("error", err.Error()))
	}

	// below check ensures that we will use our Varnish instances as Kubernetes pods
	if targetHosts == "" {
		go k8s.RunPodInformer(clientSet, varnishLabel, varnishNamespace)
	} else {
		// TODO: Check the case that Varnish instances are not running inside Kubernetes. Check them after standalone installation
		splitted := strings.Split(targetHosts, ",")
		for _, v := range splitted {
			k8s.VarnishInstances = append(k8s.VarnishInstances, &v)
		}
	}

	go metrics.RunMetricsServer(router)
	web.RunWebServer(router)
}
