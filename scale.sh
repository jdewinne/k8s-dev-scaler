#!/bin/bash

usage () { 
    echo "Scale down/up all resources in a k8s namespace"
    echo "scale.sh -c [CONTEXT] -n [NAMESPACE] -s [down|up]"; 
}

# Inputs: context - namespace - down or up
while getopts ":c:n:s:h" opt; do
  case $opt in
    c) 
      context="$OPTARG"
      ;;
    n)
      namespace="$OPTARG"
      ;;
    s)
      scale="$OPTARG"
      ;;
    h)
      usage
      exit 0
      ;;
    \?) 
      echo "Invalid option -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
    ;;
  esac
done

if ((OPTIND == 1))
then
    echo "No options specified"
    usage
fi

deployments() {
    deployments=$(kubectl get deployments --context $context --namespace $namespace -o name)
}

statefulsets() {
    statefulsets=$(kubectl get statefulsets --context $context --namespace $namespace -o name)
}

scale() {
    for deployment in $deployments
    do
        kubectl scale $deployment --replicas $1 --context $context --namespace $namespace
    done
    for statefulset in $statefulsets
    do
        kubectl scale $statefulset --replicas $1 --context $context --namespace $namespace
    done
}

deployments
statefulsets
if [  "$scale" == "down" ];
then
    scale 0
elif [ "$scale" == "up" ];
then
    scale 1
else
    echo "Invalid scaling option."
    usage
    exit 1
fi