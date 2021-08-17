package options

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var vcio = &VarnishCacheInvalidatorOptions{}

func init() {
	vcio.addFlags(pflag.CommandLine)
	pflag.Parse()
}

// GetVarnishCacheInvalidatorOptions returns the pointer of VarnishCacheInvalidatorOptions
func GetVarnishCacheInvalidatorOptions() *VarnishCacheInvalidatorOptions {
	return vcio
}

// VarnishCacheInvalidatorOptions contains frequent command line and application options.
type VarnishCacheInvalidatorOptions struct {
	// MasterUrl is the url of the kube-apiserver. Not needed if kubeconfig provided
	MasterUrl string
	// KubeConfigPath is the path of the kubeconfig file to access the cluster
	KubeConfigPath string
	// VarnishNamespace is the namespace of the target Varnish pods
	VarnishNamespace string
	// VarnishLabel is the label to select proper Varnish pods
	VarnishLabel string
	// InCluster is the boolean flag if varnish-cache-invalidator is running inside cluster or not
	InCluster bool
	// TargetHosts used when our Varnish instances are not running in Kubernetes as a pod
	// use comma separated list of instances. ex: TARGET_HOSTS=http://172.17.0.7:6081,http://172.17.0.8:6081
	// mostly this is required while running outside of the cluster for debugging purposes
	TargetHosts string
}

func (vcio *VarnishCacheInvalidatorOptions) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&vcio.MasterUrl, "masterUrl", "", "MasterUrl is the url of the kube-apiserver. " +
		"Not needed if kubeconfig provided")
	fs.StringVar(&vcio.KubeConfigPath, "kubeconfigPath", filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		"KubeConfigPath is the path of the kubeconfig file to access the cluster")
	fs.StringVar(&vcio.VarnishNamespace, "varnishNamespace", "default",
		"VarnishNamespace is the namespace of the target Varnish pods")
	fs.StringVar(&vcio.VarnishLabel, "varnishLabel", "app=varnish",
		"VarnishLabel is the label to select proper Varnish pods")
	fs.BoolVar(&vcio.InCluster, "inCluster", true,
		"InCluster is the boolean flag if varnish-cache-invalidator is running inside cluster or not")
	fs.StringVar(&vcio.TargetHosts, "targetHosts", "",
		"TargetHosts used when our Varnish instances are not running in Kubernetes as a pod")
}
