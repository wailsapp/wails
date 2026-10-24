---
title: "ローカルの音声と動画を再生する"
description: "サイズ制限付きの Blob URL で Linux 上の同梱メディアを再生し、終了時に解放します。"
slug: "guides/linux-media"
sourcePath: "guides/linux-media.md"
---

`Media.SetSource` を使用すると、Wails v3 アプリケーションで短いローカルクリップを再生できます。Linux では WebKitGTK がメディア再生を GStreamer に委ねますが、GStreamer は `wails://` URL を直接読み込めません。このヘルパーは Wails ストリーム経由でクリップを受信し、プレーヤーに Blob URL を設定します。

デスクトップでの転送は、待ち受けソケットを開かずに既存の Wails アセット転送機構を使用します。API は音声要素と動画要素で動作します。通常の HTTP/HTTPS メディアは、プレーヤーの `src` に直接設定できます。

## メディアファイルを登録する

フロントエンドで再生を許可するクリップを含むファイルシステムを公開し、名前付きストリームにメディアハンドラーを登録します。

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

アプリケーションの作成後、`app.Run()` を呼ぶ前に `registerMedia(app)` を呼び出してください。

ディスク上のファイルには `os.OpenRoot(directory)` を使用し、`root.FS()` を `media.NewHandler` に渡します。`app.Run()` が戻るまでルートを開いたままにし、その後閉じてください。フロントエンドが読み取ってよいファイルだけを含むディレクトリを選んでください。シンボリックリンクがディレクトリ外を指していても、`os.Root` はアクセスを制限します。

## クリップを読み込む

プレーヤーを作成します。

```html
<video id="player" controls></video>
```

npm ランタイムを使用するフロントエンドの場合：

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

同梱ランタイムを使用するアプリケーションでは、インポートを次のように変更します。

```javascript
import { Media } from '/wails/runtime.js';
```

第2引数は登録済みのストリーム名です。第3引数は、そのファイルシステムのルートからの相対パスで、`welcome.mp4` や `tutorials/intro.mp4` のようにスラッシュで区切ります。URL や OS のパスではありません。

ソースが設定されると Promise が解決され、その後プレーヤーがデコードします。未対応のコーデックを検出するにはプレーヤーの `error` イベントを処理してください。自分で再生を開始する場合は、ユーザー操作から `player.play()` を呼び出してください。

## クリップを置き換える、または解放する

クリップを変更するには `Media.SetSource` を再度呼び出します。そのプレーヤーの以前の未完了の読み込みがキャンセルされるため、遅い応答が最新の選択を上書きすることはありません。置き換えるクリップの読み込みが成功するまで、前のクリップを利用できます。その後、前の Blob URL が無効化されます。

プレーヤーを閉じるときやコンポーネントをアンマウントするときは、`Media.ClearSource(player)` を呼び出してください。

```javascript
Media.ClearSource(player);
```

これにより未完了の読み込みがキャンセルされ、プレーヤーがリセットされ、Blob URL が解放されます。ドキュメントから要素を削除する前に呼び出してください。フレームワークのコンポーネントでは、アンマウントまたは破棄のフックから呼び出してください。要素の削除だけでは Blob URL は解放されません。プレーヤーのソースには、これらのヘルパーを一貫して使用してください。

`<video><source ...></video>` というマークアップでは、選択したファイル名を `Media.SetSource(video, "media", name)` に渡します。ヘルパーは親プレーヤーの `src` を設定し、子の `<source>` より優先されます。この属性を消去すると、ブラウザーは再び子要素を考慮できます。この API ですべてのソースを管理する場合は、空のプレーヤーを使用してください。

ドキュメントに接続されていない音声オブジェクトでも動作します。

```javascript
const sound = new Audio();
await Media.SetSource(sound, 'media', 'notification.mp3');
sound.addEventListener('ended', () => Media.ClearSource(sound), { once: true });
// Call sound.play() from an appropriate user interaction.
// Also clear it if playback is cancelled or the owning component is disposed.
```

## ダウンロードを制限し、読み込みをキャンセルする

デフォルトの制限は **ソースごとに 32 MiB** です。Go 側で設定した制限の範囲内で、バイト単位の正の整数として、より小さい値や大きい値を指定できます。

```javascript
const controller = new AbortController();
const loading = Media.SetSource(player, 'media', 'welcome.mp4', {
    maxBytes: 8 * 1024 * 1024,
    signal: controller.signal,
});

// Call controller.abort() to cancel this load.
await loading;
```

サイズが大きすぎるファイルは `RangeError` で拒否されます。Go ハンドラーは内容を読む前にファイルサイズを確認し、自身に設定された制限とフロントエンドの制限のうち小さい方までしか転送しません。フロントエンドも受信サイズを確認し、不完全な転送を拒否します。キャンセル時には `AbortError`、または `AbortController.abort(reason)` に渡した理由で拒否されます。

ファイルは 64 KiB のフレームで転送されます。Wails のデスクトップストリーム転送はポーリング応答を 1 MiB に制限するため、完全な応答をバッファリングする WebView2 でも、メディアファイル全体が単一の応答に蓄積されることはありません。ストリームキューにも独自の上限付きバックプレッシャーがあります。これらの転送制限は、プレーヤーヘルパーが完全な Blob を扱う動作を変えません。

**再生前にファイル全体をダウンロードします。** バイト制限はファイルごとのもので、アプリケーション全体のメモリ上限ではありません。複数のプレーヤー、置き換え中の前のクリップ、Blob の構築、デコード済みメディアは追加のメモリを使用する場合があります。ヘルパーはプレーヤーの `preload` 設定に関係なく、呼び出された時点で転送します。ユーザーがクリップの読み込みを選択したときに呼び出してください。

大きなローカルファイルに対して、このヘルパーはストリーミングの解決策にはなりません。制限を上げるとメモリ使用量も増えます。HTTP/HTTPS ですでにホストされているメディアにはネイティブのメディア読み込みを使用し、ブラウザーが範囲リクエストでストリーミングやシークを行えるようにしてください。

## Linux の再生問題を調べる

- ローカルファイルの直接再生で **No URI handler implemented for "wails"** と表示される場合は、`Media.SetSource` でクリップを読み込んでください。デフォルトの GTK4 と従来の `-tags gtk3` の両方に該当します。
- 読み込みに成功してもデコードに失敗する場合は、対象システムにインストールされた GStreamer コーデックを確認してください。MP4 には通常 H.264 動画と AAC 音声のサポートが必要で、MP3 には MP3 デコーダーが必要です。配布する形式をサポート対象のディストリビューションでテストしてください。
- 転送に失敗する場合は、登録済みのストリーム名、相対ファイル名、ファイルシステムの権限を確認してください。転送中のファイル変更は不完全な転送の原因になります。書き込み完了後に再試行してください。
- アプリケーションに Content Security Policy を設定する場合は、`connect-src` で Wails アセットのオリジンを、`media-src` で `blob:` を許可してください。ローカル専用のポリシーでは `connect-src 'self'; media-src 'self' blob:` と指定できます。他のディレクティブは保持してください。
- ファイルが制限を超える場合は、より短い、または小さいクリップを選ぶか、アプリケーションのメモリ予算に合う明示的な制限を設定してください。

[audio-video サンプル](https://github.com/wailsapp/wails/tree/master/v3/examples/audio-video)を `go run .` で実行し、同梱の MP3 と MP4 サンプルをお使いの環境で確認してください。
