import { Chart, ChartProps } from "cdk8s"
import { Construct } from "constructs"
import { AppProject, Application } from "../imports/argocd-argoproj.io.js";
import { ArgoCD } from "./argo-cd.chart.js";

export class SDP extends Chart {
    constructor(scope: Construct, props: ChartProps = {}) {
        super(scope, "sdp", props)

        const argoCD = new ArgoCD(this)

        const name = "platform"
        const project = new AppProject(this, name, {
            metadata: { name, namespace: argoCD.namespace },
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


        const appName = (name: string) => `${project.name}-${name}`

        {
            const name = appName("argocd")
            new Application(this, name, {
                metadata: { name, namespace: argoCD.namespace },
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

        {
            // argo-workflows
            const name = appName("argo-workflows")
            new Application(this, name, {
                metadata: { name, namespace: argoCD.namespace },
                spec: {
                    project: project.name,
                    destination: {
                        namespace: "argo",
                        server: "https://kubernetes.default.svc",
                    },
                    source: {
                        repoUrl: "https://argoproj.github.io/argo-helm",
                        chart: "argo-workflows",
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
}