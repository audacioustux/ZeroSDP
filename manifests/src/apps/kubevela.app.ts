import { Application } from "../../imports/argocd-argoproj.io.js";
import { ArgoCD } from "../charts/argo-cd.chart.js";
import { Project } from "../projects/project.js";

export class Kubevela {
    name = "kubevela"

    constructor(scope: ArgoCD, project: Project) {
        new Application(scope, this.name, {
            metadata: { name: `${project.name}-${this.name}`, namespace: scope.namespace, finalizers: ["resources-finalizer.argocd.argoproj.io"] },
            spec: {
                project: project.name,
                destination: {
                    namespace: "vela-system",
                    server: "https://kubernetes.default.svc",
                },
                source: {
                    repoUrl: "https://kubevela.github.io/charts",
                    chart: "vela-core",
                    targetRevision: "1.9.*"
                },
                syncPolicy: {
                    automated: {
                        prune: true,
                        selfHeal: true,
                    },
                    syncOptions: [
                        "CreateNamespace=true"
                    ]
                }
            }
        })
    }
}