---
title: "BMEditor"
description: "功能齊全的程式碼編輯器，打開即用，需要更多時還能擴充"
slug: "community/showcase/bmeditor"
sourcePath: "community/showcase/bmeditor.md"
---

![開啟 Go 專案的 BMEditor](/assets/showcase-images/bmeditor-editor.webp)

![偵錯 Go 程式，每個變數的值顯示在程式碼旁邊](/assets/showcase-images/bmeditor-debug.webp)

![Docker 面板，容器依 compose 專案分組](/assets/showcase-images/bmeditor-docker.webp)

![擴充功能視窗，BMEditor 自家的擴充功能隨時可以安裝](/assets/showcase-images/bmeditor-extensions.webp)

![側邊面板中的 Claude Code，它的修改正在編輯器中等待審閱](/assets/showcase-images/bmeditor-claude.webp)

BMEditor（Bare Minimum Editor (batteries included)）是一款適用於 Windows 與 Linux 的桌面程式碼編輯器，一打開就能開始工作。偵錯器、終端機、Docker、專案中的 Git 變更，以及 Markdown、HTML 與 OpenAPI 的即時預覽都內建在編輯器裡，不必先去找工具或做設定。Go 開箱即用，編輯器本身也可擴充，其他語言與工具都可以透過擴充功能加入。

在同一個視窗中，你可以開啟 WSL 裡的資料夾、透過 SSH 連線的伺服器，或 Docker 容器內部，另一端不需要安裝任何東西。如果你使用 Claude Code 或 GitHub Copilot CLI，它會在側邊面板中執行，它改過的檔案會以審閱的形式回到你手上，每一處變更都可以接受或拒絕。版面配置與快速鍵沿用 VS Code，第一分鐘就會覺得熟悉。

它以 Wails 3、Go 與 Vue 打造，可免費使用。

[官方網站](https://cn8001.github.io/BMEditor-Releases/) | [GitHub 發行版](https://github.com/Cn8001/BMEditor-Releases/releases)
