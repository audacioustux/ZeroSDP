import { Construct } from 'constructs';
import { Chart, ChartProps, Helm, HelmProps } from 'cdk8s';

export class HelmChart extends Chart {
    constructor(scope: Construct, id: string, helmProps: HelmProps, chartProps: ChartProps = {}) {
        super(scope, id, Object.assign({ disableResourceNameHashes: true }, chartProps));

        new Helm(this, 'helm', helmProps)
    };
}