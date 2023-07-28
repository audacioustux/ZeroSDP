import { ArgoCD } from "./Components/argoCD.js";

export const argocd = new ArgoCD("argo-cd", { version: "2.7.7" }).directory;
