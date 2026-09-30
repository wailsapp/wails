---
title: "バインディングシステム"
description: "Wails v3 で、定型コードを一切記述せずに Go と JavaScript から相互に呼び出す仕組み"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> 「バインディング」は、次のようなコードを記述できるようにする<strong>型安全なコントラクト</strong>です。

```go
msg, err := chatService.Send("Hello")
```

Go では<em>、さらに</em>

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

TypeScript では、<strong>IPC の接続コードを手作業で一切記述せずに</strong>実現できます。 このドキュメントでは、その仕組みが<em>どのように</em>動作するのかを、ビルド時の<strong>静的解析</strong>から<strong>コード生成</strong>、さらに WebView 経由でバイト列を転送する<strong>ランタイムブリッジ</strong>まで順に解説します。

> 正式な詳細解説については[`contributing/architecture/bindings`](/contributing/architecture/bindings/)を参照してください。
>
> ジェネレーターパイプラインに関する決定版の詳細解説を参照してください。このページは、
>
> コントリビューター向けの概要を説明します。

---

## 1. 30秒で分かる概要

| 段階 | コンポーネント | 出力 |
| --- | --- | --- |
| **収集／解析** | `internal/generator/collect/`、`internal/generator/analyse.go` | エクスポートされた Go のサービス、メソッド、パラメーター、戻り値の型、モデルを表すメモリ内モデル |
| **生成** | `internal/generator/render/templates/*.tmpl`（`service.{js,ts}.tmpl`、`models.{js,ts}.tmpl`、`index.tmpl`、`eventcreate.js.tmpl`、`eventdata.d.ts.tmpl`、`newline.tmpl`） | `frontend/bindings/<full Go import path>/...`配下のサービス別 ES モジュール |
| **ランタイム** | `pkg/application/messageprocessor*.go`と、`internal/runtime/desktop/@wailsio/runtime/src/`配下に組み込まれた JS ランタイム（`calls.ts`、`events.ts`、…） | WebView のネイティブブリッジを介した呼び出し／イベントメッセージ |

このフローは`wails3 generate bindings`コマンドによって統括されます。このコマンドは、一連の Go パッケージを対象に、`internal/generator/generate.go`で定義された`generator.Generate`を実行します。

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. 静的解析

### エントリーポイント

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

コレクターのパスは、読み込まれたすべてのパッケージを走査し、次の情報を記録します。

- `collect.ServiceInfo` — エクスポートされ、バインドされた Go 構造体ごとに 1 つ。
- `collect.ServiceMethodInfo`／`collect.MethodInfo` — メソッドごとのシグネチャ情報（名前、パラメーター、結果、エラーの位置、レシーバー、ドキュメント）。
- `collect.ModelInfo`／`collect.StructInfo` — TS／JS モデルとして出力。
- `//wails:inject`、`//wails:include`、`//wails:internal`、`//wails:ignore`、`//wails:id <hex>`などのディレクティブコメント（`internal/generator/collect/directive.go`を参照）。

未対応の型があるとジェネレーターエラーが発生するため、誤りは実行時ではなくビルド時に明らかになります。

### モデル識別子

ランタイムの呼び出しエンベロープでは、完全修飾名（`pkg.Struct.Method`）の<strong>決定論的な FNV-1a ハッシュ</strong>によってメソッドを識別します。生成されたバインディングでは`$Call.ByID(<numeric-id>, …)`として、`-names`を指定して生成した場合は`$Call.ByName("pkg.Struct.Method", …)`として確認できます。

---

## 3. コード生成

### テンプレート

`internal/generator/render/templates/`：

| テンプレート | 目的 |
| --- | --- |
| `service.js.tmpl` | バインドされたサービスごとに 1 つの JS モジュール |
| `service.ts.tmpl` | TypeScript の関連ファイル（`-ts`指定時に生成） |
| `models.js.tmpl` | モデルクラスの出力（パッケージ単位） |
| `models.ts.tmpl` | モデルの`.d.ts`出力（パッケージ単位） |
| `index.tmpl` | パッケージ単位の`index.{js,ts}`バレル再エクスポート |
| `eventcreate.js.tmpl`／`eventdata.d.ts.tmpl` | イベントコンストラクター／ペイロードの型定義 |
| `newline.tmpl` | 末尾改行の正規化 |

出力先は`frontend/bindings/<full Go import path>/...`配下です。たとえば、`github.com/you/yourapp/services/chat`で定義されたサービスは`frontend/bindings/github.com/you/yourapp/services/chat/`配下に配置されます。v3 には`frontend/src/wailsjs/`ディレクトリはありません。

### JavaScript 出力

生成されるバインディングは、`/wails/runtime.js`からランタイムヘルパーをインポートする ES モジュールです。

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

`-names`を指定して生成すると、代わりに`$Call.ByName("pkg.Struct.Method", ...)`が出力されます。これは常に<strong>完全修飾</strong>され、単なる`"Method"`にはなりません。

