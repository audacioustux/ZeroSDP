import * as k8s from "@pulumi/kubernetes";
import { useNamespace } from "../utils";

export const config = {
    name: "cert-manager",
    manifest: {
        crds: "https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.crds.yaml",
        helm: {
            chart: "cert-manager",
            repoUrl: "https://charts.jetstack.io",
        }
    }
}

export const namespace = new k8s.core.v1.Namespace(`${config.name}-ns`, {
    metadata: {
        name: config.name,
    },
});

export const crds = new k8s.yaml.ConfigFile(`${config.name}-crds`, {
    file: config.manifest.crds
});
