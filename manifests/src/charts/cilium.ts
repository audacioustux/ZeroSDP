import { Construct } from 'constructs';
import { Chart, ChartProps, Helm } from 'cdk8s';
import { KubeNamespace } from '../../imports/k8s.js';
import { Application } from "../../imports/argoproj.io.js"

export class ArgoWorkflow extends Chart {
    constructor(scope: Construct, id: string, props: ChartProps = { disableResourceNameHashes: true }) {
        super(scope, id, props);

        const app = new Application(this, 'app', {
            metadata: {
                name: "cilium",
                namespace: "argocd",
                finalizers: ["resources-finalizer.argocd.argoproj.io"]
            },
            spec: {
                destination: {
                    namespace: "kube-system",
                    server: "https://kubernetes.default.svc"
                },
                project: "default",
                source: {
                    chart: "cilium",
                    repoUrl: "https://helm.cilium.io/",
                    targetRevision: '1.14.0',
                },
                syncPolicy: {
                    automated: {
                        prune: true,
                        selfHeal: true
                    }
                }
            }
        });
    }
}