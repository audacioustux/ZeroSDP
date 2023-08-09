import { Construct } from 'constructs';
import { Chart, ChartProps, Helm } from 'cdk8s';
import { ApplicationChart } from './_application-chart.js';

export class ArgoWorkflow extends Chart {
    constructor(scope: Construct, id: string, props: ChartProps = {}) {
        super(scope, id, props);

        const app = new ApplicationChart(this, 'app', {
            metadata: {
                name: "argo-workflows",
            },
            spec: {
                project: "default",
                destination: {
                    namespace: "argo",
                    server: "https://kubernetes.default.svc",
                },
                source: {
                    repoUrl: "https://github.com/audacioustux/sdp.git",
                    path: "manifests/helm",
                    directory: {
                        include: "argo-workflows.k8s.yaml"
                    }
                }
            }
        }, props);
    };
}