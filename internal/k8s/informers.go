package k8s

import (
	"fmt"
	"strings"
	"time"
	"varnish-cache-invalidator/internal/logging"
	"varnish-cache-invalidator/internal/options"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const PodUrl = "http://%s:%d"

var (
	logger *zap.Logger
	opts   *options.VarnishCacheInvalidatorOptions
)

func init() {
	opts = options.GetVarnishCacheInvalidatorOptions()
	logger = logging.GetLogger()
}

// RunPodInformer continuously watches kube-apiserver with shared informer for Pod resources, then does necessary updates
// on VarnishInstances slice on Add/Update/Delete conditions
func RunPodInformer(clientSet kubernetes.Interface) {
	varnishLabelKey := strings.Split(opts.VarnishLabel, "=")[0]
	varnishLabelValue := strings.Split(opts.VarnishLabel, "=")[1]

	informerFactory := informers.NewSharedInformerFactory(clientSet, 30*time.Second)
	podInformer := informerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labels := pod.GetLabels()
			for key, value := range labels {
				if key == varnishLabelKey && value == varnishLabelValue && pod.Namespace == opts.VarnishNamespace {
					if pod.Status.PodIP != "" {
						podUrl := fmt.Sprintf(PodUrl, pod.Status.PodIP, pod.Spec.Containers[0].Ports[0].ContainerPort)
						logger.Info("adding pod url to the varnishPods slice", zap.String("pod", pod.Name),
							zap.String("namespace", pod.Namespace), zap.String("podUrl", podUrl))
						addVarnishPod(&options.VarnishInstances, &podUrl)
					} else {
						logger.Warn("varnish pod does not have an ip address yet, skipping add operation",
							zap.String("pod", pod.Name), zap.String("namespace", pod.Namespace))
					}
				}
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*v1.Pod)
			newPod := newObj.(*v1.Pod)
			labels := oldPod.GetLabels()

			for key, value := range labels {
				if key == varnishLabelKey && value == varnishLabelValue && oldPod.ResourceVersion != newPod.ResourceVersion &&
					oldPod.Namespace == opts.VarnishNamespace {
					if oldPod.Status.PodIP == "" && newPod.Status.PodIP != "" {
						logger.Info("assigned an ip address to the pod, adding to varnishPods slice", zap.String("pod", newPod.Name),
							zap.String("namespace", newPod.Namespace), zap.String("ipAddress", newPod.Status.PodIP))
						podUrl := fmt.Sprintf(PodUrl, newPod.Status.PodIP, newPod.Spec.Containers[0].Ports[0].ContainerPort)
						logger.Info("adding pod url to the varnishPods slice", zap.String("podUrl", podUrl))
						addVarnishPod(&options.VarnishInstances, &podUrl)
					}
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labels := pod.GetLabels()
			for key, value := range labels {
				if key == varnishLabelKey && value == varnishLabelValue && pod.Namespace == opts.VarnishNamespace {
					logger.Info("varnish pod is deleted, removing from varnishPods slice", zap.String("pod", pod.Name),
						zap.String("namespace", pod.Namespace))
					podUrl := fmt.Sprintf(PodUrl, pod.Status.PodIP, pod.Spec.Containers[0].Ports[0].ContainerPort)
					index, found := findVarnishPod(options.VarnishInstances, podUrl)
					if found {
						removeVarnishPod(&options.VarnishInstances, index)
					}
				}
			}
		},
	})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}
