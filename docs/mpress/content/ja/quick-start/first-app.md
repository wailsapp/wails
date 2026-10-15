---
title: "初めてのアプリ"
description: "10 分で動作する Wails アプリケーションを構築する"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Wails の中核となる次の概念を示す、シンプルな挨拶アプリケーションを構築します。

- ロジックを管理する Go バックエンド
- Go 関数を呼び出すフロントエンド
- 型安全なバインディング
- 開発中のホットリロード

<strong>所要時間：</strong>10 分

@note{type="tip" title="Windows 11 ユーザー向けのパフォーマンスに関するヒント"}
プロジェクトの保存先として [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) の使用を検討してください。Dev Drive は開発者のワークロード向けに最適化されており、通常の NTFS ドライブと比べて、ビルド時間とディスクアクセス速度を最大 30% 改善できます。

@end

## プロジェクトを作成する

@steps
### プロジェクトを生成する
```bash
wails3 init -n myapp
cd myapp
```

これにより、デフォルトの Vanilla + Vite テンプレート（Vite バンドラーを使用する HTML/CSS/TypeScript）で新しいプロジェクトが作成されます。

@note{type="tip" title="その他のテンプレート"}
好みのフレームワークに応じて、`-t react`、`-t vue`、または `-t svelte` を試してください。これらは デフォルトでは TypeScript を使用します。プレーン JavaScript を使用する場合は、`-t vanilla-js` または `-t react-js` を使用してください。 利用可能なすべてのテンプレートを確認するには `wails3 init -l` を実行します。または、 [独自のフロントエンドフレームワークを使用](/guides/dev/frontend-frameworks/)できます。

@end

### プロジェクト構造を理解する
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.ts      # Frontend TypeScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   └── vite.config.ts   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### アプリを実行する
```bash
wails3 dev
```

@note{type="info" title="初回実行"}
初回実行時は、フロントエンドの依存関係のインストールやバインディングの生成などが行われるため、予想より時間がかかることがあります。2 回目以降は大幅に速くなります。

@end

アプリが開き、挨拶用のインターフェースが表示されます。名前を入力して「Greet」をクリックすると、Go バックエンドが入力を処理して挨拶を返します。

@end

## 仕組み

この仕組みを実現するコードを見ていきましょう。

### Go バックエンド

`greetservice.go` を開きます。

```go {title="greetservice.go"}
package main

import (
	"fmt"
)

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
```

**主要な概念：**

1. **サービス** — エクスポートされたメソッドを持つ Go 構造体
2. **エクスポートされたメソッド** — `Greet` は先頭が大文字であるため、フロントエンドから利用可能
3. **シンプルなロジック** — 名前を受け取り、挨拶を返す
4. **型安全性** — 入力型と出力型が定義されている

@note{type="tip" title="サービスとバインディングについて"}
<strong>サービス</strong>は、フロントエンドに機能を公開する自己完結型の Go モジュールです。実体は通常の Go 構造体であり、エクスポートされたメソッドを持ち、アプリケーション設定の `Services` フィールドに登録します。

<strong>バインディング</strong>は、フロントエンドからこれらのサービスを呼び出せるようにする、自動生成された TypeScript/JavaScript SDK です。`wails3 dev` または `wails3 build` を実行すると、Wails は登録済みのサービスを解析し、型安全なバインディングを `frontend/bindings/` に生成します。

サービスはバックエンド API、バインディングはそれと通信するクライアントライブラリと考えてください。

@end

### サービスを登録する

`main.go` を開き、サービスの登録箇所を探します。

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

これにより `GreetService` が Wails に登録され、エクスポートされたすべてのメソッドをフロントエンドから利用できるようになります。

### フロントエンド

`frontend/src/main.js` を開きます。

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.greet = async () => {
    const nameElement = document.getElementById('name');
    const resultElement = document.getElementById('result');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        resultElement.innerText = result;
    } catch (err) {
        console.error(err);
    }
};
```

**主要な概念：**

1. **自動生成されたバインディング** — `GreetService` は生成されたコードからインポートされる
2. **型安全な呼び出し** — メソッド名とシグネチャが Go コードと一致する
3. **デフォルトで非同期** — すべての Go 呼び出しが Promise を返す
4. **エラー処理** — Go からのエラーは try/catch で捕捉される

@note{type="info" title="バインディングの場所"}
生成されたバインディングは `frontend/bindings/` にあります。`wails3 dev` または `wails3 build` を実行すると、自動的に作成されます。

**これらのファイルは絶対に手動で編集しないでください**。ビルドするたびに再生成されます。

@end

## アプリをカスタマイズする

ワークフローを理解するため、新しい機能を追加してみましょう。

### 「Greet Many」機能を追加する

@steps
### GreetService にメソッドを追加する
次の内容を `greetservice.go` に追加します。

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### アプリは自動的に再ビルドされる
ファイルを保存すると、`wails3 dev` が Go コードを自動的に再ビルドし、アプリを再起動します。

@note{type="info" title="自動再ビルド"}
Go コードを変更すると、自動的に再ビルドと再起動が行われます。フロントエンドの変更は、再起動せずにホットリロードされます。

@end

### フロントエンドで使用する
次の内容を `frontend/src/main.js` に追加します。

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

ブラウザーコンソールを開いて `greetMany()` を呼び出すと、挨拶の配列が表示されます。

@end

## 本番環境向けにビルドする

アプリを配布する準備ができたら、次を実行します。

```bash
wails3 build
```

**実行される処理：**

- 最適化を有効にして Go コードをコンパイルする
- フロントエンドを本番環境向けにビルドする（ミニファイを実施）
- `bin/` にネイティブ実行ファイルを作成する

@tabs{sync-key="os"}
[Windows]
**出力：** `bin/myapp.exe`

ダブルクリックして実行します。依存関係は必要ありません（WebView2 は Windows に含まれています）。

[macOS]
**出力：** `bin/myapp.app`

Applications フォルダーにドラッグするか、ダブルクリックして実行します。

[Linux]
**出力：** `bin/myapp`

`./bin/myapp` で実行するか、ランチャー用の `.desktop` ファイルを作成します。

@end

@note{type="tip" title="クロスプラットフォームビルド"}
ほかのプラットフォーム向けにビルドする場合は、[クロスプラットフォームビルド →](/guides/build/cross-platform/) を参照してください。

@end

## ここまでに学んだこと

**プロジェクト構成**

- Go バックエンド用の `main.go`
- UI コード用の `frontend/`
- ビルドタスク用の `Taskfile.yml`

**サービス**

- エクスポートされたメソッドを持つ Go の構造体を作成
- `application.NewService()` で登録
- メソッドをフロントエンドから自動的に利用可能

**バインディング**

- TypeScript 定義を自動生成
- 型安全な関数呼び出し
- デフォルトで非同期（Promise）

**開発ワークフロー**

- ホットリロード用の `wails3 dev`
- Go の変更時に自動で再ビルドして再起動
- フロントエンドの変更を即座にホットリロード

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) に参加して、コミュニティに質問してください。
