package k8s

import (
	"context"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"sync"
	"testing"
	"time"
)

type FakeAPI struct {
	ClientSet kubernetes.Interface
	Namespace string
}

func getFakeAPI() *FakeAPI {
	client := fake.NewSimpleClientset()
	api := &FakeAPI{ClientSet: client, Namespace: "default"}
	return api
}

func (fAPI *FakeAPI) deletePod(name string) error {
	gracePeriodSeconds := int64(0)
	return fAPI.ClientSet.CoreV1().Pods(fAPI.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds})
}

func (fAPI *FakeAPI) getPod(name string) (*v1.Pod, error) {
	return fAPI.ClientSet.CoreV1().Pods(fAPI.Namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (fAPI *FakeAPI) updatePod(name, podIP string) (*v1.Pod, error) {
	pod, _ := fAPI.getPod(name)
	pod.Status.PodIP = podIP
	pod.ResourceVersion = "123456"

	pod, err := fAPI.ClientSet.CoreV1().Pods(fAPI.Namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (fAPI *FakeAPI) createPod(name, ip string) (*v1.Pod, error) {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: fAPI.Namespace,
			Labels: map[string]string{
				"app": "varnish",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "varnish",
					Image:           "varnish:7.0.1",
					ImagePullPolicy: "Always",
					Ports: []v1.ContainerPort{
						{Name: "port1", ContainerPort: 6082, Protocol: v1.ProtocolTCP},
					},
				},
			},
		},
		Status: v1.PodStatus{
			PodIP: ip,
		},
	}

	pod, err := fAPI.ClientSet.CoreV1().Pods(fAPI.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func TestRunPodInformer(t *testing.T) {
	api := getFakeAPI()
	assert.NotNil(t, api)

	cases := []struct {
		caseName, podName, ip string
	}{
		{
			caseName: "caseNonEmptyIPcreatePod",
			ip:       "10.0.0.15",
			podName:  "varnish-pod-1",
		},
		{
			caseName: "caseEmptyIPcreatePod",
			ip:       "",
			podName:  "varnish-pod-2",
		},
	}

	go func() {
		RunPodInformer(api.ClientSet)
	}()

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				pod, err := api.createPod(tc.podName, tc.ip)
				assert.Nil(t, err)
				assert.NotNil(t, pod)
			}()
			wg.Wait()

			time.Sleep(2 * time.Second)

			wg.Add(1)
			go func() {
				defer wg.Done()
				createdPod, err := api.getPod(tc.podName)
				assert.NotNil(t, createdPod)
				assert.Nil(t, err)
			}()
			wg.Wait()

			time.Sleep(2 * time.Second)

			wg.Add(1)
			go func() {
				defer wg.Done()
				updatedPod, err := api.updatePod(tc.podName, "10.0.0.15")
				assert.NotNil(t, updatedPod)
				assert.Nil(t, err)
			}()
			wg.Wait()

			time.Sleep(2 * time.Second)

			wg.Add(1)
			go func() {
				defer wg.Done()
				err := api.deletePod(tc.podName)
				assert.Nil(t, err)
			}()
			wg.Wait()
		})
	}
}

func TestGetClientSet(t *testing.T) {
	opts.IsLocal = true
	opts.KubeConfigPath = "../../mock/kubeconfig"

	restConfig, err := GetConfig()
	assert.Nil(t, err)
	assert.NotNil(t, restConfig)

	clientSet, err := GetClientSet(restConfig)
	assert.Nil(t, err)
	assert.NotNil(t, clientSet)
}
