---
title: "文件关联"
description: "为 Wails 应用配置文件关联"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

适用平台：<span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

通过文件关联，用户打开特定类型的文件时，应用可以处理这些文件。这对于文本编辑器、图像查看器或任何处理特定文件格式的应用特别有用。本指南介绍如何在 Wails v3 应用中实现文件关联。

## 概述

Wails v3 目前在以下平台和软件包中支持文件关联：

- Windows（NSIS 安装程序包）
- macOS（应用程序包）

## 配置

文件关联在项目的`build`目录下的`config.yml`文件中配置。

### 基本配置

要设置文件关联：

1. 打开`build/config.yml`
2. 在`fileAssociations`部分下添加文件关联
3. 运行`wails3 update build-assets`以更新构建资源
4. 在应用选项中设置`FileAssociations`字段
5. 使用`wails3 package`打包应用

以下是一个配置示例：

```yaml
fileAssociations:
  - ext: myapp
    name: MyApp Document
    description: MyApp Document File
    iconName: myappFileIcon
    role: Editor
  - ext: custom
    name: Custom Format
    description: Custom File Format
    iconName: customFileIcon
    role: Editor
```

### 配置属性

| 属性 | 说明 | 平台 |
| --- | --- | --- |
| ext | 不含开头句点的文件扩展名（例如`txt`） | 全部 |
| name | 文件类型的显示名称 | 全部 |
| description | 文件属性中显示的说明 | Windows |
| iconName | 构建文件夹中图标文件的名称（不含扩展名） | 全部 |
| role | 应用对于此文件类型的角色（例如`Editor`、`Viewer`） | macOS |
| mimeType | 文件的 MIME 类型（例如`image/jpeg`） | macOS |

## 监听文件打开事件

要在应用中处理文件打开事件，可以监听`events.Common.ApplicationOpenedWithFile`事件：

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
    })

    // Listen for files being used to open the application
    app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
        associatedFile := event.Context().Filename()
        app.Dialog.Info().SetMessage("Application opened with file: " + associatedFile).Show()
    })

    // Create your window and run the app...
}

```

## 分步教程

下面以一个简单的文本编辑器为例，逐步设置文件关联：

@steps
### 创建图标
- 为文件类型创建图标（建议尺寸：16x16、32x32、48x48、256x256）
- 将图标保存在项目的`build`文件夹中
- 根据`iconName`配置为图标命名（例如`textFileIcon.png`）

@note{type="tip"}
可以使用`wails3 generate icons`生成所需的图标。运行`wails3 generate icons --help`可了解更多信息。

@end

- 对于 macOS，请在`create:app:bundle:`任务中添加类似`cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources`的复制语句。

### 配置文件关联
编辑`build/config.yml`文件以添加文件关联：

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### 更新构建资源
运行以下命令以更新构建资源：

```bash
wails3 update build-assets
```

### 在应用选项中设置文件关联
在`main.go`文件的应用选项中设置`FileAssociations`字段：

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="为什么应用配置和 config.yml 中都必须指定文件扩展名？"}
在 Windows 上，通过文件关联打开文件时，系统会启动应用，并将文件名作为应用的第一个参数。应用无法得知第一个参数是文件还是命令行参数，因此会使用应用选项中的`FileAssociations`字段来判断第一个参数是否为关联文件。

@end

### 打包应用
使用以下命令打包应用：

```bash
wails3 package
```

打包后的应用将在`bin`目录中创建。随后可以安装并测试该应用。

## 其他说明

- 图标应以 PNG 格式放置在构建文件夹中
- 测试文件关联需要安装打包后的应用

@end
