import { App } from 'cdk8s';
import { ArgoCD } from './charts/argo-cd.js';

const app = new App()
new ArgoCD(app, 'argocd', { ha: false });
app.synth();
