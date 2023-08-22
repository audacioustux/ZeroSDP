import { App, AppProps } from "cdk8s"
import { ArgoCD } from "./charts/argo-cd.chart.js";
import { Platform } from "./projects/platform.project.js"

export class SDP extends App {
    constructor(props: AppProps = {}) {
        super(props)

        new Platform(new ArgoCD(this))
    }
}