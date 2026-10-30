---
title: "使用 M-Press 撰寫文件"
description: "下載並安裝 M-Press、預覽 Wails 文件，以及提交文件貢獻。"
sourcePath: "contributing/documentation.md"
---

Wails v3使用**M-Press**將 Markdown 轉換成這個文件網站。你可以修正錯字、改善說明、新增範例或指南，或修正翻譯。僅修正文件時，不需要事先建立議題，也不需要有失敗的程式碼測試。

你需要使用 Git 準備貢獻，並需要 GitHub 帳號才能提交拉取請求。可下載的 M-Press 執行檔包含文件產生器和圖表轉譯器；預覽文件不需要 Go、Node.js、Wails CLI 或桌面版建置相依項目。儲存庫的其他驗證指令碼需要 Python 3。

## 從文件頁面開始

在側邊欄中開啟**貢獻 → 文件**以進入本指南。若要處理特定頁面，請選取該頁面標頭中的**編輯圖示**，然後選取**在本機編輯此頁面**。複製適用於你作業系統的命令，並在終端機中執行。

如果已安裝 M-Press，該命令會直接使用；否則，其安裝程式會下載官方版本並驗證總和檢查碼。它會準備本機貢獻用工作目錄，並在 M-Press 中開啟所選頁面。執行儲存庫檢查前，請將`mpress version`與下方固定的版本核對：自動安裝程式可能會下載較新的版本。

本機預覽重新建置期間，你可以使用自己的編輯器編輯原始檔。精靈也將**檢閱變更**、**提交變更**、**推送貢獻分支**和**開啟草稿拉取請求**分別列為獨立動作。提交前，請遵循下方的檢查和 PR 指引。產生翻譯並非必要，且可能需要已設定的供應商；編輯和預覽 Markdown 則不需要。

如需可重現的手動設定，請依照下列步驟操作。

## 下載並安裝 M-Press

儲存庫建置指令碼會下載**最新的 M-Press 穩定版**。請在本機使用相同版本：選擇下方其中一種安裝方式、執行`mpress version`，並將結果與發行頁面顯示的版本比較。

### 下載執行檔

