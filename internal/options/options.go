package options

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var vcio = &VarnishCacheInvalidatorOptions{}

// VarnishInstances keeps pointer of varnish instances' ip:port information
var VarnishInstances []*string

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
	// IsLocal is only for local development
	IsLocal bool
	// KubeConfigPath is also only for local development to specify kubeconfig
	KubeConfigPath string
	// InCluster is the boolean flag if varnish-cache-invalidator is running inside cluster or not
	InCluster bool
	// VarnishNamespace is the namespace of the target Varnish pods
	VarnishNamespace string
	// VarnishLabel is the label to select proper Varnish pods
	VarnishLabel string
	// TargetHosts used when our Varnish instances are not running in Kubernetes as a pod
	// use comma separated list of instances. ex: TARGET_HOSTS=http://172.17.0.7:6081,http://172.17.0.8:6081
	// mostly this is required while running outside of the cluster for debugging purposes
	TargetHosts string
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
	fs.BoolVar(&vcio.IsLocal, "isLocal", false, "IsLocal is only a flag for local development")
	fs.StringVar(&vcio.KubeConfigPath, "kubeConfigPath", filepath.Join(os.Getenv("HOME"), "mock", "kubeconfig"),
		"KubeConfigPath is also only for local development to specify kubeconfig while developing varnish-cache-invalidator")
	fs.StringVar(&vcio.VarnishNamespace, "varnishNamespace", "default",
		"VarnishNamespace is the namespace of the target Varnish pods, defaults to default namespace")
	fs.StringVar(&vcio.VarnishLabel, "varnishLabel", "app=varnish",
		"VarnishLabel is the label to select proper Varnish pods, defaults to app=varnish")
	fs.BoolVar(&vcio.InCluster, "inCluster", true,
		"InCluster is the boolean flag if varnish-cache-invalidator is running inside cluster or not, defaults to true")
	fs.StringVar(&vcio.TargetHosts, "targetHosts", "http://127.0.0.1:6081",
		"TargetHosts is comma seperated list of target hosts, used when our Varnish instances are not running "+
			"in Kubernetes as a pod, required for standalone Varnish instances, defaults to 'http://127.0.0.1:6081'")
	fs.IntVar(&vcio.ServerPort, "serverPort", 3000, "ServerPort is the web server port of the "+
		"varnish-cache-invalidator, defaults to 3000")
	fs.IntVar(&vcio.WriteTimeoutSeconds, "writeTimeoutSeconds", 10, "WriteTimeoutSeconds is the write "+
		"timeout of the both web server and metrics server, defaults to 10")
	fs.IntVar(&vcio.ReadTimeoutSeconds, "readTimeoutSeconds", 10, "ReadTimeoutSeconds is the read "+
		"timeout of the both web server and metrics server, defaults to 10")
	fs.IntVar(&vcio.MetricsPort, "metricsPort", 3001, "MetricsPort is the port of the metrics "+
		"server, defaults to 3001")
}
