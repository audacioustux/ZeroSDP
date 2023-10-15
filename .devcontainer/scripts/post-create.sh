#!/usr/bin/env bash

set -eax

# set the default profile
minikube profile zerosdp
# minikube extension mounts ~/.minikube to devcontainer
# so need to delete previous profile if it exists
minikube delete
# start minikube
minikube start --cpus 8 --memory 8g --driver=docker --cni=false
# install ciliium
cilium install
# enable hubble
cilium hubble enable
# use minikube's docker daemon
echo "eval \$(minikube docker-env)" >> ~/.zshrc
