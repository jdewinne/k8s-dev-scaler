[![CI](https://github.com/jdewinne/k8s-dev-scaler/workflows/CI/badge.svg)](https://github.com/jdewinne/k8s-dev-scaler/actions?query=workflow%3ACI)
[![Security](https://github.com/jdewinne/k8s-dev-scaler/workflows/Security%20check%20using%20Snyk/badge.svg)](https://github.com/jdewinne/k8s-dev-scaler/actions?query=workflow%3A%22Security+check+using+Snyk%22)
[![Release k8s-dev-scaler](https://github.com/jdewinne/k8s-dev-scaler/workflows/Release%20k8s-dev-scaler/badge.svg)](https://github.com/jdewinne/k8s-dev-scaler/actions?query=workflow%3A%22Release+k8s-dev-scaler%22)
[![codecov](https://codecov.io/gh/jdewinne/k8s-dev-scaler/branch/main/graph/badge.svg?token=73PVIKFUFD)](https://codecov.io/gh/jdewinne/k8s-dev-scaler)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/jdewinne/k8s-dev-scaler)](https://github.com/jdewinne/k8s-dev-scaler/releases)
[![Downloads](https://img.shields.io/github/downloads/jdewinne/k8s-dev-scaler/total)](https://github.com/jdewinne/k8s-dev-scaler/releases)

# About

When developing on a local k8s instance, often you have to juggle with memory, cpu, ... And when developing with multiple branches, you sometimes have your app installed in multiple namespaces. Each branch, having it's own namespace maybe...

So in order to reduce your resource consumption by your k8s dev cluster, this tool allows to downscale all `deployments` and `statefulsets` to zero. It also allows to scale them all back up. Behind the scenes it places an annotation called `k8s.dev.scaler/desired.replicas` that keeps track of the desired number if replicas.

# Installation

+ Linux: Download from [Releases](https://github.com/jdewinne/k8s-dev-scaler/releases)
+ Linux, Mac: Install using `go get https://github.com/jdewinne/k8s-dev-scaler`

# Usage

Scale down/up all resources in a k8s namespace

```
Usage of k8s-dev-scaler:
  -context string
        (optional) k8s context to be used, current context if not provided.
  -namespace string
        (required) k8s namespace to be used, current namespace if not provided.
  -scale string
        (required) Can be one of [up|down].
```