import { App, AppProps } from 'cdk8s';
import { SDP } from "./src/index.chart.js"

const app = new App();

new SDP(app)

app.synth();
