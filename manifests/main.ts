import { App } from 'cdk8s';
import { ArgoCD, DEFAULT_HA_VALUES } from "./charts/argo-cd.js"
import { SDP } from "./charts/projects/sdp.js"

const app = new App();

const argocd = new ArgoCD(app);
const sdp = new SDP(app);

sdp.addDependency(argocd);

app.synth();

