# Local audio and video

Run from this directory:

```sh
go run .
```

Click **Load clips**, then play or seek using the controls. Click **Clear clips**
to cancel pending loads and release their blob URLs. The example uses the bundled
runtime and embedded MP3/MP4 samples from the original #4782 example.

`Media.SetSource` receives each complete clip over the registered `media` stream
with a 32 MiB default limit. Go sends 64 KiB frames over the existing Wails
transport. Both Go and the frontend enforce their byte limits. It opens no network listener. For an npm
frontend, import `Media` from `@wailsio/runtime` instead.

Linux needs GStreamer codecs for MP3 and H.264/AAC. See the
[media playback guide](https://v3.wails.io/guides/linux-media/) for limits,
component cleanup and troubleshooting.
