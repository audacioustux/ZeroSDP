#!/usr/bin/env bash

set -eax

git clean -Xdf --exclude='!**/*.env'

npm install --prefix manifests