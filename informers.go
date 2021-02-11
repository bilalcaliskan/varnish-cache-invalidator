package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"log"
	"time"
)

func runPodInformer(clientSet *kubernetes.Clientset, varnishLabelKey, varnishLabelValue, varnishNamespace string) {
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second * 30)
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labels := pod.GetLabels()
			for key, value := range labels {
				if key == varnishLabelKey && value == varnishLabelValue && pod.Namespace == varnishNamespace {
					if pod.Status.PodIP != "" {
						podString := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, pod.Spec.Containers[0].Ports[0].ContainerPort)
						log.Printf("adding podString %v to the varnishPods slice\n", podString)
						addVarnishPod(&varnishInstances, &podString)
					} else {
						log.Printf("varnish pod %v on namespace %v not have an ip address yet, skipping add " +
							"operation\n", pod.Name, pod.Namespace)
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
						log.Printf("assigned an ip address to the pod %v on namespace %v, adding to varnishPods " +
							"slice\n", newPod.Name, newPod.Namespace)
						podString := fmt.Sprintf("http://%s:%d", newPod.Status.PodIP, newPod.Spec.Containers[0].Ports[0].ContainerPort)
						log.Printf("adding podString %v to the varnishPods slice\n", podString)
						addVarnishPod(&varnishInstances, &podString)
					}
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labels := pod.GetLabels()
			for key, value := range labels {
				if key == varnishLabelKey && value == varnishLabelValue && pod.Namespace == varnishNamespace {
					log.Printf("pod %v is deleted on namespace %v, removing from varnishPods slice!\n", pod.Name,
						pod.Namespace)
					podString := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, pod.Spec.Containers[0].Ports[0].ContainerPort)
					index, found := findVarnishPod(varnishInstances, podString)
					if found {
						log.Printf("deleted pod %v found on the varnishPods slice, removing!\n", pod.Name)
						removeVarnishPod(&varnishInstances, index)
					}
				}
			}
		},
	})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}