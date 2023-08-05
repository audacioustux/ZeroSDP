import { Construct } from 'constructs';
import { App, Chart, ChartProps } from 'cdk8s';
// import { KubeDeployment, KubeService, IntOrString } from './imports/k8s.js';

export class MyChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

  }
}

const app = new App();
new MyChart(app, 'manifests');
app.synth();
