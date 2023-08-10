import { ApiObjectMetadata, App, AppProps, Chart, ChartProps, Helm, HelmProps } from 'cdk8s';
import * as helm from './charts/helm/index.js';
import * as application from './charts/application/index.js';
import { KubeNamespace, KubeNamespaceProps } from './imports/k8s.js';
import { Application, ApplicationProps } from './imports/argocd-argoproj.io.js';
import { Construct } from 'constructs';
import { MergeDeep, SetRequired } from 'type-fest';

export class HelmChart extends Chart {
    constructor(scope: Construct, helmProps: SetRequired<HelmProps, | 'releaseName'>, props: ChartProps = {}) {
        const { repo, releaseName } = helmProps;

        super(scope, `helm-${releaseName}`, { ...props, disableResourceNameHashes: true });

        new Helm(this, 'helm', helmProps);
    }
}

export class ArgoApplication extends Chart {
    constructor(scope: Construct, appProps: ApplicationProps & { metadata: SetRequired<ApiObjectMetadata, 'name'> }, props: ChartProps = {}) {
        const { metadata: { name } } = appProps;

        super(scope, `app-${name}`, props);

        new Application(this, 'app', appProps);
    }
}

export class Namespace extends Chart {
    constructor(scope: Construct, namespaceProps: KubeNamespaceProps & { metadata: SetRequired<ApiObjectMetadata, 'name'> }, props: ChartProps = {}) {
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
                    path: "manifests/dist",
                    directory: {
                        include: "argo-cd.k8s.yaml"
                    }
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