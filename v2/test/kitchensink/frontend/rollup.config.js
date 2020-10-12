import svelte from 'rollup-plugin-svelte';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import { terser } from 'rollup-plugin-terser';
import postcss from 'rollup-plugin-postcss';
import autoPreprocess from 'svelte-preprocess';
import { string } from "rollup-plugin-string";
import url from '@rollup/plugin-url';

const production = !process.env.ROLLUP_WATCH;

export default {
	input: 'src/main.js',
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: 'public/bundle.js'
	},
	onwarn: handleRollupWarning,
	plugins: [
		
		// Embed binary files
		url({	
			include: ['**/*.woff', '**/*.woff2'],
			limit: Infinity,			
		}),

		// Embed text files
		string({
			include: ["**/*.jsx","**/*.go"],
		}),
		
		svelte({
			preprocess: autoPreprocess(),
			// enable run-time checks when not in production
			dev: !production,
			// we'll extract any component CSS out into
			// a separate file - better for performance
			css: css => {
				css.write('public/bundle.css');
			},
			emitCss: true,
		}),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins. In
		// some cases you'll need additional configuration -
		// consult the documentation for details:
		// https://github.com/rollup/plugins/tree/master/packages/commonjs
		resolve({
			browser: true,
			dedupe: ['svelte']
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
						'./src/theme',
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
		!production && livereload('public'),

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
	process.exit(1);
}

function serve() {
	let started = false;

	return {
		writeBundle() {
			if (!started) {
				started = true;

				require('child_process').spawn('npm', ['run', 'start', '--', '--dev'], {
					stdio: ['ignore', 'inherit', 'inherit'],
					shell: true
				});
			}
		}
	};
}
