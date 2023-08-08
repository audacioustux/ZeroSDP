import * as k8s from "@pulumi/kubernetes";
import { Config } from "@pulumi/pulumi";

const k8sProvider = new k8s.Provider("k8s");

let config = new Config();
let manifest = config.require("manifest");

new k8s.yaml.ConfigFile("argo-cd", {
    file: manifest,
}, { provider: k8sProvider });