開啟[最新的 M-Press 穩定版](https://github.com/leaanthony/mpress/releases/latest)頁面，或下載下方相符的封存檔。ARM64 包括 Apple Silicon；AMD64 則是一般 Intel/AMD PC 使用的架構。

| 電腦 | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

從同一個發行版本下載[checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt)。在下載目錄中計算封存檔的 SHA-256雜湊值，並與`checksums.txt`中該檔名所在行比較：

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

如果你選擇了其他架構，請使用實際下載的檔名。雜湊值相符後，解壓縮並安裝執行檔：

@tabs
[macOS]
在下載目錄中使用 Apple Silicon 封存檔：

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

若使用 Intel Mac，請改用`mpress-darwin-amd64.tar.gz`。將`export PATH`這一行加入`~/.zshrc`，以便在新的終端機中繼續使用該命令。

[Linux]
在下載目錄中：

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

若使用 ARM64，請改用`mpress-linux-arm64.tar.gz`。將`export PATH`這一行加入 shell 啟動檔（例如`~/.bashrc`），以便在新的終端機中繼續使用該命令。

[Windows PowerShell]
在下載目錄中：

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

若使用 ARM64，請改用`mpress-windows-arm64.zip`。將你的使用者`bin`目錄加入 Windows **環境變數**中的使用者**Path**，然後重新開啟終端機，以便在之後的工作階段中繼續使用該命令。
@end

### 使用 Go 安裝

如果您已使用 Go，也可以改為建置最新的穩定版。請查看其`go.mod`，確認所需的工具鏈；Go 的自動工具鏈選擇功能可能會下載該工具鏈。本次文件更新所使用的發行版本需要 Go 1.26.10或更新版本。

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

Go 會將命令安裝至`GOBIN`；若未設定`GOBIN`，則會安裝至`go env GOPATH`下的`bin`目錄。如果找不到`mpress`，請將該目錄加入`PATH`，並重新開啟終端機。

## 取得文件原始碼

[建立 wailsapp/wails 的分支](https://github.com/wailsapp/wails/fork)，然後執行下列命令，並將`YOUR_USERNAME`替換成你的 GitHub 使用者名稱：

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

這會以來源儲存庫的**master**分支為起點建立你的文件分支。請以你自己的描述性分支名稱取代`docs/my-change`。

## 執行即時預覽

在包含`mpress.yaml`的 Wails 儲存庫根目錄中：

```sh
mpress dev
```

在瀏覽器中開啟[http://127.0.0.1:3000](http://127.0.0.1:3000)。若要查看本指南，請前往[http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/)。編輯期間請讓命令持續執行；儲存原始檔變更後，預覽會重新建置。在執行下方的驗證命令前，請按**Ctrl+C**停止該命令。

儲存庫已包含 M-Press 設定，不需要初始化新網站。若要使用其他連接埠，請執行`mpress dev --port 3001`，並開啟對應的 localhost 位址。

## 編輯或新增頁面

英文來源頁面位於`docs/mpress/content/`。翻譯位於語言資料夾下，例如`fr/`和`id/`。編輯圖示中的貢獻指示會指出目前頁面的原始檔。

若要修改現有頁面，請編輯其 Markdown，並保留頂端的 YAML 中繼資料。請使用標題、段落、清單、連結及圍欄式程式碼區塊。請保留成對的 M-Press 元件，例如`@note`／`@end`和`@tabs`／`@end`。[M-Press 編寫文件](https://m-press.me/authoring/)說明了支援的語法。

若要新增指南，請建立如下的檔案：`docs/mpress/content/guides/my-guide.md`

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

將頁面加入`docs/mpress/content/_nav.yaml`中的適當區段，除非該區段已使用`autogenerate`。確認讀者能透過導覽和相關頁面的連結前往該頁面。新的英文頁面也必須有完整的各語系副本，涵蓋率檢查才能通過；請與維護者或翻譯貢獻者協調。

產生的輸出位於`docs/mpress/site/`。請提交來源頁面，以及預期的導覽或靜態資產變更；產生的網站檔案和本機快取不得納入 PR。

## 檢查你的貢獻

在儲存庫根目錄中：

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

如果 Windows 上的 Python 3是以該方式安裝，請使用`python`。如果已安裝[Task](https://taskfile.dev/)，`task docs:build`會使用 PATH 中的`mpress`執行相同的檢查。

翻譯檢查會在翻譯檔案缺失、內容為空或失去對應來源時失敗。過時或損壞的翻譯段落預設會以警告形式回報；`--strict`也會讓這些檢查結果導致失敗。修改既有英文內容後，即使翻譯同步尚未完成，仍可提交。若只修正一種語言的翻譯，則可以只變更該語言。請保持命令、API 名稱、連結、程式碼範例及圖表連線不變；不要將英文預留文字當作翻譯內容。

請在瀏覽器中閱讀變更後的頁面，並執行所有經您修改的程式碼範例。如果檢查在未修改的頁面上失敗，請在 PR 中附上輸出與相關背景資訊，而不要降低檢查標準。

## 提交 Pull Request

檢查差異、提交原始碼變更，並將您的貢獻分支推送至您的分支版本。請在`wailsapp/wails`中針對**master**開啟 PR。說明變更內容、附上受影響頁面的連結，並列出您執行過的檢查。若有版面配置變更，請提供變更前後的螢幕截圖；若有程式碼範例變更，請提供平台與版本詳細資訊。

正式環境網站是靜態網站；本機預覽與驗證不需要部署認證或翻譯服務供應商。維護者會在審查後負責發布。

如需瞭解程式碼貢獻與功能提案，請參閱[為 Wails 做出貢獻](/contributing/)。如需瞭解框架內部機制，請參閱[技術概觀](/contributing/overview/)。
