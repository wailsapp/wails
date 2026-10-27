---
title: "コントリビューション"
description: "Wails に貢献する"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## コントリビューターの皆さまへ

Wails へのコントリビューションを歓迎します。バグの修正、機能の追加、ドキュメントの改善など、皆さまのご協力に感謝します。

## 貢献する方法

### 1. 問題を報告する

バグを見つけた場合は、次の情報を添えて[Issue を作成してください](https://github.com/wailsapp/wails/issues/new)。

- 明確な説明
- 再現手順
- 期待される動作と実際の動作
- システム情報
- コード例

### 2. ドキュメントを改善する

修正 PR には、事前の Issue や失敗するコードテストは必要ありません。  
ダウンロード、macOS／Linux／Windows へのインストール、ライブプレビュー、変更の提出については、[M-Press でドキュメントを作成する](/contributing/documentation/)に従ってください。

ドキュメントの改善はいつでも歓迎します。

- 誤字や誤りを修正する
- 例を追加する
- 説明を明確にする
- コンテンツを翻訳する

### 3. コードを提出する

プルリクエストを通じてコードを提供できます。

- バグ修正
- 新機能
- パフォーマンスの改善
- テスト

### 4. 機能強化を提案する（WEP）

新機能や公開動作の変更には、Wails Enhancement Proposal（WEP）プロセスを使用します。このプロセスにより、機能開発の透明性を保ち、承認されたすべての提案に実装担当者がいることを保証します。機能リクエストの Issue は作成しないでください。

1. 必要に応じて、関心の度合いを確認するため、GitHub Discussions の[Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas)カテゴリーまたは[Discord](https://discord.gg/JDdSxwjhGf)でアイデアを共有してください。
2. [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/WEP_TEMPLATE.md)を`v3/wep/proposals/<proposal name>/proposal.md`にコピーし、すべてのセクションを記入してください。
3. 提案だけを含むドラフトのプルリクエストを、タイトルを`[WEP] <title>`として作成してください。この PR が正式な議論の場になります。
4. フィードバックと支持（PR へのコメントと賛成のリアクション）を集めてください。議論には少なくとも 2 週間を設け、提案の実装担当者について合意してください。
5. PR をレビュー可能な状態にしてください。最終決定はメンテナーが行います。承認された提案には WEP 番号が割り当てられ、マージされます。

プロセス全体については、[`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)を参照してください。

## はじめに

### フォークとクローン

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### ソースからビルドする

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### テストを実行する

テストは変更の一部であり、最後に行う簡単な動作確認ではありません。単体テストは対象コードの近くに配置し、複数の入力やエッジケースで同じ動作を確認する場合は、できるだけテーブル駆動テストを使用してください。失敗時にどのシナリオか分かるよう、各ケースに名前を付けてください。

新規または変更したロジックは、完全に網羅する必要があります。PR で追加または変更するコードについて、Go のステートメントカバレッジ 100% を目標にしてください。リポジトリ全体のカバレッジ率を、その変更のテストの代わりにしてはいけません。たとえば、特定の OS でのみ発生するエラーパスや、実機なしでは再現が現実的でない条件など、カバレッジに不足があっても妥当な場合があります。その場合は、PR の説明で不足箇所と、合理的な方法ではテストできない理由を説明してください。

```bash
# Run all v3 tests
cd v3
go test ./...

# Run specific package tests
go test ./pkg/application

# Inspect coverage for the packages you changed
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
cd ..
```

統合テストスイート、競合検出、CI と同等の完全なコマンドについては、[テストと継続的インテグレーション](/contributing/testing-ci/)を参照してください。

## 変更を加える

### ブランチを作成する

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### 変更を実装する

1. Go の規約に従って**コードを書く**
2. 新機能の**テストを追加する**
3. 必要に応じて**ドキュメントを更新する**
4. 問題が生じていないことを確認するために**テストを実行する**
5. 明確なメッセージを付けて**変更をコミットする**

### コミットのガイドライン

```bash
# Good commit messages
git commit -m "fix: resolve window focus issue on macOS"
git commit -m "feat: add support for custom window chrome"
git commit -m "docs: improve bindings documentation"

# Use conventional commits:
# - feat: New feature
# - fix: Bug fix
# - docs: Documentation
# - test: Tests
# - refactor: Code refactoring
# - chore: Maintenance
```

### プルリクエストを提出する

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## プルリクエストのガイドライン

### 適切な PR の説明

```markdown
## Description
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- Tested on macOS 14
- Tested on Windows 11
- All tests passing

## Related Issues
Fixes #123
```

### PR チェックリスト

- [ ] コードが Go の規約に従っている
- [ ] テストを追加または更新した
- [ ] ドキュメントを更新した
- [ ] すべてのテストが成功している
- [ ] 破壊的変更がない（または文書化済みである）
- [ ] コミットメッセージが明確である

## コードのガイドライン

### Go コードのスタイル

```go
// ✅ Good: Clear, documented, tested
// ProcessData processes the input data and returns the result.
// It returns an error if the data is invalid.
func ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    result := process(data)
    return result, nil
}

// ❌ Bad: No docs, no error handling
func ProcessData(data string) string {
    return process(data)
}
```

### テスト

```go
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## ドキュメント

### ドキュメントの作成

ドキュメント作業の全手順については、[M-Press でドキュメントを作成する](/contributing/documentation/)に従ってください。固定されたバージョン、実行ファイルと Go のインストール方法、編集可能なソースファイル、ナビゲーション、翻訳チェック、`master`に対する PR の提出について説明しています。

ドキュメントのみの作業では、Wails アプリケーションのビルドもプラットフォーム依存関係も必要ありません。M-Press は、リポジトリのルートで`mpress dev`を使用して既存のサイトをプレビューします。

### ドキュメントのスタイル

- 国際英語の綴りを使用する
- 問題の説明から始める
- 動作する例を提示する
- トラブルシューティングを含める
- 関連コンテンツを相互参照する

## コミュニティ

### サポートを受ける

- **Discord：**[コミュニティに参加する](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions：** 質問する
- **GitHub Issues：** バグを報告する

### 行動規範

互いを尊重し、誰もが参加できるプロフェッショナルな態度を心がけてください。私たちは皆、優れたソフトウェアを一緒に作るためにここにいます。詳しくは[行動規範](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md)を参照してください。

## 貢献者の紹介

コントリビューターは、次の場所で紹介されます。

- リリースノート
- コントリビューター一覧
- GitHub Insights

Wails へのご貢献ありがとうございます！🎉

## 次のステップ

@cards{cols="2"}
◆ GitHubリポジトリ
Wailsリポジトリをご覧ください。

[GitHubで見る →](https://github.com/wailsapp/wails)

---
◆ Discordコミュニティ
コミュニティに参加しましょう。

[Discordに参加 →](https://discord.gg/JDdSxwjhGf)

---
📖 ドキュメント
ドキュメントをお読みください。

[ドキュメントを見る →](/quick-start/why-wails/)

@end
