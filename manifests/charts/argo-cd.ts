import { Chart, ChartProps, Helm, HelmProps } from "cdk8s";
import { Construct } from "constructs";
import { Namespace } from "cdk8s-plus-27";

export const DEFAULT_VALUES = {} as const;
export const DEFAULT_HA_VALUES = {
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
} as const;

export type ArgoCDProps = Pick<HelmProps, 'namespace' | 'values' | 'releaseName'>;
export class ArgoCD extends Chart {
    constructor(scope: Construct, props: ArgoCDProps = {}, chartProps: ChartProps = {}) {
        const { releaseName = "argocd", namespace = "argocd", values = DEFAULT_VALUES } = props;
        super(scope, releaseName, { namespace, ...chartProps });

        new Namespace(this, 'ns', { metadata: { name: namespace } });
        new Helm(this, "helm", {
            repo: 'https://argoproj.github.io/argo-helm',
            chart: 'argo-cd',
            releaseName,
            namespace,
            values,
            helmFlags: [
                "--skip-tests",
            ]
        });
    }
}