# README

## About

This is a basic Svelte template, using rollup to bundle the assets into a single JS file.
Rollup is configured to do the following:

- Convert imported images to base64 strings
- Convert `url()` in `@font-face` declarations to base64 strings
- Bundle all css into the JS bundle
- Copy `index.html` from `frontend/src/` to `frontend/dist/`

Clicking the button will call the backend.

## Live Development

To run in live development mode, run `wails dev` in the project directory. The frontend dev server will run
on http://localhost:34115. Open this in your browser to connect to your application.

## Building

For a production build, use `wails build`.
