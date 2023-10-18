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
# use minikube's docker daemon
LINE='eval $(minikube docker-env)'
FILE=~/.zshrc
grep -xqF -- "$LINE" "$FILE" || echo "$LINE" >> "$FILE"
