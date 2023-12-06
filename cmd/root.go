package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/k8s"
	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/metrics"
	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/web"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/spf13/cobra"

	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/logging"
	"github.com/bilalcaliskan/varnish-cache-invalidator/internal/options"
	"github.com/dimiro1/banner"
)

var (
	logger     *zap.Logger
	opts       *options.VarnishCacheInvalidatorOptions
	restConfig *rest.Config
	clientSet  *kubernetes.Clientset
	err        error
)

func init() {
	logger = logging.GetLogger()
	opts = options.GetVarnishCacheInvalidatorOptions()
	bannerBytes, _ := os.ReadFile("banner.txt")
	banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
}

func init() {
	opts = options.GetVarnishCacheInvalidatorOptions()
	rootCmd.Flags().BoolVarP(&opts.IsLocal, "isLocal", "", false,
		"IsLocal is only a flag for local development")
	rootCmd.Flags().StringVarP(&opts.KubeConfigPath, "kubeConfigPath", "",
		filepath.Join(os.Getenv("KUBECONFIG")), "KubeConfigPath is also only for local development to "+
			"specify kubeconfig while developing varnish-cache-invalidator")
	rootCmd.Flags().StringVarP(&opts.VarnishNamespace, "varnishNamespace", "",
		"default", "VarnishNamespace is the namespace of the target Varnish pods")
	rootCmd.Flags().StringVarP(&opts.VarnishLabel, "varnishLabel", "",
		"app=varnish", "VarnishLabel is the label to select proper Varnish pods")
	rootCmd.Flags().BoolVarP(&opts.InCluster, "inCluster", "", true,
		"InCluster is the boolean flag if varnish-cache-invalidator is running inside cluster or not")
	rootCmd.Flags().StringVarP(&opts.TargetHosts, "targetHosts", "",
		"http://127.0.0.1:6081", "TargetHosts is comma seperated list of target hosts, used when our "+
			"Varnish instances are not running in Kubernetes as a pod, required for standalone Varnish instances")
	rootCmd.Flags().IntVarP(&opts.ServerPort, "serverPort", "",
		3000, "ServerPort is the web server port of the varnish-cache-invalidator")
	rootCmd.Flags().IntVarP(&opts.WriteTimeoutSeconds, "writeTimeoutSeconds", "",
		10, "WriteTimeoutSeconds is the write timeout of the both web server and metrics server")
	rootCmd.Flags().IntVarP(&opts.ReadTimeoutSeconds, "readTimeoutSeconds", "",
		10, "ReadTimeoutSeconds is the read timeout of the both web server and metrics server")
	rootCmd.Flags().IntVarP(&opts.MetricsPort, "metricsPort", "",
		3001, "MetricsPort is the port of the metrics server")
	rootCmd.Flags().StringVarP(&opts.BannerFilePath, "bannerFilePath", "", "banner.txt",
		"relative path of the banner file")
}

var rootCmd = &cobra.Command{
	Use:   "varnish-cache-invalidator",
	Short: "Cache invalidation tool for Varnish opensource instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(opts.BannerFilePath); err == nil {
			bannerBytes, _ := os.ReadFile(opts.BannerFilePath)
			banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
		}

		// below check ensures that if our Varnish instances inside kubernetes or not
		if opts.InCluster {
			logger.Info("initializing kube client")
			if restConfig, err = k8s.GetConfig(); err != nil {
				logger.Error("fatal error occurred while initializing kube client", zap.String("error", err.Error()))
				return err
			}

			if clientSet, err = k8s.GetClientSet(restConfig); err != nil {
				logger.Error("fatal error occurred while getting client set", zap.String("error", err.Error()))
				return err
			}

			logger = logger.With(zap.Bool("inCluster", opts.InCluster), zap.String("masterIp", restConfig.Host),
				zap.String("varnishLabel", opts.VarnishLabel), zap.String("varnishNamespace", opts.VarnishNamespace))
			logger.Info("will use kubernetes pod instances, running pod informer to fetch pods")
			go k8s.RunPodInformer(clientSet)
		} else {
			splitted := strings.Split(opts.TargetHosts, ",")
			logger.Info("will use standalone varnish instances", zap.Any("instances", splitted))
			options.VarnishInstances = append(options.VarnishInstances, splitted...)
		}

		go func() {
			if err = metrics.RunMetricsServer(); err != nil {
				logger.Error("fatal error occured while spinning metrics server, ignoring", zap.Error(err))
			}
		}()

		if err = web.RunWebServer(); err != nil {
			logger.Error("fatal error occured while spinning web server", zap.Error(err))
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
