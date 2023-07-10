#!/usr/bin/env bash

set -eax

# clean up any untracked files
minikube delete
# start minikube
minikube start --cpus max --memory max --driver=docker --cni=cilium -p sdp
# use minikube's docker daemon
echo "eval \$(minikube -p sdp docker-env)" >> ~/.zshrc

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

# set stack name to dev if not set
: ${PULUMI_STACK_NAME:=dev}

# create stack if it doesn't exist
pulumi stack select -s $PULUMI_STACK_NAME --create
# update stack
pulumi up -y --refresh -s $PULUMI_STACK_NAME --suppress-outputs
