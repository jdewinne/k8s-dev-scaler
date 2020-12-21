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

// DeploymentsScaler allows to scale up or down all deployments
type DeploymentsScaler struct {

	// defining struct variables
	client    v1.DeploymentInterface
	namespace string
	scale     string
}

// NewDeploymentsScaler instantiates
func NewDeploymentsScaler(client kubernetes.Interface, namespace string, scale string) *DeploymentsScaler {
	p := new(DeploymentsScaler)
	p.client = client.AppsV1().Deployments(namespace)
	p.namespace = namespace
	p.scale = scale
	return p
}

func (ds *DeploymentsScaler) annotateResource(name string, replicas int32) error {
	payload := fmt.Sprintf(`{"metadata":{"annotations":{"k8s.dev.scaler/desired.replicas":"%d"}}}`, replicas)
	_, err := ds.client.Patch(context.TODO(), name, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
	return err
}

func (ds *DeploymentsScaler) getDesiredReplicas(name string) (int32, error) {
	deployment, err := ds.client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	replicas, _ := strconv.Atoi(deployment.Annotations["k8s.dev.scaler/desired.replicas"])
	return int32(replicas), nil
}

func (ds *DeploymentsScaler) getReplicaScale(name string, replicas int32) int32 {
	// store original desired number of replicas as an annotation
	if ds.scale == "down" {
		err := ds.annotateResource(name, replicas)
		if err != nil {
			panic(err.Error())
		}
	}
	// If scaling up, get the replicas from the previously stored annotation
	nreps := int32(0)
	if ds.scale == "up" {
		nreps, err := ds.getDesiredReplicas(name)
		if err != nil {
			panic(err.Error())
		}
		return nreps
	}
	return nreps
}

func (ds *DeploymentsScaler) runScaler(name string, replicas int32) {
	opts, err := ds.client.GetScale(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf(" * Scaling %s (%d to %d replicas)\n", name, opts.Spec.Replicas, replicas)
	opts.Spec.Replicas = replicas
	ds.client.UpdateScale(context.TODO(), name, opts, metav1.UpdateOptions{})
}

// ScaleDeploymentResources will scale all deployments up or down in a namespace
func (ds *DeploymentsScaler) ScaleDeploymentResources() {
	resources, err := ds.client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, r := range resources.Items {
		// If scaling up, get the replicas from the previously stored annotation, else this returns zero
		replicas := ds.getReplicaScale(r.Name, *r.Spec.Replicas)

		ds.runScaler(r.Name, replicas)
	}
}
