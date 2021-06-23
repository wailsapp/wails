import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import { terser } from 'rollup-plugin-terser';
import postcss from 'rollup-plugin-postcss';
import image from '@rollup/plugin-image';
import url from '@rollup/plugin-url';
import copy from 'rollup-plugin-copy';

const production = !process.env.ROLLUP_WATCH;

export default {
    input: 'src/main.js',
    output: {
        sourcemap: true,
        format: 'iife',
        name: 'app',
        file: 'dist/main.js'
    },
    onwarn: handleRollupWarning,
    plugins: [

        image(),

        copy({
            targets: [
                { src: 'src/index.html', dest: 'dist' },
                { src: 'src/main.css', dest: 'dist' },
            ]
        }),

        // Embed binary files
        url({
            include: ['**/*.woff', '**/*.woff2'],
            limit: Infinity,
        }),

        // If you have external dependencies installed from
        // npm, you'll most likely need these plugins. In
        // some cases you'll need additional configuration -
        // consult the documentation for details:
        // https://github.com/rollup/plugins/tree/master/packages/commonjs
        resolve({
            browser: true,
        }),
        commonjs(),

        // PostCSS preprocessing
        postcss({
            extensions: ['.css', '.scss'],
            extract: true,
            minimize: false,
            use: [
                ['sass', {
                    includePaths: [
                        './src',
                        './node_modules'
                    ]
                }]
            ],
        }),

        // In dev mode, call `npm run start` once
        // the bundle has been generated
        !production && serve(),

        // Watch the `public` directory and refresh the
        // browser on changes when not in production
        !production && livereload('dist'),

        // If we're building for production (npm run build
        // instead of npm run dev), minify
        production && terser()
    ],
    watch: {
        clearScreen: false
    }
};

function handleRollupWarning(warning) {
    console.error('ERROR: ' + warning.toString());
}

function serve() {
    let server;

    function toExit() {
        if (server) server.kill(0);
    }

    return {
        writeBundle() {
            if (server) return;
            server = require('child_process').spawn('npm', ['run', 'start', '--', '--dev'], {
                stdio: ['ignore', 'inherit', 'inherit'],
                shell: true
            });

            process.on('SIGTERM', toExit);
            process.on('exit', toExit);
        }
    };
}