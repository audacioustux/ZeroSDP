import { Construct } from 'constructs';
import { ApplicationChart } from './_application-chart.js';
import { ApplicationProps } from '../../imports/argocd-argoproj.io.js';

export class ArgoCD extends ApplicationChart {
    constructor(scope: Construct, id: string, props = {}) {
        const appProps: ApplicationProps = {
            metadata: {
                name: "argo-cd",
                namespace: "argocd",
            },
            spec: {
                project: "default",
                destination: {
                    namespace: "argocd",
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
        };

        super(scope, id, appProps, props);
    };
}