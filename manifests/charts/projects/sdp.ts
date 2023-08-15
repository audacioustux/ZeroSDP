import { Chart } from "cdk8s"
import { Construct } from "constructs"
import { Application, AppProject } from '../../imports/argocd-argoproj.io.js';

export interface SDPProps { }
export class SDP extends Chart {
    constructor(scope: Construct) {
        super(scope, "sdp")

        const project = new AppProject(this, "project", {
            metadata: {
                name: "sdp",
                namespace: "argocd",
            },
            spec: {
                description: "SDP",
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
                ]
            }
        })

        new Application(this, "app", {
            metadata: {
                name: "argo-cd",
                namespace: "argocd",
            },
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
        });
    }
}