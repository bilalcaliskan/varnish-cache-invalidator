package options

var vcio = &VarnishCacheInvalidatorOptions{}
var VarnishInstances []string

// GetOreillyTrialOptions returns the pointer of OreillyTrialOptions
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
	// BannerFilePath is the relative path to the banner file
	BannerFilePath string
}
