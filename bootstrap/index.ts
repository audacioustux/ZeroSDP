import * as k8s from "@pulumi/kubernetes";
import * as kx from "@pulumi/kubernetesx";
import { getStack, StackReference, getOrganization } from "@pulumi/pulumi";

export const other = new StackReference(`${getOrganization()}/manifests/${getStack()}`);

