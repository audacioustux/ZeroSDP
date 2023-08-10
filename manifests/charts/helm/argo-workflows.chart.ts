import { Construct } from 'constructs';
import { HelmChart } from './_helm-chart.js';

export interface ArgoWorkflowsProps {
}

export class ArgoWorkflows extends HelmChart {
    constructor(scope: Construct, id: string, props: ArgoWorkflowsProps = {}) {
        super(scope, id, {
            releaseName: "argo-workflows",
            chart: "argo/argo-workflows",
        });
    }
}