import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const config = new pulumi.Config();

const argo_cd = new kubernetes.helm.v3.Release("argo-cd", {
    chart: "argo-cd",
    repositoryOpts: {
        repo: "https://argoproj.github.io/argo-helm",
    },
    namespace: "argo-cd",
    createNamespace: true,
});
