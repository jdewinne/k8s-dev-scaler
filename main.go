package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	scaler "github.com/jdewinne/k8s-dev-scaler/scaler"
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

	// use the current context in kubeconfig
	_, client, err := GetKubeClient(*k8scontextflag)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Deployments")
	dscaler := scaler.NewDeploymentsScaler(client, *namespace, *scale)
	dscaler.ScaleDeploymentResources()

	fmt.Println("Stateful sets")
	sscaler := scaler.NewStatefulSetsScaler(client, *namespace, *scale)
	sscaler.ScaleStatefulSetResources()

}
