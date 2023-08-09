import { Construct } from 'constructs';
import { ChartProps } from 'cdk8s';
import { HelmChart } from './_helm-chart.js';

export interface ArgoWorkflowProps extends ChartProps {
}

export class ArgoWorkflow extends HelmChart {
    constructor(scope: Construct, id: string, props: ArgoWorkflowProps = {}) {
        super(scope, id, {
            releaseName: "argo-workflows",
            chart: "argo/argo-workflows",
        });
    }
}