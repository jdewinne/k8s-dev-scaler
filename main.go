package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func GetKubeClient(context string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(context string) (*rest.Config, error) {
	config, err := getConfig(context).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

// getConfig returns a Kubernetes client config for a given context.
func getConfig(context string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}

func main() {
	var k8scontextflag *string
	k8scontextflag = flag.String("context", "", "(optional) k8s context to be used, current context if not provided.")

	var namespace *string
	namespace = flag.String("namespace", "", "(required) k8s namespace to be used, current namespace if not provided.")

	var scale *string
	scale = flag.String("scale", "", "(required) Can be one of [up|down].")
	flag.Parse()

	if *namespace == "" {
		flag.Usage()
		panic("Namespace is required.")
	}

	if *scale != "up" && *scale != "down" {
		flag.Usage()
		panic("Scale must be up or down")
	}

	replicas := int32(0)
	if *scale == "up" {
		replicas = 1
	}

	// use the current context in kubeconfig
	_, client, err := GetKubeClient(*k8scontextflag)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Deployments")

	deploymentsClient := client.AppsV1().Deployments(*namespace)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * Scaling %s (%d to %d replicas)\n", d.Name, *d.Spec.Replicas, replicas)
		opts, err := deploymentsClient.GetScale(context.TODO(), d.Name, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		opts.Spec.Replicas = replicas
		deploymentsClient.UpdateScale(context.TODO(), d.Name, opts, metav1.UpdateOptions{})
	}

	fmt.Println("Stateful sets")
	statefulSetsClient := client.AppsV1().StatefulSets(*namespace)
	sslist, err := statefulSetsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, ss := range sslist.Items {
		fmt.Printf(" * Scaling %s (%d to %d replicas)\n", ss.Name, *ss.Spec.Replicas, replicas)
		opts, err := statefulSetsClient.GetScale(context.TODO(), ss.Name, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		opts.Spec.Replicas = replicas
		statefulSetsClient.UpdateScale(context.TODO(), ss.Name, opts, metav1.UpdateOptions{})
	}

}