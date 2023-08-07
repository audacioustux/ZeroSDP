#!/usr/bin/env bash

set -eax

command -v \
    helm \
    jq &>/dev/null

yq '.repos[] | .name + " " + .url' $1 | awk '{system("helm repo add "$1" "$2)}'