package k8s

import (
	"fmt"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"strings"
	"time"
	"varnish-cache-invalidator/pkg/logging"
)

// VarnishInstances keeps pointer of varnish instances' ip:port information
var VarnishInstances []*string
var logger = logging.GetLogger()

// RunPodInformer continuously watches api-server with shared informer for Pod resources, then does necessary updates
// on Add/Update/Delete conditions
func RunPodInformer(clientSet *kubernetes.Clientset, varnishLabel, varnishNamespace string) {
	varnishLabelKey := strings.Split(varnishLabel, "=")[0]
	varnishLabelValue := strings.Split(varnishLabel, "=")[1]

	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labels := pod.GetLabels()
			for key, value := range labels {
				if key == varnishLabelKey && value == varnishLabelValue && pod.Namespace == varnishNamespace {
					if pod.Status.PodIP != "" {
						podUrl := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, pod.Spec.Containers[0].Ports[0].ContainerPort)
						logger.Info("Adding pod url to the varnishPods slice", zap.String("podUrl", podUrl))
						addVarnishPod(&VarnishInstances, &podUrl)
					} else {
						logger.Warn("Varnish pod does not have an ip address yet, skipping add operation",
							zap.String("pod", pod.Name), zap.String("namespace", pod.Namespace))
					}
				}
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*v1.Pod)
			newPod := newObj.(*v1.Pod)
			labels := oldPod.GetLabels()

			// TODO: Handle all the cases

			for key, value := range labels {
				if key == varnishLabelKey && value == varnishLabelValue && oldPod.ResourceVersion != newPod.ResourceVersion &&
					oldPod.Namespace == varnishNamespace {
					if oldPod.Status.PodIP == "" && newPod.Status.PodIP != "" {
						logger.Info("Assigned an ip address to the pod, adding to varnishPods slice", zap.String("pod", newPod.Name),
							zap.String("namespace", newPod.Namespace), zap.String("ipAddress", newPod.Status.PodIP))
						podUrl := fmt.Sprintf("http://%s:%d", newPod.Status.PodIP, newPod.Spec.Containers[0].Ports[0].ContainerPort)
						logger.Info("Adding pod url to the varnishPods slice", zap.String("podUrl", podUrl))
						addVarnishPod(&VarnishInstances, &podUrl)
					}
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labels := pod.GetLabels()
			for key, value := range labels {
				if key == varnishLabelKey && value == varnishLabelValue && pod.Namespace == varnishNamespace {
					logger.Info("Varnish pod is deleted, removing from varnishPods slice", zap.String("pod", pod.Name),
						zap.String("namespace", pod.Namespace))
					podUrl := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, pod.Spec.Containers[0].Ports[0].ContainerPort)
					index, found := findVarnishPod(VarnishInstances, podUrl)
					if found {
						removeVarnishPod(&VarnishInstances, index)
					}
				}
			}
		},
	})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}
