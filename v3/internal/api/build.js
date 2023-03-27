
// jshint esversion: 9
const { build } = require("esbuild");
const { dependencies, peerDependencies } = require('./package.json');

let deps = dependencies || {};
let peerDeps = peerDependencies || {};

const sharedConfig = {
    entryPoints: ["src/index.ts"],
    bundle: true,
    minify: true,
    external: Object.keys(deps).concat(Object.keys(peerDeps)),
};

build({
    ...sharedConfig,
    platform: 'node', // for CJS
    outfile: "dist/index.js",
});

build({
    ...sharedConfig,
    outfile: "dist/index.esm.js",
    platform: 'neutral', // for ESM
    format: "esm",
});

