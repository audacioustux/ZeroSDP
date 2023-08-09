#!/usr/bin/env bash

set -eax

# clean up any untracked files
minikube delete -p sdp
# start minikube
minikube start --nodes 3 --cpus 2 --memory 4g --driver=docker --cni=false -p sdp

# if no backend url is specified - check for access token and use cloud backend
if [[ -z "$PULUMI_BACKEND_URL" ]]; then
    # if no pulumi cloud access token is specified - use local backend
    if [[ -z "$PULUMI_ACCESS_TOKEN" ]]; then
        export PULUMI_BACKEND_URL="file://~"
        # if no passphrase is specified - use default passphrase, and add to zshrc
        if [[ -z "$PULUMI_CONFIG_PASSPHRASE" ]]; then
            export PULUMI_CONFIG_PASSPHRASE="pulumi"
            echo "export PULUMI_CONFIG_PASSPHRASE=$PULUMI_CONFIG_PASSPHRASE" >> ~/.zshrc
        fi
    fi
fi

pulumi login

# set default org if specified
if [[ -n "$PULUMI_DEFAULT_ORG" ]]; then
    pulumi org set-default $PULUMI_DEFAULT_ORG
fi

export PULUMI_STACK_NAME=dev

# bootstrap sdp
npm run bootstrap
