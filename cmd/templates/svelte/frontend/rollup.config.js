import svelte from 'rollup-plugin-svelte';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import { terser } from 'rollup-plugin-terser';
import image from '@rollup/plugin-image';
import babel from 'rollup-plugin-babel';
import polyfill from 'rollup-plugin-polyfill';

const production = !process.env.ROLLUP_WATCH;

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

export default {
	input: 'src/main.js',
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: 'public/build/bundle.js'
	},
	plugins: [
		image(),
		svelte({
			// enable run-time checks when not in production
			dev: !production,
			// we'll extract any component CSS out into
			// a separate file - better for performance
			css: css => {
				css.write('bundle.css');
			}
		}),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins. In
		// some cases you'll need additional configuration -
		// consult the documentation for details:
		// https://github.com/rollup/plugins/tree/master/packages/commonjs
		resolve({
			browser: true,
			dedupe: ['svelte', 'svelte/transition', 'svelte/internal']
		}),
		commonjs(),

		// In dev mode, call `npm run start` once
		// the bundle has been generated
		!production && serve(),

		// Watch the `public` directory and refresh the
		// browser on changes when not in production
		!production && livereload('public'),

		// Credit: https://blog.az.sg/posts/svelte-and-ie11/
		babel({
			extensions: [ '.js', '.jsx', '.es6', '.es', '.mjs', '.svelte', '.html' ],
			runtimeHelpers: true,
			exclude: [ 'node_modules/@babel/**', 'node_modules/core-js/**' ],
			presets: [
			  [
				'@babel/preset-env',
				{
				  targets: '> 0.25%, not dead, IE 11',
				  modules: false,
				  useBuiltIns: 'usage',
				  forceAllTransforms: true,
				  corejs: 3,
				},

			  ]
			],
			plugins: [
			  '@babel/plugin-syntax-dynamic-import',
			  [
				'@babel/plugin-transform-runtime',
				{
				  useESModules: true
				}
			  ]
			]
		  }),
		  polyfill(['@webcomponents/webcomponentsjs']),

		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser()
	],
	watch: {
		clearScreen: false
	}
};
