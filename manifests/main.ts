import { App, AppProps, Chart, ChartProps, Helm, HelmProps } from 'cdk8s';
import { KubeNamespace, KubeNamespaceProps } from './imports/k8s.js';
import { Application, ApplicationProps } from './imports/argocd-argoproj.io.js';
import { Construct } from 'constructs';

export type HelmChartProps = HelmProps & { releaseName: string };
export class HelmChart extends Chart {
    constructor(scope: Construct, helmProps: HelmChartProps, props: ChartProps = {}) {
        const { releaseName } = helmProps;
        super(scope, `helm-${releaseName}`, { ...props, disableResourceNameHashes: true });

        new Helm(this, 'helm', helmProps);
    }
}

export type ArgoApplicationProps = ApplicationProps & { metadata: { name: string } };
export class ArgoApplication extends Chart {
    constructor(scope: Construct, appProps: ArgoApplicationProps, props: ChartProps = {}) {
        const { metadata: { name } } = appProps;
        super(scope, `app-${name}`, props);

        new Application(this, 'app', appProps);
    }
}

export type NamespaceProps = KubeNamespaceProps & { metadata: { name: string } };
export class Namespace extends Chart {
    constructor(scope: Construct, namespaceProps: NamespaceProps, props: ChartProps = {}) {
        const { metadata: { name } } = namespaceProps;
        super(scope, `namespace-${name}`, props);

        new KubeNamespace(this, 'argocd-namespace', namespaceProps);
    }
}

export class SDP extends App {
    constructor(props?: AppProps) {
        super(props);

        const ns = new Namespace(this, { metadata: { name: "argocd" } });
        const helm = new HelmChart(this, {
            repo: 'https://argoproj.github.io/argo-helm',
            chart: 'argo-cd',
            releaseName: 'argocd',
            namespace: 'argocd',
            values: {
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
            },
        });
        const app = new ArgoApplication(this, {
            metadata: {
                name: "argo-cd",
                namespace: "argocd",
            },
            spec: {
                project: "default",
                destination: {
                    server: "https://kubernetes.default.svc",
                },
                source: {
                    repoUrl: "https://github.com/audacioustux/sdp.git",
                    path: "manifests/dist"
                },
                syncPolicy: {
                    automated: {
                        prune: true,
                        selfHeal: true,
                    }
                }
            }
        });

        helm.addDependency(ns);
        app.addDependency(helm);
    }
}

const app = new SDP();
app.synth();

console.log("hello")
