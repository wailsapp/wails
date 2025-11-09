# Unreleased Changes

Improves re-build performance by optimising dependency change detection for front end builds made with `wails3 init`.

## Changed
- `common:install:frontend:deps` go-task now uses node_modules/.bin/* for change control detection

## Fixed
- Fixes an issue with the front end directory if a local npm package is installed `npm i ~/my-lib`


