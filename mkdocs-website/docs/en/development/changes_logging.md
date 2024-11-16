### Logging

Logging in v2 was confusing as both application logs and system (internal) logs
were using the same logger. We have simplified this as follows:

- Internal logs are now handled using the standard Go `slog` logger. This is
  configured using the `logger` option in the application options. By default,
  this uses the [tint](https://github.com/lmittmann/tint) logger.
- Application logs can now be achieved through the new `log` plugin which
  utilises `slog` under the hood. This plugin provides a simple API for logging
  to the console. It is available in both Go and JS.
