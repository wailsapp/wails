---
title: "ノッチウィンドウ"
description: "カメラハウジングに接続されたネイティブ macOS ウィンドウを作成します。"
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` は、カメラハウジングに接続された、形状付きでアプリケーションをアクティブ化しない macOS パネルを作成します。Wails は、ネイティブでの配置、黒い接続ウィング、透明な外側キャンバス、ウィンドウレベル、Spaces での動作、およびオプションの表示・非表示アニメーションを管理します。Web コンテンツが占めるのは、指定された内側の矩形領域だけです。ポインターをウィンドウ内に移動すると、アプリケーションをアクティブ化せずに、その WebView が直ちにキーウィンドウになります。

@note{type="caution" title="WebView の透明化にはプライベート API が使用されます"}
透明な Web コンテンツからネイティブのノッチ形状が見えるようにするには、`NewNotchWindow` に `-tags private_mac_apis` が必要です。このタグがなくても、パネル、配置、アニメーションは機能しますが、WebView は不透明なままです。[macOS のプライベート API](/guides/build/private-macos-apis/#webview-transparency-and-background)を参照してください。

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="MacBook のカメラハウジングから下にスライドしてライブのシステムメトリクスを表示し、ノッチの下へ戻って隠れる Wails のノッチ通知"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    永続的な WebView とネイティブの表示・非表示トランジションを使用するアニメーション付きノッチ通知。
  </figcaption>
</figure>

```go
alert := app.Window.NewNotchWindow(application.NotchWindowOptions{
    Width:    660,
    Height:   92,
    Animated: true,
    WindowOptions: application.WebviewWindowOptions{
        Name: "alert",
        URL:  "/alert",
    },
})

alert.Show()
alert.Hide()
visible := alert.Visibility()
alert.Close()
```

## オプション

| フィールド | 型 | デフォルト | 説明 |
| --- | --- | --- | --- |
| `Width` | `int` | `660` | ネイティブ形状の端部の内側で使用できる WebView の幅。 |
| `Height` | `int` | `92` | ネイティブ形状の端部の内側で使用できる WebView の高さ。 |
| `Animated` | `bool` | `false` | `Show` の際にウィンドウを下へ、`Hide` の際に上へスライドさせます。 |
| `AnimationSpeed` | `time.Duration` | `420ms` | 表示にかける時間。非表示にはこの値の 3 分の 2 が使用され、デフォルトは `280ms` です。 |
| `Screen` | `*Screen` | プライマリディスプレイ | 特定のディスプレイを対象にします。指定しない場合、Wails はプライマリディスプレイを使用します。 |
| `WindowOptions` | `WebviewWindowOptions` | デフォルト値 | 名前、URL または HTML、CSS、JavaScript、キーバインド、およびその他の WebView の動作を指定します。 |

`NewNotchWindow` は、外側のサイズ、位置、フレーム、透明度、サイズ変更ポリシー、ネイティブパネルクラス、ウィンドウレベル、およびコレクション動作を管理します。`WindowOptions` に指定されたこれらのフィールドの値は、意図的に置き換えられます。その他のフィールドは保持されます。ウィンドウをカメラハウジングに接続したままにするため、ネイティブの背景ドラッグと CSS のドラッグ領域は無効になります。

返される `NotchWindow` が公開するのは、意図的に `Show`、`Hide`、`Visibility`、`Close` のみです。高レベルのハンドルからネイティブのジオメトリを変更することはできません。

@note{type="note"}
ノッチウィンドウには macOS が必要です。カメラハウジングのない Mac では、Wails はメニューバーの下にある画面上部中央へウィンドウを配置します。サポートされていないプラットフォームでは、`NewNotchWindow` は何も行わないハンドルを返します。そのライフサイクルメソッドは安全な no-op であり、`Visibility` は常に false です。

@end

## ライフサイクル

- `Show` は、既存のネイティブウィンドウを表示します。アニメーションが有効な場合、ウィンドウはディスプレイの上側から下へスライドします。
- `Hide` は、再利用できるようにウィンドウを存続させます。アニメーションが有効な場合、ウィンドウはディスプレイの上側へ戻るようにスライドしてから、表示順序から外されます。WebView、JavaScript の状態、バインディング、イベントリスナーは読み込まれたままですが、非表示のウィンドウはホバー対象になりません。再表示するには、アプリケーションが `Show` を呼び出す必要があります。
- `Visibility` は、現在のネイティブの表示状態を報告します。
- `Close` は、ネイティブウィンドウを完全に破棄します。その通知を再び表示する前に、新しいウィンドウを作成してください。
- ポインターが入ると、そのノッチウィンドウが前面に移動し、アプリケーションをアクティブ化しないパネルの動作を維持したまま、直ちにキーボード操作できるように WebView にフォーカスが移ります。

`NewNotchWindow` を呼び出すたびに、独自のコンテンツ、表示状態、アニメーション状態を持つ独立したウィンドウが作成されます。各ウィンドウは同じネイティブレベルを使用するため、最後に表示された、またはポインターが入ったインスタンスが前面に現れ、それ以前のインスタンスに重なることがあります。各ウィンドウが異なる画面を対象とする場合、それぞれの画面のカメラハウジングを基準に中央へ配置されます。macOS は、異なるアプリケーションに属するノッチウィンドウ間の調整を行いません。別々のアプリケーションが同じ位置とレベルにウィンドウを表示した場合、最後に表示順序を設定されたウィンドウが前面に現れます。

通知用途では、一般にアプリケーションは非表示のウィンドウを再利用するか、1 つの WebView 内で独自のキューを管理します。Wails は、キューイング、置換、自動的な非表示、または単一ウィンドウのポリシーを強制しません。

@note{type="note"}
ホバー時に再び開く永続的な折りたたみサーフェスは、通知の非表示とは意図的に分離されており、[#6009](https://github.com/wailsapp/wails/issues/6009) で追跡されています。

@end

通知の表示と非表示を繰り返してもライブの JavaScript 状態が維持されるコンパクトなシステムモニターについては、[notch-notification のサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples/notch-notification)を参照してください。
