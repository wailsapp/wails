---
title: "更正文档"
description: "使用 M-Press 提交用于更正 Wails v3 文档的 PR。"
sourcePath: "contributing/documentation.md"
---

欢迎提交更正 PR。你可以修正拼写错误、失效链接、过时示例、不清晰的说明或翻译。对于仅更正文档的 PR，无需先创建 issue，也无需提供失败的代码测试。

## 在本地预览

复刻 [wailsapp/wails](https://github.com/wailsapp/wails/fork)，克隆你的复刻仓库，然后从 `master` 创建分支。

安装固定版本的文档生成器：

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

你也可以从[M-Press v1.0.17发布版本](https://github.com/leaanthony/mpress/releases/tag/v1.0.17)下载经过验证的二进制文件。

在 Wails 仓库根目录中运行：

```sh
mpress version
mpress dev
```

编辑 `docs/mpress/content/` 中的 `.md` 源文件。英语是默认语言，其文件直接位于该目录中。现有翻译位于 `fr/` 和 `id/` 等语言文件夹中。保存时，预览会自动重新构建。

保留每个页面顶部的元数据块以及配对的 `@...` / `@end` 组件。普通段落、标题、列表和围栏代码均可作为文本编辑。请勿编辑 `docs/mpress/site/` 中的生成文件。

## 检查更正内容

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

在浏览器中检查已更改的页面，并运行你修改过的所有代码示例。更正翻译时，请将更改后的完整段落与英文原文进行比较。确保命令、API 名称、链接、代码示例和图表连接保持不变。使用自然的技术语言，保留要求和注意事项，并翻译可见的图表标签、导航标签、图像描述及正文。

每种已发布语言都必须完整翻译每个英文页面。请勿使用英文占位内容或回退页面。如果只更正某一种翻译，可以仅更改该语言。如果更改了英文原文的含义，请更新其他已发布语言中对应的页面；无需重新生成无关页面。

## 提交拉取请求

向 `master` 发起 PR。描述原有问题，说明你的更正，并列出已执行的检查。对于可见的布局更改，请提供屏幕截图；对于已更改的代码示例，请提供平台和版本详情。

不需要 Cloudflare 凭据或私有服务。公开 PR 的检查会在无需部署凭据的情况下构建并验证静态站点。

有关代码更改和功能提案，请参阅[为 Wails 做贡献](/contributing/)。有关内部原理，请参阅[技术概述](/contributing/overview/)。
