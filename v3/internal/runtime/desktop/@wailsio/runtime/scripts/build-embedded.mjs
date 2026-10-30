import { fileURLToPath } from 'node:url';
import { resolve } from 'node:path';
import { build, normalizePath } from 'vite';

const mode = process.argv[2];
if (mode !== 'debug' && mode !== 'production') {
    throw new Error('Usage: node scripts/build-embedded.mjs debug|production [output-directory]');
}
const debug = mode === 'debug';
const root = fileURLToPath(new URL('../', import.meta.url));
const entry = normalizePath(fileURLToPath(new URL('../src/index.ts', import.meta.url)));
const outDir = process.argv[3]
    ? resolve(process.argv[3])
    : fileURLToPath(new URL('../../../../../assetserver/bundledassets/', import.meta.url));

await build({
    root,
    configFile: false,
    envDir: false,
    publicDir: false,
    define: { DEBUG: String(debug) },
    plugins: [{
        name: 'wails-bootstrap',
        transform(code, id) {
            if (id !== entry) return;
            // The cycle runs bootstrap before index.ts announces runtime readiness.
            // Append the import to leave every original source-map position intact.
            return { code: `${code}\nimport '../../../compiled/main.js';`, map: null };
        },
    }],
    build: {
        target: 'safari11',
        lib: { entry, formats: ['es'], fileName: () => debug ? 'runtime.debug.js' : 'runtime.js' },
        outDir,
        emptyOutDir: false,
        sourcemap: debug ? 'inline' : false,
        minify: false,
        rolldownOptions: {
            cwd: root,
            treeshake: { annotations: false, moduleSideEffects: true },
            output: {
                comments: { jsdoc: false },
                // Bootstrap imports the entry namespace before the entry body runs.
                strictExecutionOrder: true,
                minify: debug ? false : {
                    compress: { target: 'safari11', dropConsole: true, dropDebugger: false },
                    mangle: true,
                    codegen: { removeWhitespace: true },
                },
            },
        },
    },
});
