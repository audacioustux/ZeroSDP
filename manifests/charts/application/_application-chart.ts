import { Construct } from 'constructs';
import { Chart, ChartProps } from 'cdk8s';
import { Application, ApplicationProps } from "../../imports/argocd-argoproj.io.js"

export class ApplicationChart extends Chart {
    constructor(scope: Construct, id: string, applicationProps: ApplicationProps, chartProps: ChartProps) {
        super(scope, id, chartProps);

        const app = new Application(this, 'app', applicationProps);
    };
}