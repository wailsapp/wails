This template uses a work around as the default template does not compile due to
this issue: https://github.com/vuejs/core/issues/1228

In `tsconfig.json`, `isolatedModules` is set to `false` rather than `true` to
work around the issue.
