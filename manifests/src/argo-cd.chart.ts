import { Chart, ChartProps, Helm, HelmProps } from "cdk8s";
import { Construct } from "constructs";
import { Namespace } from "cdk8s-plus-27";

export interface ArgoCDProps extends ChartProps {
    ha?: false | {
        autoscaling?: boolean;
    }
}
export class ArgoCD extends Chart {
    constructor(scope: Construct, id: string = "argocd", chartProps: ArgoCDProps = { namespace: "argocd", ha: false }) {
        const { namespace, ha } = chartProps;
        super(scope, id, chartProps);

        let values = ha ? {
            "redis-ha": {
                enabled: true
            },
            controller: {
                replicas: 1
            },
            applicationSet: {
                replicaCount: 2
            },
            server: {
                ...(ha.autoscaling ? {
                    autoscaling: {
                        enabled: true,
                        minReplicas: 2
                    },
                } : { replicas: 2 })
            },
            repoServer: {
                ...(ha.autoscaling ? {
                    autoscaling: {
                        enabled: true,
                        minReplicas: 2
                    },
                } : { replicas: 2 })
            },
        } : {}

        new Namespace(this, 'ns', { metadata: { name: namespace } });
        new Helm(this, "helm", {
            repo: 'https://argoproj.github.io/argo-helm',
            chart: 'argo-cd',
            releaseName: id,
            namespace,
            values
        });
    }
}