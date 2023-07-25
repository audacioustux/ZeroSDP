import * as k8s from "@pulumi/kubernetes";
import { ComponentResourceOptions } from "@pulumi/pulumi";
import { ComponentManifest } from "./ComponentManifest";

export interface ArgoCDArgs {
    version: string;
}

export class ArgoCD extends ComponentManifest {
    constructor(name: string, args: ArgoCDArgs, opts?: ComponentResourceOptions) {
        super("ArgoCD", name, args, opts);

        const { version } = args;

        const manifest = new k8s.yaml.ConfigFile(`${name}-bootstrap`, {
            file: `https://raw.githubusercontent.com/argoproj/argo-cd/v${version}/manifests/install.yaml`
        }, { parent: this });

        this.registerOutputs();
    }
}