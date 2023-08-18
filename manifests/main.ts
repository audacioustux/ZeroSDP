import { App } from 'cdk8s';
import { SDP } from "./src/sdp.chart.js"

const app = new App();

new SDP(app)

app.synth();
