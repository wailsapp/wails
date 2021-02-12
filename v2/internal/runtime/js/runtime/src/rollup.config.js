import resolve from '@rollup/plugin-node-resolve';
// import commonjs from '@rollup/plugin-commonjs';
import svelte from 'rollup-plugin-svelte';
import { terser } from "rollup-plugin-terser";

export default [
    // browser-friendly UMD build
    {
        input: 'main.js',
        output: {
            name: 'bridge',
            file: '../bridge.js',
            format: 'umd',
            exports: "named"
        },
        plugins: [
            svelte({
                // Optionally, preprocess components with svelte.preprocess:
                // https://svelte.dev/docs#svelte_preprocess
                // preprocess: {
                //     style: ({content}) => {
                //         return transformStyles(content);
                //     }
                // },

                // Emit CSS as "files" for other plugins to process. default is true
                emitCss: false,

            }),
            resolve({browser: true}), // so Rollup can find `ms`
            // commonjs() // so Rollup can convert `ms` to an ES module
            terser(),
        ]
    },


];