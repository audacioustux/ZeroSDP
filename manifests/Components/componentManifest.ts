import * as k8s from "@pulumi/kubernetes";
import { Input, Inputs, ComponentResource, ComponentResourceOptions } from "@pulumi/pulumi";
import { getStack } from "@pulumi/pulumi";

export class ComponentManifest extends ComponentResource {
    directory: Input<string>;

    constructor(type: string, name: string, args?: Inputs, opts?: ComponentResourceOptions) {
        const renderYamlToDirectory = `./rendered/${getStack()}/${name}`;
        const provider = new k8s.Provider("provider", {
            kubeconfig: "",
            renderYamlToDirectory,
        })

        super(`components:${type}`, name, args, { provider, ...opts });

        this.directory = renderYamlToDirectory;
    }
}
