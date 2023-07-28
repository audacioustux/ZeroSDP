export function useNamespace(namespace) {
    return (obj) => {
        if (obj.metadata) {
            obj.metadata.namespace = namespace;
        }
    };
}
