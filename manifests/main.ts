import { App } from 'cdk8s';
import { ArgoCD } from './apps/helm/argo-cd.js';
import { ArgoWorkflow } from './apps/argocd/argo-workflow.js';
import { app as helm } from "./apps/helm.js";

helm.synth();