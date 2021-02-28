package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

func GetConfig(masterUrl, kubeConfigPath string, inCluster bool) (*rest.Config, error) {
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

func GetClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}