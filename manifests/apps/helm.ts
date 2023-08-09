import { App } from 'cdk8s';
import { ArgoCD } from './helm/argo-cd.js';
import { ArgoWorkflow } from './helm/argo-workflows.js';

export const app = new App({ outdir: "helm" })
// argo-cd
new ArgoCD(app, 'argo-cd');
new ArgoCD(app, 'argo-cd-ha', { ha: true });
// argo-workflows
new ArgoWorkflow(app, 'argo-workflows');
