package scaler

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func statefulSet(name string, replicas int32) *v1.StatefulSet {
	return &v1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"}, Spec: v1.StatefulSetSpec{Replicas: &replicas}}
}

func TestAnnotateStatefulSet(t *testing.T) {
	var tests = []struct {
		description      string
		expectedNames    []string
		expectedReplicas []int32
		statefulSets     []runtime.Object
	}{
		{"single stateful set", []string{"d1"}, []int32{1}, []runtime.Object{statefulSet("d1", 1)}},
		{"multiple stateful sets", []string{"d1", "d2"}, []int32{1, 1}, []runtime.Object{statefulSet("d1", 1), statefulSet("d2", 1)}},
		{"different replicas", []string{"d1", "d2"}, []int32{3, 1}, []runtime.Object{statefulSet("d1", 3), statefulSet("d2", 1)}},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := testclient.NewSimpleClientset(test.statefulSets...)
			sscaler := NewStatefulSetsScaler(client, "default", "down")
			for i, s := range test.statefulSets {
				sscaler.annotateResource(s.(*v1.StatefulSet).Name, *s.(*v1.StatefulSet).Spec.Replicas)
				deployment, _ := sscaler.client.Get(context.TODO(), s.(*v1.StatefulSet).ObjectMeta.Name, metav1.GetOptions{})
				replicas, _ := strconv.Atoi(deployment.Annotations["k8s.dev.scaler/desired.replicas"])
				assert.Equal(t, test.expectedNames[i], deployment.Name)
				assert.Equal(t, test.expectedReplicas[i], int32(replicas))
			}

		})
	}

}

func TestGetReplicaScaleStatefulSet(t *testing.T) {
	var tests = []struct {
		description      string
		expectedReplicas []int32
		statefulSets     []runtime.Object
	}{
		{"single stateful set", []int32{1}, []runtime.Object{statefulSet("d1", 1)}},
		{"multiple stateful sets", []int32{1, 1}, []runtime.Object{statefulSet("d1", 1), statefulSet("d2", 1)}},
		{"different replicas", []int32{3, 1}, []runtime.Object{statefulSet("d1", 3), statefulSet("d2", 1)}},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := testclient.NewSimpleClientset(test.statefulSets...)
			downscaler := NewStatefulSetsScaler(client, "default", "down")
			upscaler := NewStatefulSetsScaler(client, "default", "up")
			for i, d := range test.statefulSets {
				assert.Equal(t, int32(0), downscaler.getReplicaScale(d.(*v1.StatefulSet).Name, *d.(*v1.StatefulSet).Spec.Replicas))
				assert.Equal(t, test.expectedReplicas[i], upscaler.getReplicaScale(d.(*v1.StatefulSet).Name, *d.(*v1.StatefulSet).Spec.Replicas))
			}

		})
	}

}
