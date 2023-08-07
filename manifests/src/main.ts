import { Construct } from 'constructs';
import { App, Chart, ChartProps, Helm } from 'cdk8s';
import { KubeDeployment, KubeService, IntOrString, KubeNamespace } from '../imports/k8s.js';

export class ArgoCD extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    const namespace = new KubeNamespace(this, 'argocd-ns', {
      metadata: {
        name: 'argocd',
      },
    });

    const helm = new Helm(this, 'argocd', {
      chart: "argo/argo-cd",
      helmFlags: [
        "--set", "installCRDs=true",
        "--namespace", namespace.name
      ],
      values: {
        "redis-ha": {
          enabled: true
        },
        controller: {
          replicas: 1
        },
        server: {
          autoscaling: {
            enabled: true,
            minReplicas: 2
          }
        },
        repoServer: {
          autoscaling: {
            enabled: true,
            minReplicas: 2
          },
        },
        applicationSet: {
          replicaCount: 2
        }
      },
    });
  }
}

const app = new App();
new ArgoCD(app, 'argocd');
app.synth();
