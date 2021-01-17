package main

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func addVarnishPod(varnishPods *[]*string, pod *string) {
	_, found := findVarnishPod(*varnishPods, *pod)
	if !found {
		*varnishPods = append(*varnishPods, pod)
	}
}

func findVarnishPod(varnishPods []*string, pod string) (int, bool) {
	for i, item := range varnishPods {
		if pod == *item {
			return i, true
		}
	}
	return -1, false
}

func removeVarnishPod(varnishPods *[]*string, index int) {
	*varnishPods = append((*varnishPods)[:index], (*varnishPods)[index+1:]...)
}

func getConfig(masterUrl, kubeConfigPath string, inCluster bool) (*rest.Config, error) {
	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags(masterUrl, kubeConfigPath)
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}

func getClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func initServer(router *mux.Router, addr string, writeTimeout time.Duration, readTimeout time.Duration) *http.Server {
	registerHandlers(router)
	return &http.Server{
		Handler: router,
		Addr: addr,
		WriteTimeout: writeTimeout,
		ReadTimeout: readTimeout,
	}
}

func writeResponse(w http.ResponseWriter, response string) {
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Fatalln(err)
	}
}

func convertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("An error occured while converting %s to int. Setting it as zero.", s)
		i = 0
	}
	return i
}

func convertStringToBool(s string) bool {
	i, err := strconv.ParseBool(s)
	if err != nil {
		log.Printf("An error occured while converting %s to int. Setting it as zero.", s)
		i = false
	}
	return i
}

func getStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return convertStringToInt(value)
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return convertStringToBool(value)
}