import * as k8s from "@pulumi/kubernetes";
import { Config } from "@pulumi/pulumi";
import { argoproj } from "../crds/argocd-application";
import { useNamespace } from "../utils";

export const config = {
    name: "argo-cd",
    namespace: "argocd",
    manifest: "https://raw.githubusercontent.com/argoproj/argo-cd/v2.7.7/manifests/install.yaml"
}

export const namespace = new k8s.core.v1.Namespace(`${config.name}-ns`, {
    metadata: {
        name: config.namespace,
    },
});

export const bootstrap = new k8s.yaml.ConfigFile(`${config.name}-bootstrap`, {
    file: config.manifest,
    transformations: [useNamespace(namespace)],
});

// export const app = new argoproj.v1alpha1.Application(`${config.name}-app`, {
//     metadata: {
//         namespace: config.namespace,
//     },
//     spec: {
//         destination: {
//             namespace: config.namespace,
//             server: "https://kubernetes.default.svc",
//         },
//         source: {
//             repoURL: "https://github.com/argoproj/argo-cd.git",
//             targetRevision: "v2.7.7",
//             path: "manifests/",
//             directory: {
//                 include: "install.yaml"
//             }
//         },
//         project: "default",
//         syncPolicy: {
//             automated: {
//                 prune: true,
//                 selfHeal: true,
//             },
//         },
//     },
// });
