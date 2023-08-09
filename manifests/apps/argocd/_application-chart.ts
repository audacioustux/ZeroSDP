import { Construct } from 'constructs';
import { Chart, ChartProps } from 'cdk8s';
import { Application, ApplicationProps } from "../../imports/argocd-argoproj.io.js"
import { PartialDeep } from 'type-fest';

export class ApplicationChart extends Chart {
    constructor(scope: Construct, id: string, applicationProps: ApplicationProps, props: ChartProps) {
        super(scope, id, props);

        const defaults: PartialDeep<ApplicationProps> = {
            metadata: {
                namespace: "argocd",
                finalizers: ["resources-finalizer.argocd.argoproj.io"]
            },
            spec: {
                destination: {
                    server: "https://kubernetes.default.svc",
                }
            }
        }

        const app = new Application(this, 'app', Object.assign(defaults, applicationProps));
    };
}