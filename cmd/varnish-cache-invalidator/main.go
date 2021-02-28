package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"varnish-cache-invalidator/pkg/config"
	http2 "varnish-cache-invalidator/pkg/http"
	"varnish-cache-invalidator/pkg/k8s"
)

var (
	clientSet *kubernetes.Clientset
	serverPort, writeTimeoutSeconds, readTimeoutSeconds int
	masterUrl, kubeConfigPath, varnishNamespace, varnishLabel, targetHosts string
	inCluster bool
)

func init() {
	serverPort = config.GetIntEnv("SERVER_PORT", 3000)
	writeTimeoutSeconds = config.GetIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	readTimeoutSeconds = config.GetIntEnv("READ_TIMEOUT_SECONDS", 10)
	masterUrl = config.GetStringEnv("MASTER_URL", "")
	kubeConfigPath = config.GetStringEnv("KUBE_CONFIG_PATH", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	varnishNamespace = config.GetStringEnv("VARNISH_NAMESPACE", "default")
	varnishLabel = config.GetStringEnv("VARNISH_LABEL", "app=varnish")
	inCluster = config.GetBoolEnv("IN_CLUSTER", false)
	// targetHosts used when our Varnish instances are not running in Kubernetes as a pod
	// use comma seperated list of instances. ex: TARGET_HOSTS=http://172.17.0.7:6081,http://172.17.0.8:6081
	targetHosts = config.GetStringEnv("TARGET_HOSTS", "")
}

func main() {
	log.Println("Initializing kube client...")
	restConfig, err := k8s.GetConfig(masterUrl, kubeConfigPath, inCluster)
	if err != nil {
		log.Fatalln(err)
	}

	clientSet, err = k8s.GetClientSet(restConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// below check ensures that we will use our Varnish instances as Kubernetes pods
	if targetHosts == "" {
		k8s.RunPodInformer(clientSet, varnishLabel, varnishNamespace)
	} else {
		// TODO: Check the case that Varnish instances are not running inside Kubernetes. Check them after standalone installation
		splitted := strings.Split(targetHosts, ",")
		for _, v := range splitted {
			k8s.VarnishInstances = append(k8s.VarnishInstances, &v)
		}
	}

	router := mux.NewRouter()
	server := http2.InitServer(router, fmt.Sprintf(":%d", serverPort), time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		time.Duration(int32(readTimeoutSeconds)) * time.Second)
	log.Printf("Server is listening on port %d!", serverPort)
	log.Fatal(server.ListenAndServe())
}