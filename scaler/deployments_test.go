package scaler

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func deployment(name string, replicas int32) *v1.Deployment {
	return &v1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"}, Spec: v1.DeploymentSpec{Replicas: &replicas}}
}

func scale(name string, replicas int32) *autoscalingv1.Scale {
	return &autoscalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"}, Spec: autoscalingv1.ScaleSpec{Replicas: replicas}}
}

func TestAnnotateDeployment(t *testing.T) {
	var tests = []struct {
		description      string
		expectedNames    []string
		expectedReplicas []int32
		deployments      []runtime.Object
	}{
		{"single deployment", []string{"d1"}, []int32{1}, []runtime.Object{deployment("d1", 1)}},
		{"multiple deployment", []string{"d1", "d2"}, []int32{1, 1}, []runtime.Object{deployment("d1", 1), deployment("d2", 1)}},
		{"different replicas", []string{"d1", "d2"}, []int32{3, 1}, []runtime.Object{deployment("d1", 3), deployment("d2", 1)}},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := testclient.NewSimpleClientset(test.deployments...)
			dscaler := NewDeploymentsScaler(client, "default", "down")
			for i, d := range test.deployments {
				dscaler.annotateResource(d.(*v1.Deployment).Name, *d.(*v1.Deployment).Spec.Replicas)
				deployment, _ := dscaler.client.Get(context.TODO(), d.(*v1.Deployment).ObjectMeta.Name, metav1.GetOptions{})
				replicas, _ := strconv.Atoi(deployment.Annotations["k8s.dev.scaler/desired.replicas"])
				assert.Equal(t, test.expectedNames[i], deployment.Name)
				assert.Equal(t, test.expectedReplicas[i], int32(replicas))
			}

		})
	}

}

func TestGetReplicaScaleDeployment(t *testing.T) {
	var tests = []struct {
		description      string
		expectedReplicas []int32
		deployments      []runtime.Object
	}{
		{"single deployment", []int32{1}, []runtime.Object{deployment("d1", 1)}},
		{"multiple deployment", []int32{1, 1}, []runtime.Object{deployment("d1", 1), deployment("d2", 1)}},
		{"different replicas", []int32{3, 1}, []runtime.Object{deployment("d1", 3), deployment("d2", 1)}},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := testclient.NewSimpleClientset(test.deployments...)
			downscaler := NewDeploymentsScaler(client, "default", "down")
			upscaler := NewDeploymentsScaler(client, "default", "up")
			for i, d := range test.deployments {
				assert.Equal(t, int32(0), downscaler.getReplicaScale(d.(*v1.Deployment).Name, *d.(*v1.Deployment).Spec.Replicas))
				assert.Equal(t, test.expectedReplicas[i], upscaler.getReplicaScale(d.(*v1.Deployment).Name, *d.(*v1.Deployment).Spec.Replicas))
			}

		})
	}

}

func TestRunScalerDown(t *testing.T) {
	t.Skip("Skipping testing as autoscaling.Scale doesn't work properly yet with the fake k8s api")
	var tests = []struct {
		description      string
		expectedReplicas []int32
		scales           []runtime.Object
	}{
		{"single deployment", []int32{0}, []runtime.Object{scale("d1", 1)}},
		{"multiple deployment", []int32{0, 0}, []runtime.Object{scale("d1", 1), scale("d2", 1)}},
		{"different replicas", []int32{0, 0}, []runtime.Object{scale("d1", 3), scale("d2", 1)}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := testclient.NewSimpleClientset(test.scales...)
			downscaler := NewDeploymentsScaler(client, "default", "down")
			for i, d := range test.scales {
				downscaler.runScaler(d.(*autoscalingv1.Scale).Name, int32(0))
				replicas, _ := downscaler.client.GetScale(context.TODO(), d.(*autoscalingv1.Scale).Name, metav1.GetOptions{})
				assert.Equal(t, test.expectedReplicas[i], replicas.Spec.Replicas)
			}

		})
	}
}

func TestRunScalerUp(t *testing.T) {
	t.Skip("Skipping testing as autoscaling.Scale doesn't work properly yet with the fake k8s api")
	var tests = []struct {
		description      string
		expectedReplicas []int32
		scales           []runtime.Object
	}{
		{"single deployment", []int32{1}, []runtime.Object{scale("d1", 0)}},
		{"multiple deployment", []int32{1, 1}, []runtime.Object{scale("d1", 0), scale("d2", 0)}},
		{"different replicas", []int32{3, 1}, []runtime.Object{scale("d1", 0), scale("d2", 0)}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := testclient.NewSimpleClientset(test.scales...)
			upscaler := NewDeploymentsScaler(client, "default", "up")
			for i, d := range test.scales {
				upscaler.runScaler(d.(*autoscalingv1.Scale).Name, test.expectedReplicas[i])
				replicas, _ := upscaler.client.GetScale(context.TODO(), d.(*autoscalingv1.Scale).Name, metav1.GetOptions{})
				assert.Equal(t, test.expectedReplicas[i], replicas.Spec.Replicas)
			}

		})
	}
}
