---
title: "播放本地音频和视频"
description: "在 Linux 上通过有大小限制的 blob URL 播放随应用分发的媒体，并在结束后释放。"
slug: "guides/linux-media"
sourcePath: "guides/linux-media.md"
---

使用 `Media.SetSource` 在 Wails v3 应用中播放简短的本地片段。在 Linux 上，WebKitGTK 将媒体播放交给 GStreamer，但后者无法直接加载 `wails://` URL。此辅助函数通过 Wails 流接收片段，并将 blob URL 分配给播放器。

桌面传输使用现有的 Wails 资源传输机制，不会打开监听套接字。该 API 支持音频和视频元素。普通 HTTP/HTTPS 媒体可直接使用播放器的 `src`。

## 注册媒体文件

公开一个包含前端可播放片段的文件系统，然后在命名流上注册媒体处理器：

```go
import (
    "embed"
    "io/fs"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/media"
)

//go:embed clips
var clips embed.FS

func registerMedia(app *application.App) {
    mediaFiles, err := fs.Sub(clips, "clips")
    if err != nil {
        log.Fatal(err)
    }
    handler, err := media.NewHandler(mediaFiles, 32 << 20) // 32 MiB per file
    if err != nil {
        log.Fatal(err)
    }
    app.HandleStream("media", handler)
}
```

创建应用后、调用 `app.Run()` 之前，调用 `registerMedia(app)`。

对于磁盘文件，使用 `os.OpenRoot(directory)`，并将 `root.FS()` 传给 `media.NewHandler`。保持根目录打开，直到 `app.Run()` 返回后再关闭。选择一个只包含允许前端读取的文件的目录；即使符号链接指向目录外，`os.Root` 也会限制访问。

## 加载片段

创建播放器：

```html
<video id="player" controls></video>
```

对于使用 npm 运行时的前端：

```javascript
import { Media } from '@wailsio/runtime';

const player = document.getElementById('player');

try {
    await Media.SetSource(player, 'media', 'welcome.mp4');
} catch (error) {
    if (error.name !== 'AbortError') {
        console.error('Could not load the clip:', error);
    }
}
```

使用内置运行时的应用请将导入改为：

```javascript
import { Media } from '/wails/runtime.js';
```

第二个参数是已注册的流名称。第三个参数是相对于文件系统根目录、以斜杠分隔的文件路径，例如 `welcome.mp4` 或 `tutorials/intro.mp4`。它不是 URL，也不是操作系统路径。

分配源后，Promise 即完成，随后播放器进行解码。处理播放器的 `error` 事件以检测不支持的编解码器；如果需要自行开始播放，请在用户交互中调用 `player.play()`。

## 替换或释放片段

再次调用 `Media.SetSource` 可更换片段。它会取消该播放器之前尚未完成的加载，避免缓慢的响应覆盖最新选择。在新片段成功加载前，原片段仍然可用，随后原来的 blob URL 会被撤销。

关闭播放器或卸载组件时，调用 `Media.ClearSource(player)`：

```javascript
Media.ClearSource(player);
```

这会取消尚未完成的加载、重置播放器并释放 blob URL。请在从文档移除元素之前调用。框架组件应在卸载或销毁钩子中调用它。仅移除元素不会释放 blob URL。请始终使用这些辅助函数管理播放器的源。

对于 `<video><source ...></video>` 标记，将所选文件名传给 `Media.SetSource(video, "media", name)`。该辅助函数设置父播放器的 `src`，其优先级高于子 `<source>` 元素。清除此属性后，浏览器会重新考虑这些子元素；通过本 API 管理全部源时，请使用空播放器。

未连接到文档的音频对象也可以使用：

```javascript
const sound = new Audio();
await Media.SetSource(sound, 'media', 'notification.mp3');
sound.addEventListener('ended', () => Media.ClearSource(sound), { once: true });
// Call sound.play() from an appropriate user interaction.
// Also clear it if playback is cancelled or the owning component is disposed.
```

## 限制下载并取消加载

默认限制为**每个源 32 MiB**。你可以在 Go 中配置的限制内，选择更小或更大的正整数字节限制：

```javascript
const controller = new AbortController();
const loading = Media.SetSource(player, 'media', 'welcome.mp4', {
    maxBytes: 8 * 1024 * 1024,
    signal: controller.signal,
});

// Call controller.abort() to cancel this load.
await loading;
```

文件过大会以 `RangeError` 拒绝。Go 处理器在读取内容前检查文件大小，最多传输其配置限制与前端限制中较小的值。前端也会检查接收大小并拒绝不完整传输。取消时会以 `AbortError` 或传给 `AbortController.abort(reason)` 的原因拒绝。

文件以 64 KiB 的帧传输。Wails 桌面流传输将轮询响应限制为 1 MiB，因此 WebView2 缓冲完整响应时，不会在单个响应中缓冲整个媒体文件。流队列自身也具有有界背压机制。这些传输限制不会改变播放器辅助函数使用完整 blob 的行为。

**整个文件下载完毕后才会播放。** 字节限制针对单个文件，并非应用总内存限制。多个播放器、替换过程中的原片段、blob 构建及解码后的媒体都可能使用额外内存。辅助函数一经调用就会传输，不受播放器 `preload` 设置影响。请在用户选择加载片段时调用。

对于大型本地文件，此辅助函数不是流式播放方案。提高限制也会增加内存使用。已通过 HTTP/HTTPS 托管的媒体应使用原生媒体加载，使浏览器能通过范围请求进行流式播放和定位。

## 排查 Linux 播放问题

- 如果直接播放本地文件时报 **No URI handler implemented for "wails"**，请使用 `Media.SetSource` 加载片段。默认 GTK4 栈和旧版 `-tags gtk3` 栈均受影响。
- 如果加载成功但解码失败，请检查目标系统安装的 GStreamer 编解码器。MP4 通常需要 H.264 视频和 AAC 音频支持；MP3 需要 MP3 解码器。请在支持的发行版上测试所分发的格式。
- 如果传输失败，请检查注册的流名称、相对文件名和文件系统权限。传输期间修改文件可能导致传输不完整；请在写入结束后重试。
- 如果应用设置了内容安全策略，请在 `connect-src` 中允许 Wails 资源源站，在 `media-src` 中允许 `blob:`。对于仅本地使用的策略，可使用 `connect-src 'self'; media-src 'self' blob:`。保留其他指令。
- 如果文件超出限制，请选择更短或更小的片段，或设置符合应用内存预算的明确限制。

使用 `go run .` 运行 [audio-video 示例](https://github.com/wailsapp/wails/tree/master/v3/examples/audio-video)，在本机检查随附的 MP3 和 MP4 样例。
