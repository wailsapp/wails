---
title: "コントリビューション"
description: "Wails に貢献する"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## コントリビューターの皆さん、ようこそ！

Wails へのコントリビューションを歓迎します！バグ修正、機能追加、ドキュメント改善など、皆さんの協力に感謝します。

## 貢献する方法

### 1. 問題を報告する

バグを見つけた場合は、次の情報を添えて[Issue を作成](https://github.com/wailsapp/wails/issues/new)してください。

- 明確な説明
- 再現手順
- 期待される動作と実際の動作
- システム情報
- コード例

### 2. ドキュメントを改善する

修正 PR は、事前の Issue や失敗するコードテストがなくても歓迎します。  
[ドキュメントを修正する](/contributing/documentation/)の手順に従い、M-Press で変更をプレビューして検証してください。

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

### 4. 機能改善を提案する（WEP）

新しい機能や公開されている動作の変更には、Wails Enhancement Proposal（WEP）プロセスを使用します。このプロセスにより、機能開発の透明性が保たれ、承認されたすべての提案に実装担当者が確保されます。機能リクエストの Issue は作成しないでください。

1. 必要に応じて、関心の度合いを確認するため、GitHub Discussions の[Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas)カテゴリまたは[Discord](https://discord.gg/JDdSxwjhGf)でアイデアを共有してください。
2. [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md)を`v3/wep/proposals/<proposal name>/proposal.md`にコピーし、すべてのセクションを記入してください。
3. 提案のみを含み、タイトルを`[WEP] <title>`としたドラフトのプルリクエストを作成してください。この PR が正式な議論の場となります。
4. フィードバックと支持（PR へのコメントや賛成を示すリアクション）を集めてください。議論には少なくとも 2 週間を設け、提案の実装担当者について合意してください。
5. PR をレビュー可能な状態にしてください。最終判断はメンテナーが行います。承認された提案には WEP 番号が割り当てられ、マージされます。

プロセス全体については、[`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)に記載されています。

## はじめに

### フォークしてクローンする

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

テストは変更の一部であり、最後に行う簡単な動作確認ではありません。ユニットテストは、テスト対象のコードの近くに配置してください。複数の入力やエッジケースで同じ動作を確認する場合は、可能な限りテーブル駆動テストを使用してください。失敗時にどのシナリオか分かるよう、各テストケースに名前を付けてください。

新規または変更されたロジックは、完全にカバーすることが推奨されます。PR で追加または変更するコードについて、Go のステートメントカバレッジ 100% を目標にしてください。リポジトリ全体のカバレッジ率を、変更箇所のテストの代わりにしないでください。たとえば、特定の OS でのみ発生するエラーパスや、実機なしでは再現が現実的でない条件など、カバレッジに不足があっても妥当な場合があります。ただし、その不足と、合理的な方法ではテストできない理由を PR の説明に記載してください。

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

統合テストスイート、競合状態の検出、および CI と同等の全コマンドについては、[テストと継続的インテグレーション](/contributing/testing-ci/)を参照してください。

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

1. Go の規約に従って<strong>コードを記述する</strong>
2. 新機能の<strong>テストを追加する</strong>
3. 必要に応じて<strong>ドキュメントを更新する</strong>
4. 問題が発生しないことを確認するため、**テストを実行する**
5. 明確なメッセージを付けて<strong>変更をコミットする</strong>

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
- [ ] テストを追加または更新している
- [ ] ドキュメントを更新している
- [ ] すべてのテストが成功している
- [ ] 破壊的変更がない（ある場合は文書化している）
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

### ドキュメントの執筆

ドキュメントには M-Press を使用しています。`docs/mpress/content/` 配下の  
`.md` ファイルを編集し、リポジトリのルートからプレビューと検証を実行してください。

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

### ドキュメントのスタイル

- 国際英語の綴りを使用する
- 問題の説明から始める
- 動作する例を示す
- トラブルシューティングを含める
- 関連コンテンツへの相互参照を示す

## コミュニティ

### ヘルプを得る

- **Discord：** [コミュニティに参加する](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions：** 質問する
- **GitHub Issues：** バグを報告する

### 行動規範

互いを尊重し、多様性を受け入れ、プロフェッショナルに振る舞ってください。私たちは皆、優れたソフトウェアを共に作るためにここにいます。詳しくは、[行動規範](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md)を参照してください。

## 貢献者の紹介

貢献者は次の場所で紹介されます：

- リリースノート
- 貢献者一覧
- GitHub Insights

Wails に貢献していただき、ありがとうございます！ 🎉

## 次のステップ

@cards{cols="2"}
◆ GitHub リポジトリ
Wails リポジトリにアクセスします。

[GitHub で表示 →](https://github.com/wailsapp/wails)

---
◆ Discord コミュニティ
コミュニティに参加します。

[Discord に参加 →](https://discord.gg/JDdSxwjhGf)

---
📖 ドキュメント
ドキュメントを読みます。

[ドキュメントを見る →](/quick-start/why-wails/)

@end
