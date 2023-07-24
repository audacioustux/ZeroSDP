import { Namespace } from "@pulumi/kubernetes/core/v1";

export function useNamespace(namespace: Namespace) {
    return (obj: any) => {
        if (obj.metadata) {
            obj.metadata.namespace = namespace;
        }
    };
}
