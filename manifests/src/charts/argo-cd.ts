import { Construct } from 'constructs';
import { Chart, ChartProps, Helm } from 'cdk8s';
import { KubeNamespace } from '../../imports/k8s.js';
import { Application } from "../../imports/argocd-argoproj.io.js"

export interface ArgoCDProps extends ChartProps {
    ha: boolean;
}

export class ArgoCD extends Chart {
    constructor(scope: Construct, id: string, props: ArgoCDProps = { ha: true, disableResourceNameHashes: true }) {
        super(scope, id, props);

        const namespace = new KubeNamespace(this, 'ns', {
            metadata: {
                name: 'argocd',
            },
        });

        const ha_values = {
            "redis-ha": {
                enabled: true
            },
            controller: {
                replicas: 1
            },
            server: {
                autoscaling: {
                    enabled: true,
                    minReplicas: 2
                },
            },
            repoServer: {
                autoscaling: {
                    enabled: true,
                    minReplicas: 2
                },
            },
            applicationSet: {
                replicaCount: 2
            }
        };

        const helm = new Helm(this, 'helm', {
            releaseName: "argocd",
            chart: "argo/argo-cd",
            helmFlags: [
                "--namespace", namespace.name
            ],
            values: props.ha ? ha_values : {},
        });

        const app = new Application(this, 'app', {
            metadata: {
                name: helm.releaseName,
                namespace: namespace.name,
            },
            spec: {
                destination: {
                    namespace: namespace.name,
                    server: 'https://kubernetes.default.svc',
                },
                project: 'default',
                source: {
                    path: 'manifests/dist',
                    directory: {
                        include: "argo-cd.k8s.yaml"
                    },
                    repoUrl: "https://github.com/audacioustux/sdp.git",
                    targetRevision: 'HEAD',
                },
                syncPolicy: {
                    automated: {
                        prune: true,
                        selfHeal: true,
                    }
                },
            },
        });
    }
}