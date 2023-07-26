import * as k8s from "@pulumi/kubernetes";
import { getStack, StackReference, getOrganization } from "@pulumi/pulumi";

const manifest_stack = new StackReference(`${getOrganization()}/manifests/${getStack()}`);

const argocd = await manifest_stack.getOutputValue("argocd")

const argo_cd = new k8s.yaml.ConfigGroup("argo-cd", {
    files: [`../manifests/${argocd}/0-crd/*.yaml`, `../manifests/${argocd}/1-manifest/*.yaml`]
});
