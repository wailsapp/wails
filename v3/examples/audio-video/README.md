# Audio/Video Example

This example demonstrates HTML5 audio and video playback using the `@wailsio/runtime` npm module.

## Linux Notes

On Linux, WebKitGTK uses GStreamer for media playback. GStreamer doesn't have a URI handler for the `wails://` protocol, which means media files served from the bundled assets won't play directly.

Wails automatically works around this limitation by intercepting media elements and converting their sources to blob URLs. This happens transparently - you don't need to change your code.

See the [Linux-specific documentation](https://wails.io/docs/guides/linux-media) for details on:
- How the GStreamer workaround works
- How to disable it if needed (`DisableGStreamerFix`)
- How to enable caching for better performance (`EnableGStreamerCaching`)

## Building

```bash
cd frontend
npm install
npm run build
cd ..
go build
./audio-video
```

## Development

For development with hot-reload:

```bash
# Terminal 1: Run Vite dev server
cd frontend
npm run dev

# Terminal 2: Run Go app with dev server URL
FRONTEND_DEVSERVER_URL=http://localhost:5173 go run .
```
