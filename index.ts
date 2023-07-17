import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

import "./bootstrap"

import "./apps/cert-manager"
import "./apps/argocd"