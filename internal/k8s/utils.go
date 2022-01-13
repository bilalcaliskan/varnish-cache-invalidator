package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// addVarnishPod add a pod string to varnishPods *[]*string
func addVarnishPod(varnishPods *[]*string, pod *string) {
	_, found := findVarnishPod(*varnishPods, *pod)
	if !found {
		*varnishPods = append(*varnishPods, pod)
	}
}

// findVarnishPod finds the specific pod string in varnishPods []*string
func findVarnishPod(varnishPods []*string, pod string) (int, bool) {
	for i, item := range varnishPods {
		if pod == *item {
			return i, true
		}
	}
	return -1, false
}

// removeVarnishPod removes the Varnish pod string from varnishPods *[]*string
func removeVarnishPod(varnishPods *[]*string, index int) {
	*varnishPods = append((*varnishPods)[:index], (*varnishPods)[index+1:]...)
}

// GetConfig initializes and returns the *rest.Config
func GetConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error

	if opts.IsLocal {
		if config, err = clientcmd.BuildConfigFromFlags("", opts.KubeConfigPath); err != nil {
			return nil, err
		}
	} else {
		if config, err = rest.InClusterConfig(); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// GetClientSet initializes and returns *kubernetes.Clientset using *rest.Config
func GetClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
