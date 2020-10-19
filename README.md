# About

When developing on a local k8s instance, often you have to juggle with memory, cpu, ... And when developing with multiple branches, you sometimes have your app installed in multiple namespaces. Each branch, having it's own namespace maybe...

So in order to reduce your resource consumption by your k8s dev cluster, this script allows to downscale all `deployments` and `statefulsets` to zero. It also allows to put them all back at scale `1`.

# Usage

Scale down/up all resources in a k8s namespace

```
./scale.sh -c [CONTEXT] -n [NAMESPACE] -s [down|up]
```