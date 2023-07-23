import * as k8s from "@pulumi/kubernetes";
import { Config } from "@pulumi/pulumi";
import { argoproj } from "../crds/argocd-application";
import { useNamespace } from "../utils";

// export const namespace = new k8s.core.v1.Namespace(`${config.name}-ns`, {
//     metadata: {
//         name: config.namespace,
//     },
// }, {
//     provider: yaml_provider,
// });

// export const bootstrap = new k8s.yaml.ConfigFile(`${config.name}-bootstrap`, {
//     file: config.manifest,
//     transformations: [useNamespace(namespace)],
// }, {
//     provider: yaml_provider,
// });
