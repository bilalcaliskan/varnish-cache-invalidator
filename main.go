package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	varnishInstances []*string
	clientSet *kubernetes.Clientset
	client *http.Client

	serverPort, writeTimeoutSeconds, readTimeoutSeconds int
	masterUrl, kubeConfigPath, varnishNamespace, varnishLabelKey, varnishLabelValue, purgeDomain, targetHosts string
	inCluster bool
)

func init() {
	serverPort = getIntEnv("SERVER_PORT", 3000)
	writeTimeoutSeconds = getIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	readTimeoutSeconds = getIntEnv("READ_TIMEOUT_SECONDS", 10)
	masterUrl = getStringEnv("MASTER_URL", "")
	kubeConfigPath = getStringEnv("KUBE_CONFIG_PATH", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	varnishNamespace = getStringEnv("VARNISH_NAMESPACE", "default")
	varnishLabelKey = getStringEnv("VARNISH_LABEL_KEY", "app")
	varnishLabelValue = getStringEnv("VARNISH_LABEL_VALUE", "varnish")
	// purgeDomain will set Host header on purge requests. It must be changed to work properly on different environments.
	// A purge request hit the Varnish must match the host of the cache object.
	purgeDomain = getStringEnv("PURGE_DOMAIN", "foo.example.com")
	inCluster = getBoolEnv("IN_CLUSTER", false)
	// targetHosts used when our Varnish instances are not running in Kubernetes as a pod
	// use comma seperated list of instances. ex: TARGET_HOSTS=http://172.17.0.7:6081,http://172.17.0.8:6081
	targetHosts = getStringEnv("TARGET_HOSTS", "")
}

func main() {
	log.Println("Initializing http client...")
	client = &http.Client{}

	log.Println("Initializing kube client...")
	restConfig, err := getConfig(masterUrl, kubeConfigPath, inCluster)
	checkError(err)
	clientSet, err = getClientSet(restConfig)
	checkError(err)

	// below check ensures that we will use our Varnish instances as Kubernetes pods
	if targetHosts == "" {
		runPodInformer(clientSet, varnishLabelKey, varnishLabelValue, varnishNamespace)
	} else {
		// TODO: Check the case that Varnish instances are not running inside Kubernetes. Check them after standalone installation
		splitted := strings.Split(targetHosts, ",")
		for _, v := range splitted {
			varnishInstances = append(varnishInstances, &v)
		}
	}

	router := mux.NewRouter()
	server := initServer(router, fmt.Sprintf(":%d", serverPort), time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		time.Duration(int32(readTimeoutSeconds)) * time.Second)
	log.Printf("Server is listening on port %d!", serverPort)
	log.Fatal(server.ListenAndServe())
}