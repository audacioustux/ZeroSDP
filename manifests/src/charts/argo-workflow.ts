import { Construct } from 'constructs';
import { Chart, ChartProps, Helm } from 'cdk8s';
import { KubeNamespace } from '../../imports/k8s.js';
import { Application } from "../../imports/argoproj.io.js"

export class ArgoWorkflow extends Chart {
    constructor(scope: Construct, id: string, props: ChartProps = { disableResourceNameHashes: true }) {
        super(scope, id, props);

        const app = new Application(this, 'app', {
            metadata: {
                name: "argo-workflows",
                namespace: "argo",
            },
            spec: {
                destination: {
                    namespace: "argo",
                    server: "https://kubernetes.default.svc"
                },
                project: "default",
                source: {
                    chart: "argo-workflows",
                    repoUrl: "https://argoproj.github.io/argo-helm",
                    targetRevision: 'v3.4.9'
                },
                syncPolicy: {
                    automated: {
                        prune: true,
                        selfHeal: true
                    },
                    syncOptions: [
                        "CreateNamespace=true"
                    ]
                }
            }
        });
    }
}