/* jshint esversion: 8 */
const esbuild = require("esbuild");
const sveltePlugin = require("esbuild-svelte");

esbuild
    .build({
        entryPoints: ["main.js"],
        bundle: true,
        minify: true,
        outfile: "../sse.js",
        plugins: [sveltePlugin({compileOptions: {css: true}})],
        logLevel: "info",
        sourcemap: "inline",
    })
    .catch(() => process.exit(1));
