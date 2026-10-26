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

**所要時間:** 10 分

## プロジェクトを作成する

@steps
### プロジェクトを生成する
```bash
wails3 init -n myapp -t vanilla-js
cd myapp
```

このレッスンでは、以下のソースファイルが生成されるプロジェクトと一致するように、JavaScript の Vanilla + Vite を明示的に選択します。素の `vanilla` テンプレートでは TypeScript が使用されます。

@note{type="tip" title="その他のテンプレート"}
好みのフレームワークに応じて、`-t react`、`-t vue`、または `-t svelte` を試してください。これらはデフォルトで TypeScript を使用します。プレーン JavaScript を使用する場合は、`-t vanilla-js` または `-t react-js` を使用してください。利用可能なすべてのテンプレートを確認するには `wails3 init -l` を実行します。または、[独自のフロントエンドフレームワークを使用する](/guides/dev/frontend-frameworks/)こともできます。

@end

### プロジェクト構造を理解する
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.js      # Frontend JavaScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # Editor and JavaScript tooling configuration
│   └── vite.config.js   # Vite bundler config
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

アプリが起動し、挨拶用のインターフェイスが表示されます。名前を入力して「Greet」をクリックすると、Go バックエンドが入力を処理して挨拶を返します。

@end

@note{type="tip" title="Windows 11 ユーザー向けのパフォーマンスに関するヒント"}
プロジェクトの保存先として [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) の使用を検討してください。Dev Drive は開発者のワークロード向けに最適化されており、通常の NTFS ドライブと比べて、ビルド時間とディスクアクセス速度を最大 30% 改善できます。

@end

## 仕組み

この仕組みを実現しているコードを確認しましょう。

### Go バックエンド

`greetservice.go` を開きます。

```go {title="greetservice.go"}
package main

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}
```

**主要な概念:**

1. **サービス** - エクスポートされたメソッドを持つ Go 構造体
2. **エクスポートされたメソッド** - `Greet` は先頭が大文字になっているため、フロントエンドから利用できます
3. **シンプルなロジック** - 名前を受け取り、挨拶を返します
4. **型安全性** - 入力と出力の型が定義されています

@note{type="tip" title="サービスとバインディングについて"}
**サービス**は、機能をフロントエンドに公開する自己完結型の Go モジュールです。実体は、アプリケーション設定の `Services` フィールドに登録する、エクスポートされたメソッドを持つ通常の Go 構造体です。

**バインディング**は、フロントエンドからこれらのサービスを呼び出せるようにする、自動生成された TypeScript/JavaScript SDK です。`wails3 dev` または `wails3 build` を実行すると、Wails が登録済みのサービスを解析し、`frontend/bindings/` に型安全なバインディングを生成します。

サービスはバックエンド API、バインディングはそれと通信するクライアントライブラリと考えるとよいでしょう。

@end

### サービスを登録する

`main.go` を開き、サービスの登録箇所を見つけます。

```go {title="main.go" highlight="4-6"}
app := application.New(application.Options{
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

window.doGreet = async () => {
    const nameElement = document.getElementById('name');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        showToast(result);
    } catch (err) {
        console.error(err);
    }
};
```

この簡略化したハンドラーで、スターターの `window.doGreet` を置き換えます。既存の `showToast` ヘルパーとその他の起動コードは残してください。生成された HTML ボタンは `doGreet` を呼び出し、`showToast` がその結果を表示します。インポートパスは Go モジュール名に従います（上書きしていない場合は `changeme`）。

**主要な概念:**

1. **自動生成されたバインディング** - `GreetService` は生成されたコードからインポートされます
2. **型安全な呼び出し** - メソッド名とシグネチャが Go コードと一致します
3. **デフォルトで非同期** - すべての Go 呼び出しは Promise を返します
4. **エラー処理** - Go からのエラーは try/catch で捕捉されます

@note{type="info" title="バインディングはどこにありますか？"}
生成されたバインディングは `frontend/bindings/` にあります。`wails3 dev` または `wails3 build` を実行すると、自動的に作成されます。

**これらのファイルは決して手動で編集しないでください**。ビルドのたびに再生成されます。

@end

## アプリをカスタマイズする

ワークフローを理解するために、新しい機能を追加しましょう。

### 「Greet Many」機能を追加する

@steps
### GreetService にメソッドを追加する
次の内容を `greetservice.go` に追加します。

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = "Hello " + name + "!"
    }
    return greetings
}
```

### アプリが自動的に再ビルドされる
ファイルを保存すると、`wails3 dev` が Go コードを自動的に再ビルドし、アプリを再起動します。

@note{type="info" title="自動再ビルド"}
Go コードを変更すると、自動的に再ビルドおよび再起動されます。フロントエンドの変更は、再起動せずにホットリロードされます。

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

**実行される処理:**

- 最適化を有効にして Go コードをコンパイルする
- フロントエンドを本番環境向けにビルドする（圧縮）
- `bin/` にネイティブ実行ファイルを作成する

@tabs{sync-key="os"}
[Windows]
**出力:** `bin/myapp.exe`

ダブルクリックして実行します。依存関係は不要です（WebView2 は Windows に含まれています）。

[macOS]
**出力:** `bin/myapp.app`

Applications フォルダーにドラッグするか、ダブルクリックして実行します。

[Linux]
**出力:** `bin/myapp`

`./bin/myapp` で実行するか、ランチャー用の `.desktop` ファイルを作成します。

@end

@note{type="tip" title="クロスプラットフォームビルド"}
ほかのプラットフォーム向けにビルドする場合は、[クロスプラットフォームビルド →](/guides/build/cross-platform/) を参照してください。

@end

## 学んだこと

**プロジェクト構造**

- Go バックエンド用の `main.go`
- UI コード用の `frontend/`
- ビルドタスク用の `Taskfile.yml`

**サービス**

- エクスポートされたメソッドを持つ Go 構造体を作成
- `application.NewService()` で登録
- メソッドはフロントエンドで自動的に利用可能

**バインディング**

- 自動生成される TypeScript 型定義
- 型安全な関数呼び出し
- デフォルトで非同期（Promise）

**開発ワークフロー**

- ホットリロード用の `wails3 dev`
- Go の変更時に自動的に再ビルドして再起動
- フロントエンドの変更を即座にホットリロード

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) に参加して、コミュニティに質問してください。
