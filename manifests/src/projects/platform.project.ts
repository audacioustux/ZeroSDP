import { AppProject } from "../../imports/argocd-argoproj.io.js";
import { ArgoCD } from "../charts/argo-cd.chart.js";
import { Root } from "../apps/root.app.js";
import { Kubevela } from "../apps/kubevela.app.js";
import { Project } from "./project.js";

export class Platform implements Project {
    name = "platfrom"

    constructor(scope: ArgoCD) {
        const { namespace } = scope

        new AppProject(scope, this.name, {
            metadata: { name: this.name, namespace },
            spec: {
                sourceRepos: ["*"],
                destinations: [
                    {
                        namespace: "*",
                        server: "https://kubernetes.default.svc",
                    }
                ],
                clusterResourceWhitelist: [
                    {
                        group: "*",
                        kind: "*",
                    }
                ],
                namespaceResourceWhitelist: [
                    {
                        group: "*",
                        kind: "*",
                    }
                ],
            }
        })

        new Root(scope, this)
        new Kubevela(scope, this)
    }
}