import { App } from 'cdk8s';
import { ArgoCD } from './charts/argo-cd.js';
import { ArgoWorkflow } from './charts/argo-workflow.js';

const app = new App()
new ArgoCD(app, 'argo-cd', { ha: false });
new ArgoWorkflow(app, 'argo-workflow');
app.synth();
