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

// ScaleDeploymentResources will scale all deployments up or down in a namespace
func (ds *DeploymentsScaler) ScaleDeploymentResources() {
	resources, err := ds.client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, r := range resources.Items {
		// store original desired number of replicas as an annotation
		if ds.scale == "down" {
			err = ds.annotateResource(r.Name, *r.Spec.Replicas)
			if err != nil {
				panic(err.Error())
			}
		}
		// If scaling up, get the replicas from the previously stored annotation
		replicas := int32(0)
		if ds.scale == "up" {
			replicas, err = ds.getDesiredReplicas(r.Name)
			if err != nil {
				panic(err.Error())
			}
		}
		fmt.Printf(" * Scaling %s (%d to %d replicas)\n", r.Name, *r.Spec.Replicas, replicas)
		opts, err := ds.client.GetScale(context.TODO(), r.Name, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		opts.Spec.Replicas = replicas
		ds.client.UpdateScale(context.TODO(), r.Name, opts, metav1.UpdateOptions{})
	}
}
