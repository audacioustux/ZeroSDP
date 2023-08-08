import { Construct } from 'constructs';
import { App, Chart, ChartProps, Helm } from 'cdk8s';
import { KubeDeployment, KubeService, IntOrString, KubeNamespace } from '../imports/k8s.js';
import { Application } from "../imports/argoproj.io.js"

export class ArgoCD extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    const namespace = new KubeNamespace(this, 'ns', {
      metadata: {
        name: 'argocd',
      },
    });

    const helm = new Helm(this, 'helm', {
      releaseName: "argocd",
      chart: "argo/argo-cd",
      helmFlags: [
        "--set", "installCRDs=true",
        "--namespace", namespace.name
      ],
      values: {
        "redis-ha": {
          // enabled: true
        },
        controller: {
          replicas: 1
        },
        server: {
          // autoscaling: {
          //   enabled: true,
          //   minReplicas: 2
          // }
        },
        repoServer: {
          // autoscaling: {
          //   enabled: true,
          //   minReplicas: 2
          // },
        },
        applicationSet: {
          // replicaCount: 2
        }
      },
    });

    const app = new Application(this, 'app', {
      metadata: {
        name: helm.releaseName,
        namespace: namespace.name,
      },
      spec: {
        destination: {
          namespace: namespace.name,
          server: 'https://kubernetes.default.svc',
        },
        project: 'default',
        source: {
          path: 'manifests/dist',
          directory: {
            include: "argocd.k8s.yaml"
          },
          repoUrl: "https://github.com/audacioustux/sdp.git",
          targetRevision: 'HEAD',
        },
        syncPolicy: {
          automated: {
            prune: true,
            selfHeal: true,
          },
        },
      },
    });
  }
}

const app = new App();
new ArgoCD(app, 'argocd');
app.synth();
