import * as k8s from "@pulumi/kubernetes";
import { useNamespace } from "@audacioustux/pulumi-k8s-utils"

const provider = new k8s.Provider("provider", {
    kubeconfig: "",
    renderYamlToDirectory: "./rendered/argo-cd",
});

const config = {
    name: "argo-cd",
    namespace: "argocd",
    manifest: "https://raw.githubusercontent.com/argoproj/argo-cd/v2.7.7/manifests/install.yaml"
}

const namespace = new k8s.core.v1.Namespace(`${config.name}-ns`, {
    metadata: {
        name: config.namespace,
    },
}, { provider });

const bootstrap = new k8s.yaml.ConfigFile(`${config.name}-bootstrap`, {
    file: config.manifest,
    transformations: [useNamespace(namespace)],
}, { provider });
