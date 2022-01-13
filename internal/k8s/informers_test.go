package k8s

import (
	"context"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

type FakeAPI struct {
	ClientSet kubernetes.Interface
	Namespace string
}

/*func (fAPI *FakeAPI) deletePod(name string) error {
	// gracePeriodSeconds := int64(0)
	return fAPI.ClientSet.CoreV1().Pods(fAPI.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}*/

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

func fakeAPI() *FakeAPI {
	client := fake.NewSimpleClientset()
	api := &FakeAPI{ClientSet: client, Namespace: "default"}
	return api
}

func TestRunPodInformerCreatePod(t *testing.T) {
	t.Parallel()
	api := fakeAPI()
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
			pod, err := api.createPod(tc.podName, tc.ip)
			assert.Nil(t, err)
			assert.NotNil(t, pod)

			/*err = api.deletePod(tc.podName)
			assert.Nil(t, err)*/
		})
	}

	<-time.After(10 * time.Second)
}

/*func TestRunPodInformerDeletePod(t *testing.T) {
	t.Parallel()
	api := fakeAPI()
	assert.NotNil(t, api)

	cases := []struct {
		caseName, podName, ip string
	}{
		{
			caseName: "case1",
			ip:   "10.0.0.15",
			podName: "varnish-pod-1",
		},
	}

	go func() {
		RunPodInformer(api.ClientSet)
	}()

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			pod, err := api.createPod(tc.podName, tc.ip)
			assert.Nil(t, err)
			assert.NotNil(t, pod)

			err = api.deletePod(tc.podName)
			assert.Nil(t, err)
		})
	}

	<- time.After(10 * time.Second)
}*/

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