生成されるモデルクラスでは、フィールドごとの`if (!("X" in $$source))`デフォルト値、引用符で囲まれたフィールド名、および文字列入力に対して`JSON.parse`を実行する`static createFrom(...)`を備えた`$$source`コンストラクターパターンを使用します。

### 型マッピングの要点

`internal/generator/render/`に照らして検証済みです。

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V`（文字列以外の`K`） | `{ [_ in K]?: V }`（`Map<K, V>`でも`Record<K, V>`でもない） |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string`（JSON ISO 8601） |
| `error`（戻り値の位置） | 拒否された Promise |

### リフレクションに関する注記

`pkg/application/bindings.go`は<strong>手書き</strong>であり、`BoundMethod`レジストリからのメソッドディスパッチに`reflect`を使用します。以前の「実行時のリフレクションはゼロ」という記述を文字どおりに受け取らないでください。ジェネレーターはリフレクションを回避しますが、ランタイムディスパッチャーはリフレクションを使用します。

---

## 4. ランタイム呼び出しプロトコル

### JavaScript 側

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

ランタイムヘルパーは、`internal/runtime/desktop/@wailsio/runtime/src/calls.ts`（呼び出しのディスパッチ）、`events.ts`（イベント）などにあります。このツリーには`invoke.ts`も`errors.ts`もありません。正確な通信エンベロープは JS 側の`calls.ts`でエンコードされ、Go 側の`pkg/application/messageprocessor_call.go`でデコードされます。ブリッジをデバッグするときは、この2つのファイルを併せて確認してください。

### Go 側

1. `pkg/application/messageprocessor_call.go`が呼び出しメッセージを受信します。
2. `pkg/application/bindings.go`内で、ID または名前を使用してバインド済みメソッドを検索します（`reflect`によって駆動されます）。
3. バインド済みメソッドを呼び出し、`{result, error}`をシリアライズして JS に返します。

### エラーのマッピング

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise`は結果を伴って解決される |
| `error != nil` | `Promise`は、`message`が Go のエラー文字列である`Error`を伴って拒否される |

---

## 5. Go から JavaScript を呼び出す

バインディングジェネレーターは一方向です（Go のメソッドを JS に公開します）。Go → JS の通信には、イベントバスを使用するか、ウィンドウ内で JS を実行します。

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

JS 側では、`/wails/runtime.js`の`Events.On(name, cb)`を使用して購読します。

---

## 6. 拡張とトラブルシューティング

### サポートされていない型のエラー

```
error: field "Client" uses unsupported type: chan struct{}
```

→ チャネルをメソッド API の背後にラップするか、フィールドに`//wails:internal`を付けて、ジェネレーターがそのフィールドをスキップするようにします。

### 古いバインディング

生成された出力は、`wails3 generate bindings`、`wails3 dev`、または`wails3 build`を実行するたびに上書きされます。IDE の IntelliSense に古いスタブが表示される場合は、`frontend/bindings/`を削除してジェネレーターを再実行してください。`-clean`フラグ（現在のビルドではデフォルトが`true`）を使用すると、実行前に毎回バインディングディレクトリが消去されます。

### パフォーマンスのヒント

- 大きなバイトスライスをブリッジ経由でストリーミングしないでください。代わりに、アセットサーバー経由で配信します。
- レイテンシーが重要な場合は、複数の短い呼び出しを1つのメソッドにまとめてください。
- 小さなパラメータ構造体では、割り当てを減らすために値レシーバーを優先してください。

---

## 7. 主要ファイル一覧

| 対象 | ファイル |
| --- | --- |
| ジェネレーターのオーケストレーション | `internal/generator/generate.go` |
| セマンティックチェック | `internal/generator/analyse.go` |
| 収集（サービス、メソッド、モデル） | `internal/generator/collect/{service,method,model,struct,package}.go` |
| レンダリングテンプレート | `internal/generator/render/templates/*.tmpl` |
| 生成されたバインディングの格納場所 | `frontend/bindings/<full Go import path>/...` |
| Go 側ディスパッチャー | `pkg/application/bindings.go`、`messageprocessor_call.go` |
| JS ランタイム | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

ブリッジのバグを追跡するときに参照できるよう、この早見表を手元に置いてください。

---

## 8. まとめ

1. <strong>コレクター</strong>が Go コードをスキャン → メモリ内セマンティックモデル。
2. <strong>テンプレート</strong>が、サービスごとの ES モジュールとパッケージごとのモデル／インデックスファイルを出力。
3. <strong>メッセージプロセッサー</strong>が、バインディングレジストリを介して Go 側で呼び出しをディスパッチ。
4. <strong>JS ランタイム</strong>が、キャンセル可能で慣用的な Promise にすべてをラップ。

IPC の定型コードを一行も書く必要はありません。これが Wails v3 のバインディングシステムです。さあ、バインドしましょう！
