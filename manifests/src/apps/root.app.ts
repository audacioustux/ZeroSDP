import { Project } from "src/projects/project.js";
import { Application } from "../../imports/argocd-argoproj.io.js";
import { ArgoCD } from "../charts/argo-cd.chart.js";
import { App } from "./app.js";

export class Root implements App {
    name = "root"

    constructor(scope: ArgoCD, project: Project) {
        new Application(scope, this.name, {
            metadata: { name: `${project.name}-${this.name}`, namespace: scope.namespace },
            spec: {
                project: project.name,
                destination: {
                    server: "https://kubernetes.default.svc",
                },
                source: {
                    repoUrl: "https://github.com/audacioustux/sdp.git",
                    path: "manifests/dist"
                },
                syncPolicy: {
                    automated: {
                        prune: true,
                        selfHeal: true,
                    }
                }
            }
        })
    }
}