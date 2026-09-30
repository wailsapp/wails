---
title: "参与贡献"
description: "为 Wails 做贡献"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## 欢迎各位贡献者！

欢迎为 Wails 做贡献！无论是修复错误、添加功能还是改进文档，我们都感谢你的帮助。

## 贡献方式

### 1. 报告问题

发现了错误？请[创建议题](https://github.com/wailsapp/wails/issues/new)，并提供以下信息：

- 清晰的说明
- 复现步骤
- 预期行为与实际行为
- 系统信息
- 代码示例

### 2. 改进文档

无需预先创建议题或提供失败的代码测试，即可提交文档修正 PR。  
请按照[修正文档](/contributing/documentation/)中的说明，使用 M-Press 预览并验证更改。

我们始终欢迎以下文档改进：

- 修正拼写错误和其他错误
- 添加示例
- 澄清说明
- 翻译内容

### 3. 提交代码

通过拉取请求贡献代码：

- 错误修复
- 新功能
- 性能改进
- 测试

### 4. 提出增强建议（WEP）

新增功能和公共行为变更须遵循 Wails 增强提案（WEP）流程。该流程让功能开发保持透明，并确保每个获批提案都有实施者。请勿创建功能请求议题。

1. 你可以选择先在 GitHub Discussions 的[创意](https://github.com/wailsapp/wails/discussions/categories/ideas)类别或[Discord](https://discord.gg/JDdSxwjhGf)上提出想法，以了解大家的兴趣。
2. 将[`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md)复制到`v3/wep/proposals/<proposal name>/proposal.md`，并填写所有章节。
3. 创建一个标题为`[WEP] <title>`的草稿拉取请求，其中仅包含提案。该 PR 是讨论提案的正式场所。
4. 收集反馈和支持（PR 上的评论和点赞回应）。至少留出两周时间进行讨论，并就提案的实施者达成一致。
5. 将 PR 标记为可供审查。维护者将作出最终决定：获批的提案会分配 WEP 编号并合并。

完整流程记录在[`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)中。

## 入门

### 创建派生仓库并克隆

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### 从源代码构建

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### 运行测试

测试是更改的一部分，而不是最后才执行的冒烟检查。请将单元测试放在其所测试代码的旁边；如果要用多个输入或边界情况检查同一种行为，应优先采用表驱动测试。请为每个测试用例命名，以便在失败时说明具体场景。

新增和变更的逻辑应得到完整覆盖。对于 PR 新增或更改的代码，Go 语句覆盖率应以100% 为目标；不要用整个仓库的覆盖率百分比代替对本次更改的测试。存在覆盖缺口可能是合理的，例如仅在特定操作系统上出现的错误路径，或不使用真实硬件便难以复现的条件；但请在 PR 说明中解释该缺口，以及无法合理测试它的原因。

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

有关集成测试套件、竞态检测以及与完整 CI 等效的命令，请参阅[测试与持续集成](/contributing/testing-ci/)。

## 进行更改

### 创建分支

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### 进行更改

1. 按照 Go 约定<strong>编写代码</strong>
2. 为新功能<strong>添加测试</strong>
3. 如有需要，**更新文档**
4. **运行测试**，确保没有破坏任何功能
5. 使用清晰的消息<strong>提交更改</strong>

### 提交准则

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

### 提交拉取请求

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## 拉取请求准则

### 良好的 PR 说明

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

### PR 检查清单

- [ ] 代码遵循 Go 约定
- [ ] 已添加或更新测试
- [ ] 已更新文档
- [ ] 所有测试均已通过
- [ ] 没有破坏性变更（或已记录相关变更）
- [ ] 提交消息清晰明确

## 代码准则

### Go 代码风格

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

### 测试

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

## 文档

### 编写文档

文档使用 M-Press。编辑`docs/mpress/content/`目录下的`.md`文件，然后从仓库根目录预览并验证：

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

### 文档风格

- 使用国际英语拼写
- 从问题入手
- 提供可运行的示例
- 包含故障排除说明
- 交叉引用相关内容

## 社区

### 获取帮助

- **Discord：**[加入我们的社区](https://discord.gg/JDdSxwjhGf)
- <strong>GitHub Discussions：</strong>提问
- <strong>GitHub Issues：</strong>报告缺陷

### 行为准则

请保持尊重、包容和专业。我们齐聚于此，是为了共同打造出色的软件。详情请参阅[行为准则](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md)。

## 贡献者致谢

贡献者将在以下位置获得致谢：

- 发行说明
- 贡献者名单
- GitHub 洞察

感谢你为 Wails 做出贡献！🎉

## 后续步骤

@cards{cols="2"}
◆ GitHub 仓库
访问 Wails 仓库。

[在 GitHub 上查看 →](https://github.com/wailsapp/wails)

---
◆ Discord 社区
加入社区。

[加入 Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 文档
阅读文档。

[浏览文档 →](/quick-start/why-wails/)

@end
