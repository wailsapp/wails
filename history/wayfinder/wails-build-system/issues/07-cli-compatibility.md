# CLI compatibility and migration cutover

Type: grilling
Status: open
Blocked by: 01

## Question

What should `wails3 init`, `build`, `dev`, `package`, `sign`, `eject`,
`migrate`, `config check`, and `config show` do in manifest projects and
pre-migration Taskfile projects? Define detection, error messages, release
staging, precedence and conflict behavior when `wails.toml` and Taskfiles
coexist, how migration completion or an explicit cutover marker affects that
precedence, the exact legacy fallback condition, and the point at which
generated Taskfiles stop being supported.

## Comments
