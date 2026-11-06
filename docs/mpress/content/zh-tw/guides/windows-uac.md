---
title: "Windows UAC 設定"
description: "為 Windows Wails 應用程式設定使用者帳戶控制（UAC）"
slug: "guides/windows-uac"
sourcePath: "guides/windows-uac.md"
---

適用平台：<span class="mpress-badge mpress-badge-note">Windows</span>

<br/>

Windows 使用者帳戶控制（UAC）決定 Wails 應用程式的執行權限。Wails v3 應用程式預設會在 Windows 資訊清單中包含明確的 UAC 設定，確保在不同電腦上具有一致的行為。

## UAC 執行層級

Windows 應用程式可以透過資訊清單檔案要求不同的執行層級。Wails v3 會自動加入具有預設執行層級的 UAC 設定，而您可以依應用程式的需求自訂此層級。

### 可用的執行層級

| 層級 | 說明 | 使用情境 |
| --- | --- | --- |
| `asInvoker` | 以與父行程相同的權限執行 | 大多數應用程式的預設值 |
| `highestAvailable` | 以使用者可用的最高權限執行 | 可能需要提升存取權限的應用程式 |
| `requireAdministrator` | 一律需要系統管理員權限 | 系統公用程式、安裝程式 |

### 預設設定

Wails v3 應用程式的 Windows 資訊清單中包含預設的 UAC 設定：

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

此設定可確保您的應用程式：

- 以與啟動行程相同的權限執行
- 預設不需要提升權限
- 在不同電腦上皆有一致的行為
- 不會對一般使用者觸發 UAC 提示

## 自訂 UAC 設定

由於 Wails v3 鼓勵使用者自訂建置資產，因此您可以直接編輯 Windows 資訊清單範本來修改 UAC 設定。

### 尋找資訊清單範本

Windows 資訊清單範本位於：

```
build/windows/wails.exe.manifest
```

### 修改執行層級

若要變更執行層級，請編輯`requestedExecutionLevel`元素中的`level`屬性：

```xml {title="build/windows/wails.exe.manifest"}
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

### 範例

#### 標準應用程式（預設）

大多數應用程式應使用預設的`asInvoker`層級：

```xml
<requestedExecutionLevel level="asInvoker" uiAccess="false"/>
```

#### 系統公用程式

可用時需要提升存取權限的應用程式：

```xml
<requestedExecutionLevel level="highestAvailable" uiAccess="false"/>
```

#### 系統管理工具

一律需要系統管理員權限的應用程式：

```xml
<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
```

## UI 存取權

`uiAccess`屬性控制您的應用程式是否可以與權限較高的 UI 元素互動。在大多數情況下，此屬性應維持為`false`。

只有當您的應用程式需要執行下列操作時，才將其設為`true`：

- 將輸入傳送至其他應用程式
- 操控其他應用程式的 UI
- 存取權限較高之行程的 UI 元素

@note{type="caution" title="UI 存取權要求"}
設定`uiAccess="true"`後，您的應用程式必須：

- 使用受信任憑證授權單位所核發的憑證進行數位簽署
- 安裝在安全的位置（Program Files 或 Windows\System32）

@end

## 使用自訂 UAC 設定進行建置

修改資訊清單範本後，請照常建置應用程式：

```bash
wails3 build
```

建置程序會自動將您的自訂 UAC 設定嵌入可執行檔。

## 驗證 UAC 設定

您可以使用`go-winres`工具，確認 UAC 設定是否已正確嵌入：

```bash
go-winres extract --in your-app.exe --out extracted-resources/
```

接著檢查解出的資訊清單檔案，確認其中包含您的 UAC 設定。

@note{type="tip" title="資訊清單持續保留"}
與其他一些框架不同，Wails v3 的 UAC 設定會在編譯期間直接嵌入可執行檔，確保應用程式複製到其他電腦後仍會保留此設定。

@end

## 疑難排解

### 未出現 UAC 提示

如果您已設定`requireAdministrator`，但未看到 UAC 提示：

- 確認資訊清單已正確嵌入可執行檔
- 確認您並非從已提升權限的行程執行應用程式
- 確認資訊清單的語法是有效的 XML

### 應用程式無法啟動

如果您的應用程式在變更 UAC 設定後無法啟動：

- 檢查資訊清單語法是否有 XML 錯誤
- 確認執行層級的值有效
- 嘗試還原為`asInvoker`以釐清問題

### 不同電腦上的行為不一致

如果不同電腦上的 UAC 行為有所差異：

- 確認資訊清單已嵌入可執行檔中（而非外部檔案）
- 確認可執行檔在建置後未遭到修改
- 確認目標電腦已啟用 Windows UAC 設定
