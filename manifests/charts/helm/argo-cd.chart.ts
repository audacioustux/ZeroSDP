import { Construct } from 'constructs';
import { HelmChart } from './_helm-chart.js';

export interface ArgoCDProps {
    ha?: true | { autoscaling: boolean }
    namespace?: string
}

export class ArgoCD extends HelmChart {
    constructor(scope: Construct, id: string, props: ArgoCDProps = {}) {
        const { ha, namespace = "argocd" } = props;

        const ha_values: Record<string, any> = {
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
        const values = ha ? ha_values : {};

        super(scope, id, {
            releaseName: "argocd",
            chart: "argo/argo-cd",
            namespace,
            values
        });
    }
}