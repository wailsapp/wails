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
- 開發期間熱重新載入

<strong>完成所需時間：</strong>10分鐘

## 建立專案

@steps
### 產生專案
```bash
wails3 init -n myapp -t vanilla-js
cd myapp
```

本課程明確選用 JavaScript 的 Vanilla + Vite，讓下方的原始檔案與產生的專案一致。未加修飾的`vanilla`範本使用 TypeScript。

@note{type="tip" title="其他範本"}
依照你偏好的框架，試用`-t react`、`-t vue`或`-t svelte`。這些範本預設使用 TypeScript；若要使用純 JavaScript，請選用`-t vanilla-js`或`-t react-js`。執行`wails3 init -l`可查看所有可用的範本，也可以[使用自己的前端框架](/guides/dev/frontend-frameworks/)。

@end

### 瞭解專案結構
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

### 執行應用程式
```bash
wails3 dev
```

@note{type="info" title="首次執行"}
首次執行可能會比預期久，因為需要安裝前端相依套件、產生繫結等。之後執行會快得多。

@end

應用程式開啟後會顯示問候介面。輸入你的姓名並按一下「問候」——Go 後端會處理輸入內容並傳回問候語。

@end

@note{type="tip" title="Windows 11使用者的效能提示"}
建議使用[Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/)儲存專案。Dev Drive 已針對開發工作負載最佳化，與一般 NTFS 磁碟機相比，最多可將建置時間與磁碟存取速度改善30%。

@end

## 運作方式

接著瞭解使其運作的程式碼。

### Go 後端

開啟`greetservice.go`：

```go {title="greetservice.go"}
package main

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}
```

**核心概念：**

1. **服務**——具有匯出方法的 Go 結構
2. **匯出方法**——`Greet`以大寫字母開頭，因此可供前端使用
3. **簡單邏輯**——接收姓名並傳回問候語
4. **型別安全**——已定義輸入與輸出型別

@note{type="tip" title="瞭解服務與繫結"}
**服務**是獨立的 Go 模組，會向前端公開功能。它們只是一般的 Go 結構，具有匯出方法，並註冊於應用程式設定的`Services`欄位中。

**繫結**是自動產生的 TypeScript/JavaScript SDK，可讓前端呼叫這些服務。執行`wails3 dev`或`wails3 build`時，Wails 會分析已註冊的服務，並在`frontend/bindings/`中產生型別安全的繫結。

你可以將服務視為後端 API，將繫結視為與其通訊的用戶端程式庫。

@end

### 註冊服務

開啟`main.go`並找到服務註冊：

```go {title="main.go" highlight="4-6"}
app := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

這會向 Wails 註冊`GreetService`，讓其所有匯出方法都可供前端使用。

### 前端

開啟`frontend/src/main.js`：

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

這個簡化的處理常式會取代入門範本的`window.doGreet`；請保留現有的`showToast`輔助函式及其他啟動程式碼。產生的 HTML 按鈕會呼叫`doGreet`，而`showToast`會顯示其結果。匯入路徑會依照 Go 模組名稱（除非覆寫，否則為`changeme`）。

**核心概念：**

1. **自動產生的繫結**——從產生的程式碼匯入`GreetService`
2. **型別安全的呼叫**——方法名稱與簽章會與 Go 程式碼一致
3. **預設採用非同步**——所有 Go 呼叫都會傳回 Promise
4. **錯誤處理**——來自 Go 的錯誤會由 try/catch 擷取

@note{type="info" title="繫結位於何處？"}
產生的繫結位於`frontend/bindings/`。執行`wails3 dev`或`wails3 build`時，系統會自動建立這些繫結。

**切勿手動編輯這些檔案**——每次建置時都會重新產生。

@end

## 自訂應用程式

接著新增一項功能，以瞭解工作流程。

### 新增「問候多人」功能

@steps
### 將方法新增至 GreetService
將以下內容新增至`greetservice.go`：

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = "Hello " + name + "!"
    }
    return greetings
}
```

### 應用程式將自動重新建置
儲存檔案後，`wails3 dev`會自動重新建置 Go 程式碼並重新啟動應用程式。

@note{type="info" title="自動重新建置"}
Go 程式碼變更會觸發自動重新建置與重新啟動。前端變更則會熱重新載入，不需重新啟動。

@end

### 在前端使用此功能
將以下內容新增至`frontend/src/main.js`：

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

開啟瀏覽器主控台並呼叫`greetMany()`，你將會看到問候語陣列。

@end

## 建置生產版本

準備好發佈應用程式時：

```bash
wails3 build
```

**此操作會：**

- 以最佳化設定編譯 Go 程式碼
- 建置生產版本的前端（最小化）
- 在`bin/`中建立原生可執行檔

@tabs{sync-key="os"}
[Windows]
**輸出：**`bin/myapp.exe`

按兩下即可執行。不需要任何相依項目（WebView2 是 Windows 的一部分）。

[macOS]
**輸出：**`bin/myapp.app`

拖曳至 Applications 資料夾，或按兩下執行。

[Linux]
**輸出：**`bin/myapp`

使用`./bin/myapp`執行，或為啟動器建立`.desktop`檔案。

@end

@note{type="tip" title="跨平台建置"}
想要為其他平台建置嗎？請參閱[跨平台建置 →](/guides/build/cross-platform/)

@end

## 我們學到的內容

**專案結構**

- `main.go`用於 Go 後端
- `frontend/` 用於 UI 程式碼
- `Taskfile.yml` 用於建置工作

**服務**

- 建立具有匯出方法的 Go 結構
- 使用 `application.NewService()` 註冊
- 方法會自動供前端使用

**繫結**

- 自動產生的 TypeScript 型別定義
- 型別安全的函式呼叫
- 預設採用非同步方式（Promise）

**開發工作流程**

- `wails3 dev` 用於熱重新載入
- Go 變更會自動重新建置並重新啟動
- 前端變更會立即熱重新載入

---

<strong>有問題嗎？</strong>歡迎加入 [Discord](https://discord.gg/JDdSxwjhGf)，向社群提問。
