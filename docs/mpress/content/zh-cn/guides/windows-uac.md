---
title: "Windows UAC 配置"
description: "为 Windows Wails 应用程序配置用户账户控制（UAC）"
slug: "guides/windows-uac"
sourcePath: "guides/windows-uac.md"
---

适用平台：<span class="mpress-badge mpress-badge-note">Windows</span>

<br/>

Windows 用户账户控制（UAC）决定 Wails 应用程序的执行权限。默认情况下，Wails v3 应用程序会在 Windows 清单中包含明确的 UAC 配置，以确保在不同计算机上的行为一致。

## UAC 执行级别

Windows 应用程序可以通过其清单文件请求不同的执行级别。Wails v3 会自动包含采用默认执行级别的 UAC 配置，你可以根据应用程序的需要进行自定义。

### 可用的执行级别

| 级别 | 说明 | 使用场景 |
| --- | --- | --- |
| `asInvoker` | 以与父进程相同的权限运行 | 大多数应用程序的默认选项 |
| `highestAvailable` | 以用户可获得的最高权限运行 | 可能需要提升访问权限的应用程序 |
| `requireAdministrator` | 始终需要管理员权限 | 系统实用程序、安装程序 |

### 默认配置

Wails v3 应用程序的 Windows 清单中包含默认的 UAC 配置：

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

此配置可确保应用程序：

- 以与启动进程相同的权限运行
- 默认不需要提升权限
- 在不同计算机上的行为一致
- 不会为普通用户触发 UAC 提示

## 自定义 UAC 配置

Wails v3 鼓励用户自定义构建资源，因此你可以直接编辑 Windows 清单模板来修改 UAC 配置。

### 查找清单模板

Windows 清单模板位于：

```
build/windows/wails.exe.manifest
```

### 修改执行级别

要更改执行级别，请编辑`requestedExecutionLevel`元素中的`level`属性：

```xml {title="build/windows/wails.exe.manifest"}
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

### 示例

#### 标准应用程序（默认）

大多数应用程序应使用默认的`asInvoker`级别：

```xml
<requestedExecutionLevel level="asInvoker" uiAccess="false"/>
```

#### 系统实用程序

在可用时需要提升访问权限的应用程序：

```xml
<requestedExecutionLevel level="highestAvailable" uiAccess="false"/>
```

#### 管理工具

始终需要管理员权限的应用程序：

```xml
<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
```

## UI 访问权限

`uiAccess`属性控制应用程序能否与权限更高的 UI 元素交互。在大多数情况下，此属性应保持为`false`。

仅当应用程序需要执行以下操作时，才将其设置为`true`：

- 向其他应用程序发送输入
- 操控其他应用程序的 UI
- 访问权限更高的进程的 UI 元素

@note{type="caution" title="UI 访问要求"}
设置`uiAccess="true"`后，应用程序必须：

- 使用受信任证书颁发机构签发的证书进行数字签名
- 安装在安全位置（Program Files 或 Windows\System32）

@end

## 使用自定义 UAC 设置进行构建

修改清单模板后，照常构建应用程序：

```bash
wails3 build
```

构建过程会自动将自定义 UAC 配置嵌入可执行文件。

## 验证 UAC 配置

可以使用`go-winres`工具验证 UAC 设置是否已正确嵌入：

```bash
go-winres extract --in your-app.exe --out extracted-resources/
```

然后检查提取出的清单文件，确认其中包含 UAC 配置。

@note{type="tip" title="清单持久性"}
与其他一些框架不同，Wails v3 的 UAC 配置会在编译期间直接嵌入可执行文件，从而确保将应用程序复制到其他计算机后，该配置仍然保留。

@end

## 故障排除

### 未出现 UAC 提示

如果已设置`requireAdministrator`，但没有看到 UAC 提示：

- 验证清单是否已正确嵌入可执行文件
- 确认应用程序并非从已提升权限的进程启动
- 确保清单语法是有效的 XML

### 应用程序无法启动

如果应用程序在更改 UAC 设置后无法启动：

- 检查清单语法是否存在 XML 错误
- 确认执行级别值有效
- 尝试恢复为`asInvoker`，以便隔离问题

### 不同计算机上的行为不一致

如果不同计算机上的 UAC 行为存在差异：

- 确保清单已嵌入可执行文件中（而非外置）
- 确认可执行文件在构建后未被修改
- 确认目标计算机已启用 Windows UAC 设置
