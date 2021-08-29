package options

import (
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
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
	// PurgeDomain will set Host header on purge requests. It must be changed to work properly on different environments.
	PurgeDomain string
	// ServerPort is the web server port of the varnish-cache-invalidator, defaults to 3000
	ServerPort int
	// WriteTimeoutSeconds is the write timeout of the both web server and metrics server, defaults to 10
	WriteTimeoutSeconds int
	// ReadTimeoutSeconds is the read timeout of the both web server and metrics server, defaults to 10
	ReadTimeoutSeconds int
	// MetricsPort is the port of the metrics server, defaults to 3001
	MetricsPort int
}

func (vcio *VarnishCacheInvalidatorOptions) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&vcio.MasterUrl, "masterUrl", "", "MasterUrl is the additional url to kube-apiserver. "+
		"There is no need to specify anything here since we are using serviceAccount to access kube-apiserver. "+
		"Defaults to ''")
	fs.StringVar(&vcio.KubeConfigPath, "kubeconfigPath", filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		"KubeConfigPath is the path of the kubeconfig file to access the cluster. Defaults to ~/.kube/config")
	fs.StringVar(&vcio.VarnishNamespace, "varnishNamespace", "default",
		"VarnishNamespace is the namespace of the target Varnish pods, defaults to default namespace")
	fs.StringVar(&vcio.VarnishLabel, "varnishLabel", "app=varnish",
		"VarnishLabel is the label to select proper Varnish pods, defaults to app=varnish")
	fs.BoolVar(&vcio.InCluster, "inCluster", true,
		"InCluster is the boolean flag if varnish-cache-invalidator is running inside cluster or not, defaults to true")
	fs.StringVar(&vcio.TargetHosts, "targetHosts", "",
		"TargetHosts used when our Varnish instances are not running in Kubernetes as a pod, required for standalone "+
			"Varnish instances, defaults to ''")
	fs.StringVar(&vcio.PurgeDomain, "purgeDomain", "foo.example.com", "PurgeDomain will set Host header "+
		"on purge requests. It must be changed to work properly on different environments.")
	fs.IntVar(&vcio.ServerPort, "serverPort", 3000, "ServerPort is the web server port of the "+
		"varnish-cache-invalidator, defaults to 3000")
	fs.IntVar(&vcio.WriteTimeoutSeconds, "writeTimeoutSeconds", 10, "WriteTimeoutSeconds is the write "+
		"timeout of the both web server and metrics server, defaults to 10")
	fs.IntVar(&vcio.ReadTimeoutSeconds, "readTimeoutSeconds", 10, "ReadTimeoutSeconds is the read "+
		"timeout of the both web server and metrics server, defaults to 10")
	fs.IntVar(&vcio.MetricsPort, "metricsPort", 3001, "MetricsPort is the port of the metrics "+
		"server, defaults to 3001")
}
