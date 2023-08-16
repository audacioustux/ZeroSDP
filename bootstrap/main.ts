import * as k8s from "@pulumi/kubernetes";
import { Config } from "@pulumi/pulumi";

const k8sProvider = new k8s.Provider("k8s");

let config = new Config();
let manifests = config.requireObject<string[]>("manifests");

new k8s.yaml.ConfigGroup("argo-cd", {
    files: manifests,
}, { provider: k8sProvider });
