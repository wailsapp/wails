---
title: "你的第一個應用程式"
description: "在10分鐘內建置可運作的 Wails 應用程式"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

我們將建置一個簡單的問候應用程式，示範 Wails 的核心概念：

- 由 Go 後端管理邏輯
- 前端呼叫 Go 函式
- 型別安全的繫結
- 開發期間的熱重新載入

<strong>完成所需時間：</strong>10分鐘

@note{type="tip" title="給 Windows 11使用者的效能提示"}
建議使用[Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/)儲存專案。Dev Drive 已針對開發人員工作負載最佳化；相較於一般 NTFS 磁碟機，可大幅縮短建置時間並提升磁碟存取速度，改善幅度最高可達30%。

@end

## 建立專案

@steps
### 產生專案
```bash
wails3 init -n myapp
cd myapp
```

這會使用預設的 Vanilla + Vite 範本（HTML/CSS/TypeScript 搭配 Vite 打包工具）建立新專案。

@note{type="tip" title="其他範本"}
你可以依偏好的框架嘗試`-t react`、`-t vue`或`-t svelte`。這些範本預設使用 TypeScript；若要使用純 JavaScript，請選擇`-t vanilla-js`或`-t react-js`。 執行`wails3 init -l`即可查看所有可用範本，或 [使用你自己的前端框架](/guides/dev/frontend-frameworks/)。

@end

### 瞭解專案結構
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

### 執行應用程式
```bash
wails3 dev
```

@note{type="info" title="首次執行"}
首次執行時會安裝前端相依套件、產生繫結等，因此耗時可能比預期更久。後續執行會快得多。

@end

應用程式開啟後會顯示問候介面。輸入你的姓名並按一下「問候」；Go 後端會處理輸入內容並傳回問候語。

@end

## 運作方式

接下來瞭解讓它運作的程式碼。

### Go 後端

開啟`greetservice.go`：

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

**核心概念：**

1. **服務**：具有匯出方法的 Go 結構
2. **匯出方法**：`Greet`以大寫字母開頭，因此前端可以使用它
3. **簡單邏輯**：接收姓名並傳回問候語
4. **型別安全**：已定義輸入和輸出型別

@note{type="tip" title="瞭解服務與繫結"}
<strong>服務</strong>是獨立的 Go 模組，可將功能公開給前端。它們只是一般的 Go 結構，具有匯出方法，並在應用程式設定的`Services`欄位中註冊。

<strong>繫結</strong>是自動產生的 TypeScript/JavaScript SDK，讓前端能呼叫這些服務。執行`wails3 dev`或`wails3 build`時，Wails 會分析已註冊的服務，並在`frontend/bindings/`中產生型別安全的繫結。

你可以將服務視為後端 API，將繫結視為與該 API 通訊的用戶端程式庫。

@end

### 註冊服務

開啟`main.go`並找到服務註冊：

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

這會向 Wails 註冊你的`GreetService`，讓前端可以使用其所有匯出方法。

### 前端

開啟`frontend/src/main.js`：

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

**核心概念：**

1. **自動產生的繫結**：從產生的程式碼匯入`GreetService`
2. **型別安全的呼叫**：方法名稱和簽章與你的 Go 程式碼一致
3. **預設為非同步**：所有 Go 呼叫都會傳回 Promise
4. **錯誤處理**：使用 try/catch 擷取來自 Go 的錯誤

@note{type="info" title="繫結位於何處？"}
產生的繫結位於`frontend/bindings/`。執行`wails3 dev`或`wails3 build`時，系統會自動建立這些繫結。

**切勿手動編輯這些檔案**，每次建置時都會重新產生它們。

@end

## 自訂應用程式

接下來新增一項功能，以瞭解整個工作流程。

### 新增「多人問候」功能

@steps
### 將方法新增至 GreetService
將以下內容新增至`greetservice.go`：

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### 應用程式將自動重新建置
儲存檔案後，`wails3 dev`會自動重新建置 Go 程式碼並重新啟動應用程式。

@note{type="info" title="自動重新建置"}
變更 Go 程式碼會觸發自動重新建置和重新啟動。變更前端時則會熱重新載入，不需重新啟動。

@end

### 在前端中使用
將以下內容新增至`frontend/src/main.js`：

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

開啟瀏覽器主控台並呼叫`greetMany()`，你將看到問候語陣列。

@end

## 建置正式版本

準備好發佈應用程式後：

```bash
wails3 build
```

**執行的作業：**

- 使用最佳化設定編譯 Go 程式碼
- 建置正式環境使用的前端（已縮小）
- 在`bin/`中建立原生可執行檔

@tabs{sync-key="os"}
[Windows]
**輸出：**`bin/myapp.exe`

按兩下即可執行。不需要安裝任何相依套件（WebView2 是 Windows 的一部分）。

[macOS]
**輸出：**`bin/myapp.app`

拖曳至「應用程式」資料夾，或按兩下執行。

[Linux]
**輸出：**`bin/myapp`

使用`./bin/myapp`執行，或建立`.desktop`檔案供應用程式啟動器使用。

@end

@note{type="tip" title="跨平台建置"}
想要為其他平台建置嗎？請參閱[跨平台建置 →](/guides/build/cross-platform/)

@end

## 我們學到了什麼

**專案結構**

- `main.go`用於 Go 後端
- `frontend/`用於 UI 程式碼
- `Taskfile.yml`用於建置工作

**服務**

- 建立具有匯出方法的 Go 結構體
- 使用`application.NewService()`註冊
- 方法會自動供前端使用

**繫結**

- 自動產生的 TypeScript 型別定義
- 型別安全的函式呼叫
- 預設採用非同步方式（Promise）

**開發工作流程**

- `wails3 dev`用於熱重新載入
- Go 變更後會自動重新建置並重新啟動
- 前端變更會立即熱重新載入

---

<strong>有問題嗎？</strong>加入[Discord](https://discord.gg/JDdSxwjhGf)並向社群提問。
