import { Construct } from 'constructs';
import { ChartProps } from 'cdk8s';
import { HelmChart } from './_helm-chart.js';

export interface ArgoCDProps extends ChartProps {
    ha?: boolean;
}

export class ArgoCD extends HelmChart {
    constructor(scope: Construct, id: string, props: ArgoCDProps = {}) {
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

        super(scope, id, {
            releaseName: "argocd",
            chart: "argo/argo-cd",
            values: props.ha ? ha_values : {},
        });
    }
}