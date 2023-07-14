import * as k8s from "@pulumi/kubernetes";
import { useNamespace } from "../utils";

export const config = {
    name: "argocd",
    namespace: "argocd",
    manifest: {
        file: "https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml"
    }
}

export const namespace = new k8s.core.v1.Namespace(`${config.name}-ns`, {
    metadata: {
        name: config.namespace,
    },
});

export const bootstrap = new k8s.yaml.ConfigFile(`${config.name}-bootstrap`, {
    file: config.manifest.file,
    transformations: [useNamespace(namespace)],
});
