import { ArgoCD } from "./Components/ArgoCD";

export const argocd = new ArgoCD("argo-cd", { version: "2.7.7" }).path;
