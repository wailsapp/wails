---
title: "最初のアプリケーション"
description: "初めての Wails デスクトップアプリケーションを手順に沿って作成する"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

このガイドでは、初めての Wails v3 アプリケーションを作成する方法を、プロジェクトのセットアップ、ビルド、開発ワークフローまで含めて説明します。

<br/>

<br/>

@steps
### 新しいプロジェクトの作成
ターミナルを開き、次のコマンドを実行して新しい Wails プロジェクトを作成します。

```bash
wails3 init -n myfirstapp -t vanilla-js
```

このレッスンでは、JavaScript Vanilla テンプレートを明示的に選択します。標準の `vanilla` テンプレートは TypeScript を使用します。このコマンドにより、必要なファイルがすべて含まれた `myfirstapp` が作成されます。

   <video src="/assets/wails_init.mp4" controls></video>

### プロジェクト構成の確認
`myfirstapp` ディレクトリに移動します。次のファイルとフォルダーがあります。

@filetree
- build/           ビルドプロセスで使用するファイルを格納
  - appicon.png  アプリケーションアイコン
  - config.yml   ビルド設定
  - Taskfile.yml Build tasks
  - darwin/      macOS 固有のビルドファイル
    - Info.dev.plist Development configuration
    - Info.plist    本番環境用の設定
    - Taskfile.yml  macOS のビルドタスク
    - icons.icns    macOS のアプリケーションアイコン
  - linux/       Linux 固有のビルドファイル
    - Taskfile.yml  Linux のビルドタスク
    - appimage/     AppImage パッケージング
      - build.sh  AppImage ビルドスクリプト
    - nfpm/        NFPM パッケージング
      - nfpm.yaml Package configuration
      - scripts/  ビルドスクリプト
  - windows/     Windows 固有のビルドファイル
    - Taskfile.yml        Windows のビルドタスク
    - icon.ico           Windows のアプリケーションアイコン
    - info.json          アプリケーションのメタデータ
    - wails.exe.manifest Windows manifest file
    - nsis/              NSIS インストーラーファイル
      - project.nsi                    NSIS プロジェクトファイル
      - wails_tools.nsh               NSIS ヘルパースクリプト
- frontend/        フロントエンドアプリケーションのファイル
  - index.html   メイン HTML ファイル
  - src/         フロントエンドのソースコード
    - main.js    メイン JavaScript ファイル
  - package.json NPM package configuration
  - public/      静的アセット
  - Inter Font License.txt Font license
- .gitignore      Git の除外設定ファイル
- README.md       プロジェクトのドキュメント
- Taskfile.yml    プロジェクトのタスク
- go.mod          Go モジュールファイル
- go.sum          Go モジュールのチェックサム
- greetservice.go Greeting service
- main.go         アプリケーションのメインコード
@end

少し時間を取ってこれらのファイルを確認し、プロジェクト構成を把握しておきましょう。

@note{type="info"}
Wails v3 ではデフォルトのビルドシステムとして [Task](https://taskfile.dev/) を使用しますが、`make` やその他のビルドシステムを使用することもできます。

@end

### アプリケーションのビルド
アプリケーションをビルドするには、次のコマンドを実行します。

```bash
wails3 build
```

生成された Taskfile を使用すると、このコマンドによって `bin` に本番用ビルドが作成されます。開発モードには `wails3 dev` を、明示的に開発用ビルドを作成するには `wails3 build DEV=true` を使用します。

@note{type="info"}
`wails3 build` は `wails3 task build` の短縮形で、`Taskfile.yml` にある `build` タスクを実行します。

@end

     <video src="/assets/wails_build.mp4" controls></video>

ビルドが完了したら、通常のアプリケーションと同じように実行できます。

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

スターター UI が表示されます。本番用ビルドでは、開発者向けインスペクターやコンソールの動作は保証されません。デバッグ時には開発モードを使用してください。

### 開発モード
アプリケーションは開発モードでも実行できます。このモードでは、アプリケーション全体を再ビルドせずにフロントエンドコードを変更し、その変更を実行中のアプリケーションに反映できます。

1. 新しいターミナルウィンドウを開きます。
2. `wails3 dev` を実行します。アプリケーションがコンパイルされ、デバッグモードで起動します。
3. 任意のエディターで `frontend/index.html` を開きます。
4. 名前入力欄の `placeholder="Your name"` を `placeholder="Your full name"` に変更し、ファイルを保存します。
5. ファイルを保存します。

この変更はアプリケーションに即座に反映されます。

バックエンドコードを変更すると、再ビルドが実行されます。

1. `greetservice.go` を開きます。
2. `return "Hello " + name + "!"` がある行を `return "Hello there " + name + "!"` に変更します。
3. ファイルを保存します。

数秒以内にアプリケーションが更新されます。

     <video src="/assets/wails_dev.mp4" controls></video>

### アプリケーションのパッケージ化
アプリケーションを配布する準備が整ったら、プラットフォーム固有のパッケージを作成できます。

@tabs{sync-key="platform"}
[Mac]
`.app` バンドルを作成するには、次のコマンドを実行します。

```bash
wails3 package
```

これにより本番用ビルドが作成され、`bin` ディレクトリ内の `.app` バンドルとしてパッケージ化されます。

[Windows]
NSIS インストーラーを作成するには、次のコマンドを実行します。

```bash
wails3 package
```

これにより本番用ビルドが作成され、`bin` ディレクトリ内の NSIS インストーラーとしてパッケージ化されます。

[Linux]
Wails は、Linux 配布用に複数のパッケージ形式をサポートしています。

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

パッケージ化のオプションと設定について詳しくは、[ビルドおよびパッケージ化ガイド](/guides/build/building/)を参照してください。

### バージョン管理とモジュール名の設定
プロジェクトは、仮のモジュール名 `changeme` で作成されます。リポジトリの URL と一致するように更新することをお勧めします。

1. GitHub（または任意の Git ホスティングサービス）に新しいリポジトリを作成します
2. プロジェクトディレクトリで Git を初期化します。
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. リモートリポジトリを設定します（実際のリポジトリ URL に置き換えてください）。
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. `go.mod` のモジュール名を、リポジトリの URL と一致するように更新します。
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. コードをプッシュします。
  ```bash
  git push -u origin main
  ```


これにより、Go モジュール名が Go のモジュール命名規則に準拠し、コードを共有しやすくなります。

@note{type="tip" title="ヒント"}
プロジェクトの作成時に `-git` フラグを使用すると、すべての初期化手順を自動化できます。

```bash
wails3 init -n myfirstapp -t vanilla-js -git github.com/username/myfirstapp
```

さまざまな Git URL 形式に対応しています。

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` または `ssh://git@github.com/username/project`
- Git プロトコル: `git://github.com/username/project`
- ファイルシステム: `file:///path/to/project.git`

@end

@end

## おめでとうございます！

これで、初めての Wails アプリケーションの作成、開発、パッケージ化が完了しました。Wails で実現できることは、まだ始まったばかりです v3。

## 次のステップ

Wails を初めて使用する場合は、次にチュートリアルを一通り読むことをお勧めします。Wails のさまざまな機能を実践的に学べます。最初のチュートリアルは、[サービスの作成](/tutorials/01-creating-a-service/)です。

より高度な知識をお持ちの場合は、Wails の使用方法に関する詳細情報について、[ビルドおよびパッケージ化ガイド](/guides/build/building/)を参照してください。
