---
title: "ノッチウィンドウ"
description: "カメラハウジングに接続されたネイティブ macOS ウィンドウを作成します。"
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` は、カメラハウジングに接続された、アプリケーションをアクティブ化しない成形済みの macOS パネルを作成します。ネイティブでの配置、左右の黒い接続部、透明な外側キャンバス、ウィンドウレベル、Spaces での動作、およびオプションの表示・非表示アニメーションは Wails が管理します。Web コンテンツは、指定された内側の矩形領域のみを占有します。ポインターをウィンドウ内に移動すると、アプリケーションをアクティブ化せずに、その WebView が直ちにキーウィンドウになります。

@note{type="caution" title="WebView の透明化にはプライベート API を使用します"}
透明な Web コンテンツからネイティブのノッチ形状を見えるようにするには、`NewNotchWindow` に `-tags private_mac_apis` が必要です。このタグがなくても、パネル、配置、アニメーションは機能しますが、WebView は不透明なままです。[macOS のプライベート API](/guides/build/private-macos-apis/#webview-transparency-and-background)を参照してください。

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="MacBook のカメラハウジングから下へスライドしてライブのシステム指標を表示し、再びノッチの下へ隠れる Wails ノッチ通知"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    永続的な WebView とネイティブの表示・非表示トランジションを使用する
    アニメーション付きノッチ通知。
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
| `Width` | `int` | `660` | ネイティブの成形済みエッジ内で使用可能な WebView の幅。 |
| `Height` | `int` | `92` | ネイティブの成形済みエッジ内で使用可能な WebView の高さ。 |
| `Animated` | `bool` | `false` | `Show` ではウィンドウを下へ、`Hide` では上へスライドさせます。 |
| `AnimationSpeed` | `time.Duration` | `420ms` | 表示にかける時間。非表示にはこの値の 3 分の 2 が使用され、デフォルトは `280ms` です。 |
| `Screen` | `*Screen` | プライマリディスプレイ | 特定のディスプレイを対象にします。指定しない場合、Wails はプライマリディスプレイを使用します。 |
| `WindowOptions` | `WebviewWindowOptions` | デフォルト値 | 名前、URL または HTML、CSS、JavaScript、キーバインド、およびその他の WebView の動作を指定します。 |

外側のサイズ、位置、フレーム、透明度、サイズ変更ポリシー、ネイティブパネルクラス、ウィンドウレベル、およびコレクション動作は、`NewNotchWindow` が管理します。`WindowOptions` にあるこれらのフィールドの値は、意図的に置き換えられます。その他のフィールドは保持されます。ウィンドウをカメラハウジングに接続したままにするため、ネイティブの背景ドラッグと CSS ドラッグ領域は無効化されます。

返される `NotchWindow` が公開するのは、意図的に `Show`、`Hide`、`Visibility`、`Close` のみに限定されています。高レベルのハンドルからネイティブジオメトリを変更することはできません。

@note{type="note"}
ノッチウィンドウには macOS が必要です。カメラハウジングのない Mac では、Wails はウィンドウをメニューバーの下、画面上部中央に配置します。サポート対象外のプラットフォームでは、`NewNotchWindow` は何も行わないハンドルを返します。そのライフサイクルメソッドは安全な no-op であり、`Visibility` は常に false です。

@end

## ライフサイクル

- `Show` は既存のネイティブウィンドウを表示します。アニメーションが有効な場合、ディスプレイの上方から下へスライドします。
- `Hide` は、再利用できるようにウィンドウを存続させます。アニメーションが有効な場合、ウィンドウはディスプレイの上方へスライドして戻った後、画面から除外されます。WebView、JavaScript の状態、バインディング、およびイベントリスナーは読み込まれたままですが、非表示のウィンドウにはホバー対象がありません。再表示するには、アプリケーションが `Show` を呼び出す必要があります。
- `Visibility` は、現在のネイティブでの表示状態を返します。
- `Close` はネイティブウィンドウを完全に破棄します。その通知を再び表示する前に、新しいウィンドウを作成してください。
- ポインターが入ると、そのノッチウィンドウが前面に移動し、WebView にフォーカスが設定されます。これにより、アプリケーションをアクティブ化しないパネルの動作を維持しながら、直ちにキーボード操作を行えます。

`NewNotchWindow` を呼び出すたびに、固有のコンテンツ、表示状態、アニメーション状態を持つ独立したウィンドウが作成されます。各ウィンドウは同じネイティブレベルを使用するため、最後に表示された、またはポインターが入ったインスタンスが前面に現れ、それ以前のインスタンスに重なることがあります。各ウィンドウが異なる画面を対象とする場合、それぞれの画面のカメラハウジングを基準に中央揃えされます。macOS は、異なるアプリケーションに属するノッチウィンドウ間の調整を行いません。別々のアプリが同じ位置とレベルにウィンドウを表示した場合、最後に表示順が設定されたウィンドウが前面に現れます。

通知処理では、非表示のウィンドウを再利用するか、1 つの WebView 内でアプリケーション独自のキューを管理するのが一般的です。Wails は、キューイング、置換、自動消去、または単一ウィンドウのポリシーを強制しません。

@note{type="note"}
ホバー時に再度開く、折りたたまれた永続的なサーフェスは、意図的に通知の非表示処理とは分離されており、[#6009](https://github.com/wailsapp/wails/issues/6009) で追跡されています。

@end

通知の表示と非表示を繰り返しても JavaScript のライブ状態が維持されるコンパクトなシステムモニターについては、[notch-notification のサンプル](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notch-notification)を参照してください。
