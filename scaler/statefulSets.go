package scaler

import (
	"context"
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

// StatefulSetsScaler allows to scale up or down all statefulSets
type StatefulSetsScaler struct {

	// defining struct variables
	client    v1.StatefulSetInterface
	namespace string
	scale     string
}

// NewStatefulSetsScaler instantiates
func NewStatefulSetsScaler(client kubernetes.Interface, namespace string, scale string) *StatefulSetsScaler {
	p := new(StatefulSetsScaler)
	p.client = client.AppsV1().StatefulSets(namespace)
	p.namespace = namespace
	p.scale = scale
	return p
}

// annotateResource places the k8s.dev.scaler/desired.replicas annotation
func (ss *StatefulSetsScaler) annotateResource(name string, replicas int32) error {
	payload := fmt.Sprintf(`{"metadata":{"annotations":{"k8s.dev.scaler/desired.replicas":"%d"}}}`, replicas)
	_, err := ss.client.Patch(context.TODO(), name, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
	return err
}

// getDesiredReplicas fetches the value from the k8s.dev.scaler/desired.replicas annotation
func (ss *StatefulSetsScaler) getDesiredReplicas(name string) (int32, error) {
	statefulSet, err := ss.client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	replicas, _ := strconv.Atoi(statefulSet.Annotations["k8s.dev.scaler/desired.replicas"])
	return int32(replicas), nil
}

// ScaleStatefulSetResources will scale all deployments up or down in a namespace
func (ss *StatefulSetsScaler) ScaleStatefulSetResources() {
	resources, err := ss.client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, r := range resources.Items {
		// store original desired number of replicas as an annotation
		if ss.scale == "down" {
			err = ss.annotateResource(r.Name, *r.Spec.Replicas)
			if err != nil {
				panic(err.Error())
			}
		}
		// If scaling up, get the replicas from the previously stored annotation
		replicas := int32(0)
		if ss.scale == "up" {
			replicas, err = ss.getDesiredReplicas(r.Name)
			if err != nil {
				panic(err.Error())
			}
		}
		fmt.Printf(" * Scaling %s (%d to %d replicas)\n", r.Name, *r.Spec.Replicas, replicas)
		opts, err := ss.client.GetScale(context.TODO(), r.Name, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		opts.Spec.Replicas = replicas
		ss.client.UpdateScale(context.TODO(), r.Name, opts, metav1.UpdateOptions{})
	}
}
