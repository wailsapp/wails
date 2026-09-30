---
title: "建立自訂範本"
description: "如何產生、自訂及託管您自己的 Wails v3 專案範本"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails 隨附一組內建範本，但您也可以建立自己的範本並與社群分享。自訂範本其實就是一個 Git 儲存庫——公開託管後，任何人都能使用一行指令，以它為基礎建立專案骨架。

## 產生範本骨架

`wails3 generate template` 指令會產生可立即開始自訂的範本目錄：

```bash
wails3 generate template -name MyTemplate
```

所有旗標：

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-name` | 範本名稱（必填） | — |
| `-author` | 作者名稱 | — |
| `-description` | 顯示於 CLI 中的簡短說明 | — |
| `-helpurl` | 此範本的說明文件 URL | — |
| `-version` | 初始版本 | `v0.0.1` |
| `-frontend` | 將現有的前端目錄複製到範本中 | — |
| `-dir` | 範本目錄的寫入位置 | 目前目錄 |

使用所有旗標的範例：

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

產生的目錄結構如下：

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="閱讀 NEXTSTEPS.md"}
產生的 `NEXTSTEPS.md` 包含範本各部分的詳細指引。請先閱讀，再開始自訂。發布前請將其刪除——使用您的範本建立的專案中不得出現此檔案。

@end

## 設定範本中繼資料

開啟 `template.yaml`，設定範本的中繼資料：

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

`wailsVersion` 欄位為<strong>必填</strong>，且必須是 `3`。頂端的 `# yaml-language-server` 註解可在 VS Code（搭配[YAML 擴充功能](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)）與 JetBrains IDE 中啟用自動完成和行內驗證。您可以保留或移除該註解；它不會影響執行階段。

## 自訂範本

### 前端

`frontend/` 目錄會原封不動地複製到使用您範本建立的每個專案中。請將預留內容替換為實際的前端：

@tabs
[從零開始]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

依照提示操作，然後安裝相依套件：

```bash
npm install
```

[使用現有專案]
產生範本時傳入 `-frontend`，即可一步將現有前端複製進去：

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

您也可以稍後手動將其複製到 `frontend/` 目錄中。

@end

### 建置工作

`Taskfile.tmpl.yml` 定義建置工作流程。請更新 `install:frontend:deps` 與 `build:frontend` 工作，使其符合您的前端工具鏈：

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Go 應用程式

`main.go.tmpl` 檔案是應用程式進入點。建立專案時，Wails 的範本引擎會處理此檔案，並將 `{{.ProductName}}` 等範本變數替換為使用者提供的值。

若要將其當作真正的 Go 檔案編輯（並獲得 IDE 支援），請暫時將它重新命名為 `main.go`，完成變更後，再於提交前將它重新命名回 `main.go.tmpl`。

#### 範本變數

任何 `.tmpl` 檔案中都可使用以下變數：

| 變數 | 說明 | 範例 |
| --- | --- | --- |
| `{{.ProjectName}}` | 使用者提供的專案名稱 | `"MyApp"` |
| `{{.BinaryName}}` | 二進位檔案名稱 | `"myapp"` |
| `{{.ProductName}}` | 產品顯示名稱 | `"My Application"` |
| `{{.ProductDescription}}` | 產品說明 | `"An awesome application"` |
| `{{.ProductVersion}}` | 產品版本 | `"1.0.0"` |
| `{{.ProductCompany}}` | 公司／作者名稱 | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | 著作權字串 | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | 其他產品註解 | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | 反向 DNS 產品識別碼 | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Go 模組路徑 | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | 用來建立專案的 Wails 版本 | `"3.0.0"` |
| `{{.Typescript}}` | 如果範本名稱以`-ts`結尾，則為`true` | `true` |
| `{{.Opn}}` | 常值`{{` — 在範本內進行逸出 | `{{` |
| `{{.Cls}}` | 常值`}}` — 在範本內進行逸出 | `}}` |

@note{type="tip"}
範本中的任何檔案都可以是`.tmpl`檔案，包括 HTML、JSON 和 YAML 檔案。沒有`.tmpl`副檔名的檔案會原樣複製。

@end

## 在本機測試範本

發布前，請從本機路徑建立專案以測試範本：

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

接著確認專案能正常運作：

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

請檢查：

- 前端熱重新載入功能正常
- 變更 Go 程式碼後，應用程式會重新建置並重新啟動
- `bin/`中的正式環境二進位檔能正確執行

## 發布至 GitHub

@steps
### **為範本建立公開的 GitHub 儲存庫**。儲存庫根目錄必須包含 `template.yaml`。
### **刪除 `NEXTSTEPS.md`** — 此檔案是提供給範本作者的指引，不得出現在使用者以您的範本建立的專案中。
### **提交並推送**範本目錄的內容，作為儲存庫根目錄的內容：
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **使用語意化版本控制標記發行版本**：
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

使用者現在可以從您的範本建立專案：

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="第三方範本警告"}
使用者安裝遠端範本時，Wails 會顯示警告，說明該範本是第三方程式碼，且 Wails 專案不對其內容承擔任何責任。使用者必須明確確認，系統才會建立專案。

身為範本作者，您須對範本中所有程式碼的安全性與正確性負責。

@end

## 最佳實務

- **撰寫清楚的`README.md`** — 使用者建立專案後會看到此內容。請說明如何執行、建置及自訂專案。
- **填寫`helpurl`** — 請連結至您的儲存庫或專屬文件。使用者會在 Wails CLI 範本清單中看到此資訊。
- **固定前端相依套件的版本** — 在`package.json`中固定版本，以免上游更新導致安裝失敗。
- **加上標籤前先進行測試** — 向社群公布前，請從已加上標籤的發行版本建立全新專案。
- **保留`wailsVersion: 3`** — 此欄位會告知 Wails 該範本的目標主要版本。請勿變更。
- **定期更新** — 讓相依套件保持最新，並針對新版 Wails 進行測試。
