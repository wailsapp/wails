/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/**
 * True when running inside a browser/webview with a DOM available.
 * False under server-side rendering (e.g. `next build` prerendering),
 * where application code may import the runtime module even though no
 * Wails APIs can actually be used (#4679). Modules must not touch
 * `window`/`document` at import time except behind this guard.
 */
export const hasDOM = typeof window !== "undefined" && typeof document !== "undefined";
