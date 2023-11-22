#!/usr/bin/env bash

set -eax

# set the default profile
minikube profile zerosdp
# minikube feature mounts ~/.minikube to devcontainer
# so need to delete previous profile if it exists
minikube delete
# start minikube
minikube start --cpus 6 --memory 8g --driver=docker --cni=false
# install ciliium
cilium install
# enable hubble
cilium hubble enable
# enable volume snapshots
minikube addons enable volumesnapshots
# use csi hostpath driver
minikube addons enable csi-hostpath-driver
# set the default storage class
minikube addons disable storage-provisioner
minikube addons disable default-storageclass
kubectl patch storageclass csi-hostpath-sc -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
# use minikube's docker daemon
LINE='eval $(minikube docker-env)'
FILE=~/.zshrc
grep -xqF -- "$LINE" "$FILE" || echo "$LINE" >> "$FILE"