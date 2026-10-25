---
title: "播放本機音訊與影片"
description: "在 Linux 上使用有大小限制的 blob URL 播放隨附媒體，並在結束後釋放。"
slug: "guides/linux-media"
sourcePath: "guides/linux-media.md"
---

使用 `Media.SetSource` 在 Wails v3 應用程式中播放簡短的本機片段。在 Linux 上，WebKitGTK 將媒體播放交給 GStreamer，但後者無法直接載入 `wails://` URL。此輔助函式透過 Wails 串流接收片段，並將 blob URL 指派給播放器。

桌面傳輸使用現有的 Wails 資源傳輸機制，不會開啟監聽通訊端。此 API 支援音訊與影片元素。一般 HTTP/HTTPS 媒體可直接使用播放器的 `src`。

## 註冊媒體檔案

公開包含前端可播放片段的檔案系統，再將媒體處理常式註冊到具名串流：

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

建立應用程式後、呼叫 `app.Run()` 之前，請呼叫 `registerMedia(app)`。

對於磁碟上的檔案，使用 `os.OpenRoot(directory)`，並將 `root.FS()` 傳給 `media.NewHandler`。保持根目錄開啟，直到 `app.Run()` 返回後再關閉。請選擇只包含允許前端讀取的檔案的目錄；即使符號連結指向目錄外，`os.Root` 也會限制存取。

## 載入片段

建立播放器：

```html
<video id="player" controls></video>
```

對於使用 npm 執行階段的前端：

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

使用內建執行階段的應用程式，請將匯入改為：

```javascript
import { Media } from '/wails/runtime.js';
```

第二個引數是已註冊的串流名稱。第三個是相對於檔案系統根目錄、以斜線分隔的檔案路徑，例如 `welcome.mp4` 或 `tutorials/intro.mp4`。它不是 URL，也不是作業系統路徑。

指定來源後，Promise 即完成，接著由播放器解碼。請處理播放器的 `error` 事件，以偵測不支援的編解碼器；若需要自行開始播放，請在使用者互動中呼叫 `player.play()`。

## 取代或釋放片段

再次呼叫 `Media.SetSource` 即可更換片段。它會取消該播放器先前尚未完成的載入，避免緩慢的回應覆寫最新選擇。在替代片段成功載入之前，原片段仍可使用，之後其 blob URL 會被撤銷。

關閉播放器或卸載元件時，請呼叫 `Media.ClearSource(player)`：

```javascript
Media.ClearSource(player);
```

這會取消尚未完成的載入、重設播放器並釋放 blob URL。請在從文件移除元素之前呼叫。框架元件應在卸載或銷毀掛鉤中呼叫它。只移除元素不會釋放 blob URL。請一致地使用這些輔助函式管理播放器來源。

對於 `<video><source ...></video>` 標記，將選取的檔名傳給 `Media.SetSource(video, "media", name)`。輔助函式會設定父播放器的 `src`，其優先順序高於子 `<source>` 元素。清除此屬性後，瀏覽器會再次考慮這些子元素；透過此 API 管理所有來源時，請使用空的播放器。

未連接到文件的音訊物件也能使用：

```javascript
const sound = new Audio();
await Media.SetSource(sound, 'media', 'notification.mp3');
sound.addEventListener('ended', () => Media.ClearSource(sound), { once: true });
// Call sound.play() from an appropriate user interaction.
// Also clear it if playback is cancelled or the owning component is disposed.
```

## 限制下載並取消載入

預設限制為**每個來源 32 MiB**。您可以在 Go 設定的限制內，選擇更小或更大的正整數位元組上限：

```javascript
const controller = new AbortController();
const loading = Media.SetSource(player, 'media', 'welcome.mp4', {
    maxBytes: 8 * 1024 * 1024,
    signal: controller.signal,
});

// Call controller.abort() to cancel this load.
await loading;
```

檔案過大會以 `RangeError` 拒絕。Go 處理常式會在讀取內容前檢查檔案大小，傳輸量最多為其設定限制與前端限制中的較小值。前端也會檢查接收大小，並拒絕不完整的傳輸。取消時會以 `AbortError` 或傳給 `AbortController.abort(reason)` 的原因拒絕。

檔案以 64 KiB 的訊框傳輸。Wails 桌面串流傳輸將輪詢回應限制為 1 MiB，因此 WebView2 緩衝完整回應時，不會在單一回應中緩衝整個媒體檔案。串流佇列本身也有具上限的背壓機制。這些傳輸限制不會改變播放器輔助函式使用完整 blob 的行為。

**整個檔案下載完成後才會播放。** 位元組限制針對個別檔案，並非應用程式的總記憶體限制。多個播放器、取代過程中的原片段、blob 建立與解碼後的媒體，可能使用額外的記憶體。輔助函式一經呼叫就會傳輸，不受播放器的 `preload` 設定影響。請在使用者選擇載入片段時呼叫。

對於大型本機檔案，此輔助函式不是串流播放方案。提高限制也會增加記憶體用量。已透過 HTTP/HTTPS 託管的媒體應使用原生媒體載入，讓瀏覽器能透過範圍請求進行串流播放與定位。

## 排解 Linux 播放問題

- 如果直接播放本機檔案時出現 **No URI handler implemented for "wails"**，請使用 `Media.SetSource` 載入片段。預設 GTK4 堆疊與舊版 `-tags gtk3` 堆疊都受影響。
- 如果載入成功但解碼失敗，請檢查目標系統安裝的 GStreamer 編解碼器。MP4 通常需要 H.264 影片與 AAC 音訊支援；MP3 需要 MP3 解碼器。請在支援的發行版上測試您提供的格式。
- 如果傳輸失敗，請檢查已註冊的串流名稱、相對檔名與檔案系統權限。傳輸期間修改檔案可能造成不完整的傳輸；請在寫入完成後重試。
- 如果應用程式設定了內容安全性政策，請在 `connect-src` 中允許 Wails 資源來源，在 `media-src` 中允許 `blob:`。對於僅限本機的政策，可使用 `connect-src 'self'; media-src 'self' blob:`。請保留其他指令。
- 如果檔案超過限制，請選擇更短或更小的片段，或設定符合應用程式記憶體預算的明確上限。

使用 `go run .` 執行 [audio-video 範例](https://github.com/wailsapp/wails/tree/master/v3/examples/audio-video)，在您的電腦上檢查隨附的 MP3 與 MP4 範例。
