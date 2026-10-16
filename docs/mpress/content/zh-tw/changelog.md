---
title: "變更日誌"
description: "Wails v3 的版本歷史與版本資訊"
slug: "changelog"
sourcePath: "changelog.md"
---

圖例：

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- 此專案的所有重大變更都將記錄於此檔案中。

此格式以[維護變更日誌](https://keepachangelog.com/en/1.0.0/)為基礎， 且此專案遵循[語意化版本](https://semver.org/spec/v2.0.0.html)。

- `Added`表示新功能。
- `Changed`表示現有功能的變更。
- `Deprecated`表示即將移除的功能。
- `Removed`表示現已移除的功能。
- `Fixed`表示任何錯誤修正。
- `Security`表示安全漏洞。

_/

/_   * 請勿更新此檔案 *   更新應新增至`v3/UNRELEASED_CHANGELOG.md`   謝謝！ _/

## [尚未發布]

## v3.0.0-beta.21 - 2026-09-13

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/6116)中使用 M-Press 提供 Wails v3文件，由 @leaanthony 貢獻

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/6118)中解析 MPD 前置資料中的 JSON slug 值以產生變更日誌，由 @leaanthony 修正
- 在[PR](https://github.com/wailsapp/wails/pull/6080)中，更新程式會清除輔助程式的環境變數，並在備份失敗後重新啟動原始目標，由 @cnmax 修正
- 在[PR](https://github.com/wailsapp/wails/pull/6098)中，於 App.Run 期間啟動預設訊號處理常式，由 @leaanthony 修正
- 在[PR](https://github.com/wailsapp/wails/pull/6112)中，Windows 選單可處理 nil 選單、釋放被取代的資源，並重新繪製選單列，由 @taliesin-ai 修正
- 在[PR](https://github.com/wailsapp/wails/pull/6115)中，為使用共用 YAML 組態的新建專案恢復 MSIX 封裝功能，由 @leaanthony 修正
- 修正泛型模型建立函式參照後續宣告的輔助項目時，產生的 JavaScript 與 TypeScript 繫結無法載入的問題，並防止建立相互依賴的泛型模型時發生堆疊溢位（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 已變更

- 由 @01xR4in 在[PR](https://github.com/wailsapp/wails/pull/6082)中，將 Clave 展示範例的連結更新為目前的網站與儲存庫

## 已修正

- 取消已中止的 Windows 資產請求（包括 worker 請求），同時在導覽期間保留 keepalive 處理常式。在 Apple 平台上，透過應用程式包裝函式轉送原生請求上下文。（#5963、#5969）
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6094)中加入重試機制，在相互競爭的推送之間保留變更日誌項目
- 由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6031)中，移除參照模組從未隨附之二進位檔的嵌入項目，修正`go mod vendor`在所有平台上皆因`pattern arm64/WebView2Loader.dll: no matching files found`而失敗的問題，並修正[#5782](https://github.com/wailsapp/wails/issues/5782)與[#5376](https://github.com/wailsapp/wails/issues/5376)

## 已移除

- 由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6031)中移除原生 WebView2 載入器支援，改由純 Go 載入器取代。這會移除內嵌的`WebView2Loader.dll`二進位檔與`github.com/jchv/go-winloader`相依性。`native_webview2loader`建置標籤仍會被接受且不再引發錯誤，但對 v3 建置沒有任何作用
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6097)中，從 macOS API 指南移除未使用的建置標籤與 FPS 選項

## v3.0.0-beta.19 - 2026-09-09

## 已新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6087)中，使用建置標籤限制私有 macOS API，僅供選擇加入時使用——請參閱[文件](https://v3.wails.io/features/browser/integration)、[文件](https://v3.wails.io/features/environment/info)、[文件](https://v3.wails.io/features/windows/basics)、[文件](https://v3.wails.io/features/windows/frameless)、[文件](https://v3.wails.io/features/windows/notch-windows)、[文件](https://v3.wails.io/features/windows/options)、[文件](https://v3.wails.io/guides/build/macos)、[文件](https://v3.wails.io/guides/build/private-macos-apis)及[文件](https://v3.wails.io/reference/overview)

## 已修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6091)中，以 HTTP 413拒絕超過64 MiB 的執行階段請求

## 安全性

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6092)中，以權杖驗證強化 MCP 來源與遠端存取的安全性

## v3.0.0-beta.18 - 2026-09-08

## 已修正

- 由 @4RH1T3CT0R7 在[PR](https://github.com/wailsapp/wails/pull/6083)中使用指標接收器，修正 Linux 與 Darwin 上的 Calloc 記憶體洩漏

## v3.0.0-beta.17 - 2026-09-06

## 已修正

- Windows：WebResourceRequested 處理常式中失敗或為 nil 的`GetRequest`不再終止程序（`log.Fatal`／nil 解除參照 panic）；由 @midagedev 在[PR](https://github.com/wailsapp/wails/pull/6006)中改為捨棄該請求並寫入記錄

## v3.0.0-beta.16 - 2026-08-29

## 已變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6029)中，改為在新的終端機視窗提示輸入公證密碼

## 已修正

- 由 @ChewbaccaCookie 在[PR](https://github.com/wailsapp/wails/pull/5919)中，正確處理 macOS 上的系統匣點選類型
- 由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6041)中，讓 CI 在更新前移除未使用的 Microsoft apt 軟體庫

## v3.0.0-beta.15 - 2026-08-27

## 已修正

- 由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6043)中，將 WebView2 嵌入逾時時間延長至60秒

## v3.0.0-beta.14 - 2026-08-26

## 已修正

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6032)中，正確命名 macOS 上的 Control 與字母鍵組合鍵輸入
- 由 @nik9play 在[PR](https://github.com/wailsapp/wails/pull/6016)中，修正 ICO 系統匣圖示，並在 Windows 上跟隨工作列佈景主題

## v3.0.0-beta.13 - 2026-08-25

## 已修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6026)中，讓 macOS 在執行強制回應迴圈時仍持續處理主執行緒工作
- 由 @mortenolsrud 在[PR](https://github.com/wailsapp/wails/pull/5923)中，讓行動裝置的安全儲存空間可回報失敗，並在失敗時採取封閉策略
- 由 @archy-rock3t-cloud 在[PR](https://github.com/wailsapp/wails/pull/5999)中，讓應用程式事件掛鉤即使未註冊監聽器也會執行
- 由 @haoku123 在[PR](https://github.com/wailsapp/wails/pull/6023)中，修正註解與本地化文件中的錯字
- 由 @4RH1T3CT0R7 在[PR](https://github.com/wailsapp/wails/pull/6025)中，移除提交於`v3/examples`下的預先編譯 macOS 二進位檔

## v3.0.0-beta.12 - 2026-08-21

## 已新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6010)中，新增具生命週期與遙測範例的 macOS 瀏海通知視窗——請參閱[文件](https://v3.wails.io/features/windows/notch-windows)
- 新增 macOS NSPanel 視窗支援，包含新選項與原生整合，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/6008)中實作——請參閱[文件](https://v3.wails.io/features/windows/options)

## 修正

- 防止並行呼叫 SQLite Prepare 時發生停滯，見 @archy-rock3t-cloud 的 [PR](https://github.com/wailsapp/wails/pull/5998)

## v3.0.0-beta.11 - 2026-08-20

## 移除

- 從文件中移除已過時的實作追蹤器，見 @leaanthony 的 [PR](https://github.com/wailsapp/wails/pull/6005)

## v3.0.0-beta.10 - 2026-08-19

## 修正

- 修正 GTK4 Linux 主程式遺漏自訂通訊協定與檔案關聯啟動引數的問題，見 @midagedev 的 [PR](https://github.com/wailsapp/wails/pull/6000)
- 在變更記錄驗證中正確處理已刪除的行與同來源修正，見 @taliesin-ai 的 [PR](https://github.com/wailsapp/wails/pull/5993)

## v3.0.0-beta.9 - 2026-08-16

## 新增

- 新增安全的 wails3 mcp 伺服器，用於代理程式輔助的專案管理，見 @leaanthony 的 [PR](https://github.com/wailsapp/wails/pull/5896)
- 新增繫結中模型的文件——請參閱 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5988) 中提供的[文件](https://v3.wails.io/features/bindings/models)
- 支援在不可變式 Linux 系統上透過 rpm-ostree 安裝，見 @leaanthony 的 [PR](https://github.com/wailsapp/wails/pull/5987)
- 新增原生的每週 Star 歷史圖表產生與發布功能——請參閱 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5986) 中提供的[文件](https://v3.wails.io/credits)、[文件](https://v3.wails.io/de/credits)、[文件](https://v3.wails.io/fr/credits)、[文件](https://v3.wails.io/id/credits)、[文件](https://v3.wails.io/ja/credits)、[文件](https://v3.wails.io/ko/credits)、[文件](https://v3.wails.io/pt/credits)、[文件](https://v3.wails.io/ru/credits)、[文件](https://v3.wails.io/zh-cn/credits)及[文件](https://v3.wails.io/zh-tw/credits)
- 新增僅適用於 Darwin 的 mac 套件，用於解析應用程式套件資源——請參閱 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5965) 中提供的[文件](https://v3.wails.io/guides/build/macos)
- 新增 Condui 展示頁面與索引項目——請參閱 @mgueregath 在 [PR](https://github.com/wailsapp/wails/pull/5962) 中提供的[文件](https://v3.wails.io/community/showcase/condui)及[文件](https://v3.wails.io/community/showcase)
- 新增 Redis Viewer 展示頁面，包含螢幕擷取畫面與專案連結——請參閱 @redisviewer 在 [PR](https://github.com/wailsapp/wails/pull/5984) 中提供的[文件](https://v3.wails.io/community/showcase)及[文件](https://v3.wails.io/community/showcase/redisviewer)

## 變更

- 將 Linux 的 GTK 應用程式旗標更新為 G<em>APPLICATION</em>NON_UNIQUE，見 @overlordtm 的 [PR](https://github.com/wailsapp/wails/pull/5971)
- 將找不到視窗事件的記錄層級從警告改為偵錯，見 @julianstorer 的 [PR](https://github.com/wailsapp/wails/pull/5914)

## 修正

- 允許已註冊的 macOS 鍵盤快速鍵優先於 WebView 處理，見 @julianstorer 的 [PR](https://github.com/wailsapp/wails/pull/5902)
- 修復文件側邊欄中失效的連結，見 @northes 的 [PR](https://github.com/wailsapp/wails/pull/5937)
- WebKit 中止相符的自訂 URL scheme 工作時，取消 macOS 與 iOS 的資產請求上下文（#5963）
- 處理 WindowSetFullscreenButtonEnabled 訊息，見 @archy-rock3t-cloud 的 [PR](https://github.com/wailsapp/wails/pull/5976)
- 在 preact-ts 範本中匯入 Fragment，以解決建置失敗問題，見 @haoku123 的 [PR](https://github.com/wailsapp/wails/pull/5979)
- 防止舊版僅服務型 GTK3 應用程式在作用中視窗或顯示器可用之前執行螢幕探索時當機（#5966）
- 即使未發布的變更記錄為空，也允許明確指定版本的發布執行繼續進行（#5977）

## 安全性

- 將網站的 nanoid 鎖定檔更新至已修補的 3.3.18，以解決安全性公告所述問題，見 @taliesin-ai 的 [PR](https://github.com/wailsapp/wails/pull/5985)

## v3.0.0-beta.8 - 2026-08-12

## 新增

- 為自動變更記錄項目新增文件 URL 產生功能，見 @taliesin-ai 的 [PR](https://github.com/wailsapp/wails/pull/5957)
- 新增 Streams：採用 WebSocket 程式設計模型，在 Go 與 JavaScript 之間提供雙向位元組串流，且不使用監聽通訊端。在 Go 中以 `app.HandleStream(name, handler)` 宣告串流，並從前端以 `Stream(name)` 連線；後者會傳回形狀如 `WebSocket` 的物件。Go→JS 透過資產伺服器，使用每個視窗一個持續等待的輪詢來傳輸；JS→Go 則使用一般的 POST。過程中不會繫結任何 TCP 連接埠，也不會有任何資料通過 `evaluateJavaScript`。在伺服器建置（`-tags server`）中，改以真正的 WebSocket 提供相同的處理常式，因此不同建置可使用完全相同的應用程式碼。作者：@leaanthony
- 將 mailbox 變更記錄項目移至 Unreleased，見 @leaanthony 的 [PR](https://github.com/wailsapp/wails/pull/5935)

## 變更

- 更新文件側邊欄自動產生功能與部落格作者型別推導，見 @leaanthony 的 [PR](https://github.com/wailsapp/wails/pull/5938)

## 修正

- WebView2 初始化改用截止期限與訊息泵浦，見 @leaanthony 的 [PR](https://github.com/wailsapp/wails/pull/5952)
- 除非明確選擇啟用，否則在 CI 中略過 WebView2 Cookie 測試；測試執行時會鎖定於目前的作業系統執行緒，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5951)
- Windows 選單建構器會還原子選單父項目的命令 ID，見 @gilad-ch 的 [PR](https://github.com/wailsapp/wails/pull/5944)
- 使官方交叉編譯映像檔符合 GTK 4.14+ Linux 支援基準（#5928）
- 設定 iOS Xcode 專案以保留繼承的連結器旗標，並加入 -ObjC，見 @mortenolsrud 的 [PR](https://github.com/wailsapp/wails/pull/5915)
- 修正大型前端中 `wails3 dev` 資產 Proxy 過度頻繁建立及關閉 TCP 連線的問題；此問題可能耗盡主機的暫時性連接埠，並導致不相關的程序發生 `EADDRNOTAVAIL` 錯誤
- 將各視窗的事件 JavaScript 排入佇列，以便依序分派並提供背壓，見 @leaanthony 的 [PR](https://github.com/wailsapp/wails/pull/5934)

## 移除

- 移除桌面二進位檔發布管線：v3 發布僅建立標籤，並使用 `go install` 安裝 `wails3` CLI。同時刪除 `release-v3.yml`，以及每晚作業中分派該管線的步驟，見 @leaanthony 的 [PR](https://github.com/wailsapp/wails/pull/5946)

## v3.0.0-beta.7 - 2026-08-11

## 新增

- 新增 macOS 自動播放偏好設定，可停用媒體播放需要使用者操作的要求，見 [PR](https://github.com/wailsapp/wails/pull/5512)，作者：@Eyalm321
- 將信箱變更記錄項目移至「未發布」，見 [PR](https://github.com/wailsapp/wails/pull/5935)，作者：@leaanthony

## 變更

- macOS 縮放動畫改用 CADisplayLink 或 NSTimer，以提升流暢度，見 [PR](https://github.com/wailsapp/wails/pull/5945)，作者：@savely-krasovsky

## 修正

- 設定 iOS Xcode 專案以保留繼承的連結器旗標，並加入 -ObjC，見 [PR](https://github.com/wailsapp/wails/pull/5915)，作者：@mortenolsrud
- 修正大型前端中 `wails3 dev` 資源代理伺服器過度頻繁地建立與關閉 TCP 連線的問題。此問題可能耗盡主機的暫時連接埠，導致不相關的處理程序因 `EADDRNOTAVAIL` 而失敗
- 將各視窗的事件 JavaScript 排入佇列，以依序分派並提供背壓機制，見 [PR](https://github.com/wailsapp/wails/pull/5934)，作者：@leaanthony

### 新增

- 實作通用非同步 FIFO 信箱，以依序傳遞事件，見 [PR](https://github.com/wailsapp/wails/pull/5851)，作者：@savely-krasovsky 與 @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## 新增

- 實作有容量上限的主機端儲存空間，用於儲存過大的事件並依序傳遞 JavaScript，見 [PR](https://github.com/wailsapp/wails/pull/5930)，作者：@leaanthony
- 實作 macOS Dock 圖示彈跳，以支援視窗閃爍通知，見 [PR](https://github.com/wailsapp/wails/pull/5921)，作者：@julianstorer

## 修正

- 資源伺服器在清空緩衝區時會保留內容類型偵測器錯誤和尚未寫入的前置內容，見 [PR](https://github.com/wailsapp/wails/pull/5931)，作者：@leaanthony
- 防止 macOS 應用程式從 Wails 回呼取代應用程式選單時當機
- 修正 Windows 10 1809／Windows Server 2019（組建 17763）上的原生選單無法閱讀問題。深色模式的 uxtheme 匯出受到組建 18334 條件限制，因此應用程式層級的深色模式選擇加入機制從未在這些主機上執行：選單背景已繪製為深色，但 Windows 仍以淺色佈景主題繪製選單文字，導致深色背景上出現深色文字。這些序號從 17763 起即已存在，因此條件限制現已與其相符。
- 修正 `w32.GetStockObject` 呼叫 `GetDeviceCaps` 而非 `GetStockObject` 的問題；此問題會使其針對每個預設物件都傳回 0。
- 改善 WebView2 引導安裝程式下載錯誤的處理與回報，見[PR](https://github.com/wailsapp/wails/pull/5924)，作者：@jannskiee

## v3.0.0-beta.5 - 2026-08-07

## 修正

- macOS 應用程式啟用作業現在僅針對一般應用程式遵循啟用原則，見 [PR](https://github.com/wailsapp/wails/pull/5897)，作者：@julianstorer
- 在 Linux 組建中防護尚未初始化的 GTK 視窗，見 [PR](https://github.com/wailsapp/wails/pull/5898)，作者：@julianstorer
- 在載入 URL 前，為 Linux WebKit 視窗明確設定不透明背景色彩，見 [PR](https://github.com/wailsapp/wails/pull/5899)，作者：@julianstorer

## v3.0.0-beta.4 - 2026-08-05

## 變更

- Android 建置工作預設以 arm64 為目標，而 deploy-emulator 會選擇主機架構，見 [PR](https://github.com/wailsapp/wails/pull/5890)，作者：@mortenolsrud

## 修正

- 拖曳期間保留 macOS 視窗的縮放狀態並減少動態效果，見 [PR](https://github.com/wailsapp/wails/pull/5900)，作者：@leaanthony
- 將 `!server` 加入 `webview_window_windows_nonclient.go` 建置條件，以修正 Windows 伺服器模式組建

## v3.0.0-beta.3 - 2026-08-03

## 新增

- 在實作詳細資料中記錄第 10 階段的 beta 驗證已完成，見 [PR](https://github.com/wailsapp/wails/pull/5881)，作者：@leaanthony

## 修正

- 將視窗控制代碼傳遞給 Windows 深色模式 API，並驗證引數，見 [PR](https://github.com/wailsapp/wails/pull/5877)，作者：@leaanthony
- 集中處理 macOS 無框線視窗標題列按鈕的狀態解析，見 [PR](https://github.com/wailsapp/wails/pull/5870)，作者：@taliesin-ai
- 防止 Windows 應用程式要求使用深色模式、但 Windows 應用程式佈景主題為淺色時，原生選單文字無法閱讀。在 Windows 能夠呈現深色選單文字前，選單現在會使用相符的淺色原生背景。
- 修正 Windows 10 1809／Windows Server 2019（組建 17763）上的原生選單無法閱讀問題。深色模式的 uxtheme 匯出受到組建 18334 條件限制，因此應用程式層級的深色模式選擇加入機制從未在這些主機上執行：選單背景已繪製為深色，但 Windows 仍以淺色佈景主題繪製選單文字，導致深色背景上出現深色文字。這些序號從 17763 起即已存在，因此條件限制現已與其相符。

## v3.0.0-beta.2 - 2026-08-02

## 變更

- 將 v3 從 alpha 升級為 beta
- 記錄系統匣的智慧型預設值與快顯視窗自動隱藏行為，並為點擊處理常式的選擇加入迴歸測試涵蓋範圍（#5840）。
- GitHub 更新程式現在預設排除 Windows 安裝程式資源，見 [PR](https://github.com/wailsapp/wails/pull/5861)，作者：@leaanthony
- 新增 macOS 無框線視窗的圓角、方角及自訂圓角半徑支援，見 [PR](https://github.com/wailsapp/wails/pull/5866)，作者：@leaanthony

## 修正

- 回報即時 GTK4 視窗大小，並從已設定的表面發出調整大小、最大化、最小化及全螢幕狀態事件（#5830）。
- 修正在 fetch 要求中傳送 Blob 或 FormData 時導致 Linux WebKit 當機的問題，見 [PR](https://github.com/wailsapp/wails/pull/5854)，作者：@taliesin-ai
- 缺少 Blob/FormData 標頭時，Fetch 相容層會傳遞 undefined，見 [PR](https://github.com/wailsapp/wails/pull/5865)，作者：@leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## 新增

## 變更

- 新增 macOS 無框線視窗的圓角、方角及自訂圓角半徑支援，見 [PR](https://github.com/wailsapp/wails/pull/5866)，作者：@leaanthony

## 修正

- 缺少 Blob/FormData 標頭時，Fetch 相容層會傳遞 undefined，見 [PR](https://github.com/wailsapp/wails/pull/5865)，作者：@leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## 新增

- 新增 macOS DMG 封裝支援，包括新的選項和建置任務，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5857)

## 變更

- GitHub 更新程式現在預設排除 Windows 安裝程式資產，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5861)

## 修正

- 修正 Linux WebKit 在 fetch 請求中傳送 Blob 或 FormData 時當機的問題，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5854)

## v3.0.0-alpha2.120 - 2026-07-31

## 新增

- 實作 macOS 標題列按兩下時最大化或最小化的動作，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5853)

## 變更

- 更新 QR 服務教學課程以使用 NewServiceWithOptions，並加入間距，見 @jeongkyu 的[PR](https://github.com/wailsapp/wails/pull/5849)

## 修正

- 讓 WKWebView 在 macOS 縮放期間維持回應，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5856)
- 修正 GTK4 視窗大小查詢，並從設定的`GdkSurface`發出調整大小、最大化、最小化及全螢幕狀態事件。

## v3.0.0-alpha2.119 - 2026-07-27

## 修正

- 更新多語言文件以納入架構圖，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5833)

## v3.0.0-alpha2.118 - 2026-07-26

## 新增

- 為圖示產生作業的輸入與輸出提供預設路徑，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5825)
- 將原始碼進入點模組加入執行階段 package.json 的 sideEffects，見 @savely-krasovsky 的[PR](https://github.com/wailsapp/wails/pull/5797)
- 在貢獻指南中新增授權與來源章節，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5816)

## 修正

- 套用限定範圍的 GTK4 無框架 CSS，以移除圓角，見 @savely-krasovsky 的[PR](https://github.com/wailsapp/wails/pull/5800)
- 在 Windows 上，於快顯功能表與螢幕列舉期間妥善處理游標位置取得失敗，見 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5789)
- macOS 開啟檔案對話方塊現在能正確篩選副檔名，並依檔名後綴驗證允許的檔案，見 @phergul 的[PR](https://github.com/wailsapp/wails/pull/5678)
- 防止 Windows 深色模式初始化時呼叫 nil API，見 @roachadam 的[PR](https://github.com/wailsapp/wails/pull/5793)
- 修正更新程式的32位元建置失敗：在`GOARCH=386`上將`maxArchiveTotalSize`常數（2 GiB）傳給`fmt.Errorf`時，該常數超出平台`int`的範圍。現在已將其明確指定為`int64`型別。
- 修正啟動時發生的 nil 指標 panic：當視窗使用深色（或系統深色）標題列，而 Windows 建置未載入深色模式 uxtheme API 時（例如 Windows 10 1809／Windows Server 2019，組建17763），便會發生此問題。視窗佈景主題設定中的`AllowDarkModeForWindow`呼叫現在已加入 nil 防護，與`w32.SetMenuTheme`中既有的防護一致。

## v3.0.0-alpha2.117 - 2026-07-08

## 新增

- 為 Windows 上的非工作區域實作自訂命中測試邏輯，見 @savely-krasovsky 的[PR](https://github.com/wailsapp/wails/pull/5462)

## 變更

- 根據 UseVisualHosting 設定 WebView2 的顯示器縮放偵測，見 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5761)

## v3.0.0-alpha2.116 - 2026-07-07

## 新增

- 更新常見問題文件，使其聚焦於 Wails v3 的功能與指引，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5763)

## v3.0.0-alpha2.115 - 2026-07-06

## 修正

- 修正`Menu.Update()`未在 GTK4 Linux 上重新建置原生功能表的問題（#5659，@puneetdixit200 另行診斷並於 #5539中修正）
- 透過複製螢幕 ID／名稱字串並建立數量快照，修正顯示器變更時列舉 macOS 螢幕會當機的問題（#5565，@x-haose 另行診斷並於 #5584中修正）
- 修正 Windows 上的當機問題：在最小化／還原轉換期間，`GetClientRect`傳回 nil，而`WM_ERASEBKGND`繪製純色背景時會發生此問題（防護措施由 @sinspired 於 #5636中回報）
- 修正前端繫結錯誤一律解析為文字的問題，由 @mbaklor 於 #5690中修正
- 修正使用`server`建置標籤時 Windows 建置失敗的問題；原因是 Windows GUI 檔案缺少 macOS 與 Linux 對應檔案已有的`!server`建置條件（#5680）

## v3.0.0-alpha2.114 - 2026-07-05

## 新增

- 實作 Update Manifest 通訊協定與端點提供者，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5720)

## 變更

- 將`webview2`繫結併入 v3 模組並命名為`v3/internal/webview2`，同時移除獨立模組、其每夜發行／同步工作流程，以及 go.mod 版本調整程序（v3 是其唯一使用者），見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5711)

## 修正

- 將 WebView2 顯示器縮放偵測及 DPI 變更時的主機重新同步修正移至「未發行」章節，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5750)
- 更新 WebView2 對 float64 與 BOOL 參數的 COM 封送處理，見 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5741)
- 防止 Windows 系統匣圖示在更新及銷毀時發生 panic 與 nil 解參照，見 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5703)
- 修正 Windows 上已隱藏的視窗未能正確再次隱藏的問題，見 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5743)
- 在視窗最小化、最大化及還原時同步 WebView2 控制器的可見性，見 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5742)

### 修正

- 重新啟用 WebView2 顯示器縮放偵測，並僅在 DPI 變更時執行主機重新同步，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5734)；此變更以 @randalmurphal 驗證的修正為基礎，根本原因由 @eleclin 驗證，硬體測試則由 @qq540491950 進行

## v3.0.0-alpha2.113 - 2026-07-04

## 新增

- 建置發布版 AAB 時，若未設定`ANDROID_KEYSTORE_FILE`，則顯示警告（Google Play 會拒絕以偵錯金鑰簽署的套件），並在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5730)中記錄 App Bundle 封裝與簽署方式
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5739)中新增多語言的 Why Wails 文件
- 在 @fbbdev 的[PR](https://github.com/wailsapp/wails/pull/5398)中，支援在繫結中將 Go time.Time 對應至 JS Date 或字串
- 在 @mortenolsrud 的[PR](https://github.com/wailsapp/wails/pull/5728)中新增供提交至 Play Store 使用的 Android App Bundle（AAB）封裝工作（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）；APK 工作仍保留供本機／模擬器測試使用（修正[#5726](https://github.com/wailsapp/wails/issues/5726)）
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5735)中新增 Android 實體裝置工作目標，並恢復相機／位置權限

## 變更

- 將`webview2`升級至 v1.0.28（[版本資訊](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5730)中，將 Android 範本的`compileSdk`/`targetSdk`從34升級至35；Google Play 規定新提交的應用程式必須使用此版本

## 修正

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5745)中修正 sponsorkit 中已烘焙的頭像遮罩
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5730)中修正 Android AVD 自動建立功能因按字典順序排序版本，而選取錯誤系統映像或 cmdline-tools 版本的問題
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5730)中修正設定精靈建議過時 Android NDK 版本的問題（現為26.3.11579264，與文件所述需求一致）
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5744)中更新 SvelteKit 與選項的法文文件
- 在 @flofreud 的[PR](https://github.com/wailsapp/wails/pull/5516)中修正顯示器變更期間，macOS 列舉螢幕時發生 SIGSEGV 的問題

## v3.0.0-alpha2.112 - 2026-07-03

## 新增

- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5724)中新增以 Go 為基礎的貢獻者 SVG 產生器，並更新文件／網站的致謝頁面
- 在 @fbbdev 的[PR](https://github.com/wailsapp/wails/pull/5398)中，支援在繫結中將 Go time.Time 對應至 JS Date 或字串

## 變更

- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5719)中，以 Go 產生器取代以 Node 為基礎的贊助者圖片處理流程

## 修正

- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5729)中修正 Android 建置資產相依項目的安裝指令碼
- 在`ValidateAndSanitizeURL`中拒絕 U+0085（NEXT LINE）控制字元，補齊 URL 驗證器對空白字元的涵蓋範圍
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/4785)中，於無框線視窗的 DPI 變更時重新計算 DWM 窗框
- 在 @yulesxoxo 的[PR](https://github.com/wailsapp/wails/pull/4632)中修正 Windows 使用非100% 縮放比例時，DnD 放置區域偵測失敗的問題
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5714)中，為 Darwin 對話方塊、選單、系統匣和通知所使用的 Cocoa 物件新增明確的 Objective-C 記憶體管理
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5718)中修正 Linux CGO 後端錯誤與系統匣問題

## v3.0.0-alpha2.111 - 2026-07-01

## 新增

- 在 @Aliuyanfeng 的[PR](https://github.com/wailsapp/wails/pull/5061)中將 HappyTools 新增至社群展示區
- 在 @triadmoko 的[PR](https://github.com/wailsapp/wails/pull/5643)中新增印尼語地區設定支援與完整文件
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/4813)中為 WindowsWindow 新增 DisableMenu 選項

## 變更

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5617)中更新 Taskfile 範本與 CLI，使其使用 GOOS 和 ARCH 分派建置／封裝工作

## 修正

- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5708)中修正 Mac 視窗分頁問題

## 移除

- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5702)中，從貢獻、功能及指南中移除德文翻譯的 MDX 檔案

## v3.0.0-alpha2.110 - 2026-06-30

## 新增

- 在 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5129)中實作 macOS WebView 重新載入與強制重新載入，並新增 WebContent 行程終止後的復原機制
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5396)中新增涵蓋貢獻、功能及指南的完整德文文件
- 在 @popaprozac 的[PR](https://github.com/wailsapp/wails/pull/5333)中強化通知功能，新增音效、附件、排程與更新 API

## 修正

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/4785)中，於無框線視窗的 DPI 變更時重新計算 DWM 窗框
- 在 @yulesxoxo 的[PR](https://github.com/wailsapp/wails/pull/4632)中修正 Windows 使用非100% 縮放比例時，DnD 放置區域偵測失敗的問題

## v3.0.0-alpha2.109 - 2026-06-29

## 新增

- 在 @iamhabbeboy 的[PR](https://github.com/wailsapp/wails/pull/5026)中為 EventsEmit 文件新增程式碼範例
- 在 @MerIijn 的[PR](https://github.com/wailsapp/wails/pull/5380)中新增 Windows WebView2 視覺化託管選項
- 在 @SametKUM 的[PR](https://github.com/wailsapp/wails/pull/5536)中將 Klustr 新增至社群展示文件
- 在 @thiennguyen93 的[PR](https://github.com/wailsapp/wails/pull/5685)中，以新頁面及變更記錄項目將 Kira 新增至社群展示區
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5694)中為 MCP 服務指南新增意見回饋章節

## 變更

- 伺服器模式現在具備正式支援的一流生產環境建置流程，與桌面版建置工作一致（#5693）。`task build:server`預設會建置生產環境二進位檔（`-tags server,production`、`-trimpath`、移除符號），並接受`DEV=true`（開發伺服器）、`OBFUSCATED=true`（garble）及`EXTRA_TAGS`。`task run:server`會執行開發伺服器。`Dockerfile.server`／`task build:docker`會先建置生產環境伺服器（`-tags server,production`）與生產環境前端；映像檔預設是在 distroless/static 上使用純 Go 靜態建置，並公開`CGO_ENABLED`、`GO_IMAGE`及`RUNTIME_IMAGE`，作為可供 CGO 應用程式覆寫的建置引數。

## 已修正

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/4435)中，防止關閉仍有待處理非同步呼叫的視窗時發生當機
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5249)中，防止在 Windows 上開啟隱藏的應用程式時啟用視窗
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5668)中，確保 WebKit 要求中繼資料、回應完成作業及本文串流處理皆在 GTK 主執行緒上執行
- 修正`Menu.Update()`不會在 GTK4 Linux 上重建原生選單的問題（#5659；@puneetdixit200 也在 #5539中獨立診斷並修正此問題）
- 透過複製螢幕 ID／名稱字串並擷取數量快照，修正顯示器變更時列舉 macOS 螢幕所造成的當機（#5565；@x-haose 也在 #5584中獨立診斷並修正此問題）
- 透過在`WM_DPICHANGED`處理常式中重新套用控制器邊界，採用與取消最小化時 DPI 重新同步相同的方式，修正在 Windows 上將視窗拖曳穿越不同 DPI 的顯示器後，WebView2 內容先縮小再消失的問題（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 新增

- 透過`app.GlobalShortcut`新增全域（系統範圍）鍵盤快速鍵（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）。即使應用程式未取得焦點，快速鍵仍會觸發。各平台皆採用原生實作，且不含第三方相依套件：macOS 使用 Carbon 快速鍵、Windows 使用`RegisterHotKey`、X11 使用`XGrabKey`，Wayland 則使用 XDG Desktop Portal 全域快速鍵介面。
- 新增內建 MCP 伺服器：當應用程式使用`mcp`標記建置時，這個 Model Context Protocol 伺服器會自動啟動，讓 LLM 代理程式測試及控制執行中的 Wails 應用程式，包括視窗控制、DOM 檢查、JavaScript 求值、繫結方法呼叫、事件，以及以畫面上動畫游標呈現的模擬滑鼠／鍵盤輸入。不需要使用者程式碼：設定`WAILS_MCP=1`時，`wails3 build`／`wails3 dev`會自動加入`mcp`標記。完全透過環境變數設定（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）。

## 已修正

- 修正`Menu.Update()`不會在 GTK4 Linux 上重建原生選單的問題（#5659；@puneetdixit200 也在 #5539中獨立診斷並修正此問題）
- 透過複製螢幕 ID／名稱字串並擷取數量快照，修正顯示器變更時列舉 macOS 螢幕所造成的當機（#5565；@x-haose 也在 #5584中獨立診斷並修正此問題）

## v3.0.0-alpha2.107 - 2026-06-27

## 新增

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5613)中，新增包含側邊欄導覽的實驗性 Wake 文件

## v3.0.0-alpha2.106 - 2026-06-24

## 已變更

- 將`webview2`升級至 v1.0.27。
  - ci(webview2)：修正發行版本建置（交叉編譯 Windows + 完整的 go.sum）（#5671）\

  **完整差異：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5672)中，從 webview2 發行工作流程的交叉編譯作業移除 go vet
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5670)中，將 auto-changelog 的 OpenRouter 模型更新為 google/gemini-2.5-flash-lite
- 將`webview2`升級至 v1.0.26。

### 修正項目

- **從暫時性執行階段 COM 錯誤復原，而非結束程式**（#5658、#5580）。`Chromium.errorCallback`先前會對<em>任何</em> COM 錯誤呼叫`os.Exit(1)`，因此啟動後可復原的短暫異常會終止整個應用程式。執行階段路徑（`Resize`／`GetClientRect`、`Navigate`／`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）現在會記錄錯誤並復原。尤其是`MessageReceived`中的格式錯誤／不受信任 Web 訊息，現在會直接捨棄，而不會導致處理程序終止。這解決了跨越不同 DPI 顯示器時發生的一類當機（#5544、#5650）。環境／控制器建立路徑發生錯誤時，仍會終止程式。\

**完整差異：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 已修正

- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5671)中，修正 release-webview2 工作流程，使其能正確處理 go.sum 檔案
- 在 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5659)中，透過清除並重建原生選單來修正 Linux GTK4 選單更新

## v3.0.0-alpha2.105 - 2026-06-21

## 新增

- 新增`application.System`，以便從共用程式碼偵測執行階段平台：`System.IsMobile()`（iOS／Android）、`System.IsDesktop()`（macOS／Windows／Linux）、`System.IsServer()`（`server`建置標記），以及可直接測試單一目標的`System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`。它可在每個目標上編譯，因此你不需要建置標記即可進行條件分支。對應的前端輔助工具（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）可在`@wailsio/runtime`中使用
- 新增「使用其他前端框架」指南，說明如何將自己的 Vite 專案直接放入`frontend/`（涵蓋 Solid、Preact、Lit、SvelteKit、Qwik、Angular 等）
- `wails3 setup`精靈現在會檢查行動平台（iOS／Android）工具鏈，包括 Xcode 與 iOS Simulator 執行階段、JDK、Android SDK／NDK 和模擬器；在適用情況下，還提供一鍵安裝及可複製的 shell 設定修正方式
- 產生的專案隨附`frontend/.npmrc`，其中設定7天的`minimum-release-age`，以降低接觸剛發佈（且可能已遭入侵）套件的風險（pnpm 與 bun 會遵循此設定；npm 會安全地忽略）

## 已變更

- 以全新的霓虹山景主視覺重新設計所有內建入門範本（Web、iOS 及 Android）
- **TypeScript 現在是入門範本的預設語言，並使用不含語言限定詞的範本名稱。**`wails3 init`（不含`-t`）會建立 TypeScript 專案骨架；`-t vanilla`、`-t react`、`-t vue`及`-t svelte`均使用 TypeScript，而其 JavaScript 變體則為`-t vanilla-js`、`-t react-js`、`-t vue-js`及`-t svelte-js`。內建範本會在`template.yaml`中以`typescript:`宣告所用語言；使用`-ts`後綴的社群範本仍可作為後備方式運作
- 採用霓虹「digital Wails」主題（以山景為背景的霧面玻璃鮮明效果），重新設計`wails3 setup`精靈

## 修正

- 修正在 Windows 上，應用程式最小化的時間長到足以讓 WebView2 暫停或回收其算繪/GPU 處理程序後，還原應用程式時發生的當機。最小化/還原時的 DPI 重新同步（#5544）現在只會在視窗的 DPI 確實變更時存取 WebView2 控制器，避免在常見的相同 DPI 還原情況下，對已暫停的控制器進行致命的 COM 呼叫（#5605）
- 修正在頻繁載入資產/媒體的長時間執行 Linux 應用程式中，反覆發生的原生 `SIGABRT`/`SIGSEGV` 當機（通常發生在 GTK 主迴圈期間的 `g_object_unref` 內）。資產伺服器原本從工作 goroutine 完成 `WebKitURISchemeRequest`，因而在 GTK 主執行緒以外呼叫非執行緒安全的 WebKit2GTK 函式；現在完成作業（`webkit_uri_scheme_request_finish_with_response`/`finish_error`）會在主執行緒上執行。這項修正補完了 #5566 中的部分修正。GTK3 與 GTK4/WebKitGTK 6.0 組建均受影響（#5631、#5557）
- 修正 Linux/GTK3 上 `setupSignalHandlers` 中間歇性出現的 `fatal error: invalid pointer found on stack`。以訊號 `user_data` 傳遞的視窗 ID 原本儲存在 Go 的 `unsafe.Pointer` 區域變數中，因此垃圾回收器在複製堆疊期間掃描到這個（非指標）值時會中止。現在 Go 端會將 ID 保持為整數型別（`uintptr_t`），把 #4958 套用於 GTK4 路徑的相同修正回移植至舊版 GTK3 路徑（該修正將 C 訊號函式改用 `uintptr_t`，以消除 `-race`/checkptr 錯誤）（#5631）

## 移除

- 移除 `react-swc`、`preact`、`lit`、`solid`、`qwik` 和 `sveltekit` 起始範本（以及其 `-ts` 變體）。目前支援的內建範本為 `vanilla`、`react`、`vue` 和 `svelte`；每個範本預設都使用 TypeScript，另有 `-js` JavaScript 變體。仍可透過[使用您自己的前端](https://v3.wails.io/guides/dev/frontend-frameworks)或自訂範本來使用任何其他框架

## v3.0.0-alpha2.104 - 2026-06-18

## 修正

- 修正繫結的 Go 服務方法傳回空字串時發生的 iOS 當機（SIGABRT）。iOS 資產回應寫入器原本使用 `buf != nil` 而非主體長度來防護主體指標，因此長度為零的主體會導致 `&buf[0]` 發生 panic；現在改為依長度防護，與桌面版寫入器一致

## v3.0.0-alpha2.103 - 2026-06-15

## 變更

- 將 iOS 與 Android 原生功能移至平台管理器：請透過 `application.IOS.*` 和 `application.Android.*` 呼叫（例如 `application.IOS.Haptic("medium")`、`application.Android.Share(payload)`），不要再使用舊的 `application.IOS*`/`application.Android*` 自由函式（#5602）
- 重新命名行動版橋接事件：跨平台事件現在使用 `common:*` 前綴（例如 `common:haptic`、`common:location`），平台專屬事件則使用 `ios:*` / `android:*`（例如 `ios:backgroundTask`、`android:foregroundService`）；不再使用 `native:*` 前綴（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 新增

- 新增實驗性的 `wails3 setup` 精靈，用於互動式專案設定與相依性檢查
- 為 `wails3 doctor` 新增 `--json` 旗標，以輸出機器可讀格式
- 在 `wails3 doctor` 命令中新增簽署狀態區段

## 修正

- 修正 Linux 上的 npm 偵測，除了套件管理員之外也會檢查 PATH

## v3.0.0-alpha.101 - 2026-06-13

## 新增

- iOS：原生訊息對話方塊（UIAlertController），以及開啟檔案/多個檔案/目錄的對話方塊（UIDocumentPickerViewController）；儲存對話方塊會傳回明確的錯誤
- iOS：透過 UIPasteboard 支援剪貼簿
- iOS：透過 UIScreen 取得實際螢幕度量（點、像素、縮放比例、安全區域工作區）
- iOS：裝置組建（`IOS_PLATFORM=device`）、程式碼簽署身分/佈建描述檔/權利支援、`.ipa` 封裝，以及透過 devicectl 執行 `deploy-device`
- iOS：可設定最低 iOS 版本（build/config.yml 中的 `ios.minIOSVersion`）
- iOS：`wails3 doctor` 會回報 macOS 上 Xcode 與 iOS SDK 的可用性
- iOS：系統事件——電池、網路、佈景主題、螢幕鎖定與記憶體不足會以 `events.IOS.*` 和平台中立的 `events.Common.*` 應用程式事件公開
- iOS：原生行動功能橋接（匯出的 `application.IOS*`）——分享面板、開啟 URL、防止休眠、手電筒、安全區域內距、亮度、應用程式資訊、方向鎖定、狀態列、生物辨識（Face ID/Touch ID）、本機通知，以及 Keychain 安全儲存空間
- iOS：感測器與硬體——觸覺回饋、單次地理位置、加速度計、鄰近感測器、文字轉語音、儲存空間資訊、電源/電池狀態、網路狀態、鍵盤內距，以及螢幕擷取偵測
- iOS：文件（IOS.md 與文件網站指南）
- Android：原生訊息對話方塊（AlertDialog），以及開啟檔案/多個檔案的對話方塊（Storage Access Framework，匯入為快取副本）；開啟目錄與儲存對話方塊會傳回明確的錯誤
- Android：透過 ClipboardManager 支援剪貼簿
- Android：透過 WindowMetrics/DisplayMetrics 取得實際螢幕度量（dp、像素、縮放比例、系統列工作區）
- Android：觸覺回饋（`Android.Haptics.Vibrate`）、裝置資訊（`Android.Device.Info`）與快顯通知（`Android.Toast.Show`）執行階段方法
- Android：具型別的生命週期事件（`events.Android.*`，由 events.txt 產生），並將 `ActivityCreated` 對應至 `Common.ApplicationStarted`
- Android：組建管線會產生可安裝的偵錯版與發行版 APK（`android:run`、`android:package`、`android:package:fat`）；發行版預設使用偵錯金鑰庫簽署，也可透過 `ANDROID_KEYSTORE_*` 環境變數使用實際的金鑰庫
- Android：`wails3 doctor` 會回報 Android SDK、NDK 與 JDK
- Android：系統事件——電池、網路、佈景主題、螢幕鎖定與記憶體不足會以 `events.Android.*` 和平台中立的 `events.Common.*` 應用程式事件公開
- Android：原生行動功能橋接（匯出的 `application.Android*`）——分享、開啟 URL、防止休眠、手電筒、安全區域內距、亮度、應用程式資訊、方向鎖定、狀態列、生物辨識（BiometricPrompt）、本機通知，以及 EncryptedSharedPreferences 安全儲存空間
- Android：感測器與硬體——觸覺回饋、單次地理位置、加速度計、鄰近感測器、文字轉語音、儲存空間資訊、電源/電池狀態、網路狀態、鍵盤內距，以及透過 FLAG_SECURE 封鎖螢幕擷取
- Android：文件（ANDROID.md 與文件網站指南）
- 範例：`mobile` 綜合展示新增「行動裝置」與「硬體」分頁，示範跨 iOS 與 Android 的原生功能橋接（膠囊式分頁會換行成多列）
- 行動版：電池——應用程式進入背景時會暫停加速度計、鄰近感測器、手電筒和範例中的週期性時鐘，返回時再恢復（Android 會讓處理程序繼續在背景執行，而手電筒在 iOS 上是會持續維持的硬體狀態）；此外，Android 系統事件接收器只會在應用程式位於前景時註冊
- iOS：相機拍攝 — `application.IOSCapturePhoto`/`IOSCaptureVideo`（UIImagePickerController → 含有 base64 縮圖的 `native:capture` 事件）
- iOS：背景執行 — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask`（UIApplication 背景工作時段），以及可設定的 `ios.backgroundModes`（build/config.yml），用來將 `UIBackgroundModes` 套用至產生的 Info.plist
- Android：相機拍攝 — `application.AndroidCapturePhoto`/`AndroidCaptureVideo`（透過 FileProvider 使用系統相機 → `native:capture` 事件）
- Android：前景服務 — `application.AndroidStartForegroundService`/`AndroidStopForegroundService`（具有持續顯示通知的 `WailsForegroundService`，可讓處理程序持續運作，以執行長時間的背景工作）
- 範例：新增相機分頁，示範相片／影片拍攝與背景執行（Android 使用前景服務，iOS 使用背景工作時段）

## 已修正

- 修正 Linux 上 `getUserMedia` 總是因 `NotAllowedError` 而失敗的問題：WebKitGTK 會拒絕未由任何處理常式處理的權限要求，而 `permission-request` 訊號並未連接。現在會依據新的跨平台 `WebviewWindowOptions.Permissions` 對應表（`map[PermissionType]Permission`）處理相機／麥克風權限，且 Linux（WebKitGTK）與 Windows（WebView2）皆會遵循此設定。在沒有原生提示視窗的 Linux 上，相機／麥克風現在預設為允許（恢復 `getUserMedia`），並可透過 `PermissionDeny` 關閉（#5552）
- iOS：`GOOS=ios` 現在可再次編譯（匯出 `events.IOS`、加入行動版方法名稱存根），使用 production 標籤的組建也可編譯（修正 pkg/application 與數個服務中的建置標籤）
- iOS：Go→JS 事件與 ExecJS 現在可正常運作 — 頁面啟動時不再載入兩次，且 `wails:runtime:ready` 交握不會再遺失
- iOS：`ApplicationDidFinishLaunching`/`ApplicationStarted` 不再與應用程式啟動產生競爭；已移除固定等待 2 秒的啟動延遲
- iOS：修正每次執行 Go→JS JavaScript 時發生的 C 字串記憶體洩漏
- iOS：`hasListeners` 現在會反映實際的接聽器註冊狀態
- iOS：正式版組建在編譯時會排除框架偵錯記錄
- Android：`GOOS=android` 現在可再次編譯 — 已定義 `events.Android`、移除越界的 `events_android.go` 接聽器陣列、加入行動版方法名稱存根，並防止桌面版 Linux 檔案（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）混入 Android 組建
- Android：JS→Go 繫結現在可正常運作 — WebView 無法將 `fetch()` POST 主體傳送至 `shouldInterceptRequest`，因此執行階段呼叫現在會透過 JavascriptInterface 傳輸層（`nativeHandleRuntimeCall`）路由，而不會因要求主體為 nil 而當機
- Android：`Screens.*` 執行階段呼叫現在會傳回實際資料 — ScreenManager 現在會在啟動時填入資料（先前從未接線，因此 `GetAll` 會傳回 nil）
- Android：正式版組建在編譯時會排除框架偵錯記錄；偵錯組建則會以 `Wails` 標籤將記錄傳送至 logcat
- Android：實作真正的 `hasListeners` 登錄機制、JNI 參照／例外處理，以及僅載入一次的頁面生命週期（不再重複導覽）
- 修正 Vite 開發伺服器執行時，`wails3 generate bindings` 在 Windows 上因「Access is denied」而失敗的問題：現在會將產生的檔案同步至輸出目錄，而不是透過重新命名覆蓋該目錄（#5515）
- 修正 macOS 在顯示器變更後讀取螢幕資訊時偶發的嚴重當機：儲存的螢幕 ID 與名稱指標指向自動釋放的 `UTF8String` 緩衝區，而這些緩衝區可能在 Go 複製前就已釋放（釋放後使用）。現在會對字串執行 `strdup`，並在轉換後釋放；螢幕列舉也會在明確的自動釋放集區中執行，因此從 Go goroutine 呼叫時不再發生記憶體洩漏（#5556）
- 修正 Linux 上 assetserver 關閉`WebKitURISchemeRequest`時偶發的 SIGSEGV：最後一次`g_object_unref`會在 assetserver goroutine 上執行，導致 WebKit GObject 在 GTK 主執行緒之外完成終結。現在會透過`g_main_context_invoke`將 unref 派送至 GTK 主執行上下文執行（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 已新增

- 擴充 `MacWebviewPreferences`，加入其他 WKWebView 設定選項：`EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize` 與 `ApplicationNameForUserAgent`（#5549）

## 已修正

- 修正 Vite 開發伺服器執行時，`wails3 generate bindings` 在 Windows 上因「Access is denied」而失敗的問題：現在會將產生的檔案同步至輸出目錄，而不是透過重新命名覆蓋該目錄（#5561）
- 修正 Linux 上無框視窗不會觸發 JS 調整大小事件的問題；修正無框視窗的捲動條邊緣偵測（#5368）
- 修正 Windows 上暫存目錄與安裝目錄位於不同磁碟區時，更新程式因「invalid cross-device link」而失敗的問題（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 已修正

- 修正 Vite 開發伺服器執行時，`wails3 generate bindings` 在 Windows 上因「Access is denied」而失敗的問題：現在會將產生的檔案同步至輸出目錄，而不是透過重新命名覆蓋該目錄（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 已修正

- 修正 Linux 上 WebKit 閒置時（例如開啟檢查器時）UI 凍結的問題：不再對 `SIGUSR1` 強制設定 `SA_ONSTACK`，因為這會破壞 JavaScriptCore 的 GC 執行緒同步（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 已新增

- 新增偵錯頁面以及使用 `runtime/trace` 的相關內容

## 已變更

- 移除部分不必要的 `_ "embed"` 匯入，稍微整理程式碼

## 已修正

- 修正 Windows 上視窗取消最大化後，未強制套用最小寬度／高度限制的問題（#4593）
- 修正同時使用 Frameless 與 Transparent 視窗選項時，全螢幕模式下滑鼠點擊會穿透的問題（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 已新增

- 新增 Garble 混淆支援（[#4563](https://github.com/wailsapp/wails/issues/4563)）：穩定的繫結方法 ID、建置／Taskfile 串接（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`），以及每個面向執行階段的承載資料所需的 JSON 結構標籤（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`），使線路格式在 Garble 重新命名匯出欄位後仍能維持不變。

## v3.0.0-alpha.95 - 2026-05-20

## 已新增

- 新增遺漏的專案結構頁面

## 已變更

- 文件：變更架構頁面上的數個圖表，改用循序圖以獲得更簡潔的顯示效果
- 文件：加入執行前須先安裝 D2 的說明

## 修正

- 修正 GTK4 預設設定中的`wails3 generate appimage`：打包工具現在會先從二進位檔偵測 GTK 技術堆疊，再搜尋執行階段檔案，因此會為 GTK4 組建選取`libwebkitgtkinjectedbundle.so`（位於`webkitgtk-6.0/`下），並為`-tags gtk3`組建選取`libwebkit2gtkinjectedbundle.so`（位於`webkit2gtk-4.1/`下）。`.relr.dyn`探測現在也會檢查`libgtk-4.so.1`，因此無論使用哪種技術堆疊，在現代工具鏈上都能正確停用符號剝除。（#5475）
- 修正以相對`-builddir`呼叫`wails3 generate appimage`時失敗的問題：打包工具現在會預先將`-binary`、`-icon`、`-desktopfile`、`-builddir`和`-outputdir`解析為絕對路徑，避免流程中途的`s.CD`破壞 AppRun 下載 goroutine 或複製後的`ldd`探測。
- 修正桌面檔案的`Name=`欄位與二進位檔基本名稱不符時，`wails3 generate appimage`無法將最終 AppImage 移至`-outputdir`的問題：打包工具現在會透過`OUTPUT`環境變數，強制 linuxdeploy 的 appimage 外掛程式將 AppImage 寫入`<binary>-<arch>.AppImage`，而非使用根據桌面檔案衍生的名稱。
- 修正 GTK4 + WebKitGTK 6.0技術堆疊在 alpha.93升為預設值後，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`和`Common.SystemDidWake`未在 Linux 上觸發的問題。新的預設`application_linux.go` `run()`未呼叫`setupCommonEvents()`（此函式會將`Linux.*`事件轉送至對應的`Common.*`事件）或`monitorPowerEvents()`。現在 GTK3 與 GTK4 組建路徑會透過`application_linux_dbus.go`共用 DBus 電源監控輔助程式。（#5474）

## v3.0.0-alpha.94 - 2026-05-19

## 修正

- 修正 GTK4 + WebKitGTK 6.0技術堆疊在 alpha.93升為預設值後，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`和`Common.SystemDidWake`未在 Linux 上觸發的問題。新的預設`application_linux.go` `run()`未呼叫`setupCommonEvents()`（此函式會將`Linux.*`事件轉送至對應的`Common.*`事件）或`monitorPowerEvents()`。現在 GTK3 與 GTK4 組建路徑會透過`application_linux_dbus.go`共用 DBus 電源監控輔助程式。（#5474）

## v3.0.0-alpha.93 - 2026-05-17

## 新增

- 由 @leaanthony 在 Linux 的`wails3 doctor`輸出中加入`XDG_SESSION_TYPE`

## 修正

- 由 @leaanthony 修正 appmenu-gtk-module 存取尚未具現化的視窗而導致 Wayland 上的視窗選單當機問題（#4769）
- 由 @leaanthony 修正應用程式名稱含有無效字元（空格、括號等）時 GTK 應用程式當機的問題
- 由 @overlordtm 修正 Windows 上初始化拖放功能時出現「記憶體不足」錯誤的問題（#4701）
- 由 @leaanthony 修正主執行緒回呼儲存區在刪除 map 項目時誤用 RLock 所造成的競爭條件（Linux、macOS、iOS）（#4424）
- 修正將命令列引數傳遞給工作時的變數處理。現在能正確初始化以 KEY=VALUE 配對指定的 CLI 變數，並在整個工作執行期間傳遞這些變數。
- 修正 macOS 上的 NSWindowZoomButton 衝突：`MaximiseButtonState`和`FullscreenButtonState`現在會在啟動及執行階段套用限制較嚴格的狀態，任一 setter 都無法再悄悄覆寫另一個（#5319）
- 修正 CodeRabbit 在 #5463中揭露的舊版 GTK3 組建路徑（`-tags gtk3`）既有問題群組：透過檔案關聯啟動時不再略過啟動處理常式；`getTheme`現在具備邊界與型別安全性；`appName`不再釋放 GLib 擁有的記憶體；`clipboardGet`不再洩漏 GTK 傳回的`gchar*`；`Calloc`現在使用指標接收器（且`NewCalloc`會傳回`*Calloc`），讓集區能實際追蹤配置；`zoomOut`改用`zoomInFactor`的倒數，不再使用會被限制為1.0的負乘數；`execJS`會重複使用預先配置的空 world-name，不再每次呼叫都洩漏一個`C.CString("")`；並已從`menuItem.setAccelerator`移除開發用的`fmt.Println`。解決 #5465。
- 修正預設 GTK4 組建路徑（`linux_cgo.go`）中相同的`Calloc`值接收器洩漏：改用指標接收器及`NewCalloc() *Calloc`，讓每個視窗的`c.String(...)`配置能被實際追蹤並釋放。

## v3.0.0-alpha.92 - 2026-05-15

## 新增

- 修改 Taskfile，以允許透過`PACKAGE_MANAGER`選項控制所使用的前端套件管理員
- 在範本資料中加入`{{.Opn}}`和`{{.Cls}}`，使 Taskfile 範本的撰寫更可預測

## 變更

- 修改數個現有的 Taskfile 以使用`{{.Opn}} and {{.Cls}}`

## 修正

- 修正面板讀取系統匣選單時更新該選單，會在`linuxSystemTray`中引發`concurrent map read and map write`執行階段致命錯誤的問題。
- WebView2 的錯誤與堆疊追蹤輸出改用`log`而非`fmt`，避免應用程式在 Windows 上未連接主控台執行時遺失訊息。

## v3.0.0-alpha.91 - 2026-05-12

## 變更

- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5414)中更新贊助者 SVG
- <strong>重大變更（macOS）：</strong>統一 macOS 座標系統，使`GetScreens`、`Position`和`SetPosition`全都使用相同的空間：以邏輯點為單位、Y 軸向下，且`(0,0)`位於主要螢幕左上角。這與 Windows、GTK，以及 Electron 和 Web 的公開 API 一致。實體位置高於主要螢幕的螢幕現在會回報負`Bounds.Y`（先前為正值），且`Position()`／`SetPosition()`值現在使用邏輯點，而非`points × primaryScale`。仍可保證`Position()` → `SetPosition()`往返轉換不變；從先前 alpha 組建記錄的絕對值，或手動計算的因應措施（例如乘以`primaryScale`，或依螢幕高度反轉 Y 軸）都需要更新。解決[#5117](https://github.com/wailsapp/wails/issues/5117)。

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5416)中以防禦性方式驗證 DBus 訊號名稱與主體長度，以防止 panic
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5363)中修正 Linux 上 GTK 選單處理的記憶體安全性問題
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5295)中加入 NVIDIA GPU 偵測，並在 Linux 上停用 DMA-BUF 渲染器
- 修正 macOS 上`SetPosition`的跨螢幕 Y 座標轉換：使用主要螢幕高度作為全域基準，讓視窗在與主要顯示器垂直錯位的顯示器上落於正確位置，詳見[#5117](https://github.com/wailsapp/wails/issues/5117)
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5109)中修正 git PR 範本，使其指向正確的意見回饋 URL
- 修正在 Windows 系統匣中因損壞的`DestroyMenu`系統呼叫所造成的一系列`SetMenu`當機問題；該呼叫誤傳了四個引數而非一個，導致每次呼叫都傳回 FALSE，且未釋放任何資源。現在也會在重建選單時釋放 HMENU 與 HBITMAP 控制代碼（包括執行階段透過`MenuItem.SetBitmap`配置的控制代碼）、重設`Win32Menu.Update`中過時的核取方塊／選項按鈕對應，並移除`systemtray.updateMenu`中會使配置量加倍的多餘`Update()`呼叫。長時間執行的系統匣應用程式不再於每次重建選單時洩漏 GDI/USER 物件。

## v3.0.0-alpha.90 - 2026-05-11

## 新增

- 新增可設定 macOS 上 WKWebView User-Agent 應用程式名稱的功能，見[PR](https://github.com/wailsapp/wails/pull/5261)，由 @vinhvoit225 貢獻
- 在 gin-service 範例中新增間接相依套件 github.com/coder/websocket，見[PR](https://github.com/wailsapp/wails/pull/5400)，由 @taliesin-ai 貢獻
- 為建置資產測試新增深度相等比較支援，見[PR](https://github.com/wailsapp/wails/pull/5402)，由 @leaanthony 貢獻

## 變更

- 將建置輸出整合至 assets 目錄，見[PR](https://github.com/wailsapp/wails/pull/5401)，由 @taliesin-ai 貢獻
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5399)中更新贊助者 SVG

## 修正

- 在 macOS 單一執行個體訊息中使用通知物件，見[PR](https://github.com/wailsapp/wails/pull/5289)，由 @overlordtm 貢獻
- 批次處理 Windows 回呼，以避免高負載時遺失 Promise，見[PR](https://github.com/wailsapp/wails/pull/5383)，由 @taliesin-ai 貢獻

## v3.0.0-alpha.89 - 2026-05-10

## 新增

- 新增 go<em>test</em>results 作業以彙整 Go 測試結果，見[PR](https://github.com/wailsapp/wails/pull/5316)，由 @leaanthony 貢獻

## 變更

- 依條件將大型 RPC 承載資料拆分為分塊 POST 要求，見[PR](https://github.com/wailsapp/wails/pull/5369)，由 @leaanthony 貢獻
- 將所有前端範本中的 Vite 從5.x.x 升級至8.0.0，見[PR](https://github.com/wailsapp/wails/pull/5386)，由 @leaanthony 貢獻
- 將 Vite 開發伺服器的連接埠設定移轉至環境變數，見[PR](https://github.com/wailsapp/wails/pull/5365)，由 @leaanthony 貢獻
- 將所有範本中的 Vite 開發伺服器設定為繫結至127.0.0.1，見[PR](https://github.com/wailsapp/wails/pull/5361)，由 @leaanthony 貢獻
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5384)中更新贊助者 SVG

## 修正

- 在更新建置資產時清理 Info.plist 範本存根，見[PR](https://github.com/wailsapp/wails/pull/5312)，由 @leaanthony 貢獻
- 在主執行緒上同步套用選單項目修改方法（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`），修正 macOS 選單中的過時狀態；這消除了快速重新開啟選單時仍會呈現先前狀態的`dispatch_async`競爭條件（#5002）
- 在開發模式中忽略`*_test.go`檔案，以避免不必要的重新建置，見[PR](https://github.com/wailsapp/wails/pull/5203)，由 @leaanthony 貢獻
- 避免應用程式未執行時 Menu.Update() 發生記憶體區段錯誤，見[PR](https://github.com/wailsapp/wails/pull/5291)，由 @wucm667 貢獻
- 在 Windows 上使用 lastSizeWParam 控制是否重繪選單列，見[PR](https://github.com/wailsapp/wails/pull/5382)，由 @taliesin-ai 貢獻

## v3.0.0-alpha.88 - 2026-05-09

## 變更

- 將 HiddenOnTaskbar 改為使用 WS<em>EX</em>TOOLWINDOW，見[PR](https://github.com/wailsapp/wails/pull/5371)，由 @leaanthony 貢獻
- 重新排列相依套件，並移除 go.mod 中的 webview2 replace 指令，見[PR](https://github.com/wailsapp/wails/pull/5370)，由 @atterpac 貢獻
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5358)中更新贊助者 SVG

## 修正

- 移除泛型間接別名並整合對應鍵型別，見[PR](https://github.com/wailsapp/wails/pull/5331)，由 @fbbdev 貢獻

## 移除

- 刪除 PR-master 工作流程，並移除文件、Go 測試及略過測試，見[PR](https://github.com/wailsapp/wails/pull/5377)，由 @leaanthony 貢獻

## v3.0.0-alpha.87 - 2026-05-07

## 新增

- 新增 Wails v3 韓文文件，見[PR](https://github.com/wailsapp/wails/pull/5352)，由 @leaanthony 貢獻
- 新增安裝與快速入門的法文文件，見[PR](https://github.com/wailsapp/wails/pull/5354)，由 @leaanthony 貢獻
- 新增快速入門、概念與社群的葡萄牙文文件，見[PR](https://github.com/wailsapp/wails/pull/5355)，由 @leaanthony 貢獻

## v3.0.0-alpha.86 - 2026-05-06

## 新增

- 新增法文文件本地化，見[PR](https://github.com/wailsapp/wails/pull/5328)，由 @leaanthony 貢獻
- 在文件網站中新增德文語系，見[PR](https://github.com/wailsapp/wails/pull/5343)，由 @leaanthony 貢獻

## 變更

- 在文件設定中註冊全部8個已翻譯語系，見[PR](https://github.com/wailsapp/wails/pull/5347)，由 @leaanthony 貢獻
- 更新多個與 Windows 上 WebView2 相關的檔案，見[PR](https://github.com/wailsapp/wails/pull/5317)，由 @leaanthony 貢獻

## 修正

- 在 Linux 上將對話方塊分派拆分為 GTK3 與 GTK4，見[PR](https://github.com/wailsapp/wails/pull/5340)，由 @leaanthony 貢獻
- 確保對話方塊回呼在 GTK 執行緒上執行，以修正記憶體區段錯誤，見[PR](https://github.com/wailsapp/wails/pull/5339)，由 @leaanthony 貢獻

## v3.0.0-alpha.85 - 2026-05-05

## 新增

- 將 PR 範本 URL 新增至儲存庫，見[PR](https://github.com/wailsapp/wails/pull/5179)，由 @leaanthony 貢獻
- 新增 Wails v3 德文文件，見[PR](https://github.com/wailsapp/wails/pull/5330)，由 @leaanthony 貢獻

## v3.0.0-alpha.84 - 2026-05-03

## 新增

- 新增可停用 macOS 上按 Escape 鍵退出全螢幕的選項，見[PR](https://github.com/wailsapp/wails/pull/5307)，由 @leaanthony 貢獻
- 新增可停用 macOS 上按 Escape 鍵退出全螢幕的選項，見[PR](https://github.com/wailsapp/wails/pull/5310)，由 @leaanthony 貢獻
- 由 @yuseferi 在[PR](https://github.com/wailsapp/wails/pull/5288)中新增 Pausa 社群展示文件

## 變更

- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5308)中更新贊助者 SVG
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5309)中更新圖示產生指令，以處理不支援的平台
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5224)中以三態 ButtonState 取代布林值全螢幕 API，並實作平台繫結

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5315)中防止 WebView2 在控制器狀態為 nil 時執行焦點操作
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5313)中更新 GitHub Actions 工作流程，使其正確參照 PR 的基底分支
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5203)中於開發模式下忽略`*_test.go`檔案，以避免不必要的重新建置
- 由 @wucm667 在[PR](https://github.com/wailsapp/wails/pull/5291)中防止應用程式未執行時 Menu.Update() 發生記憶體區段錯誤

## v3.0.0-alpha.83 - 2026-05-02

## 新增

- 由 @symball 在[PR](https://github.com/wailsapp/wails/pull/5094)中新增 InstallScope 旗標及電腦／使用者安裝的建置選項
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5294)中為 BrowserWindow 新增無操作的 SetScreen 方法，以滿足 Window 介面

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5295)中偵測 NVIDIA GPU，並在 Linux 上停用 DMA-BUF 轉譯器
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5109)中修正 Git PR 範本，使其指向正確的意見回饋 URL
- 修正一系列 Windows 系統匣`SetMenu`當機問題。其肇因於損壞的`DestroyMenu`系統呼叫傳入四個引數而非一個，導致每次呼叫皆傳回 FALSE，且未釋放任何資源。另在重建選單時釋放 HMENU 與 HBITMAP 控制代碼（包括執行階段透過`MenuItem.SetBitmap`配置的控制代碼）、重設`Win32Menu.Update`中過時的核取方塊／選項按鈕對應，並移除`systemtray.updateMenu`中造成配置量加倍的多餘`Update()`呼叫。長時間執行的系統匣應用程式不再於每次重建選單時洩漏 GDI/USER 物件。

## v3.0.0-alpha.82 - 2026-05-01

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5232)中修正桌面檔案產生功能，使其正確處理桌面名稱

## v3.0.0-alpha.81 - 2026-04-30

## 變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5286)中將每夜版本發佈排程調整為 UTC 15:00

## 修正

- 修正 Retina Mac 上 Screen Bounds、WorkArea 及 Size 數值減半的問題 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5285)中更新文件相依套件及內容集合載入器

## v3.0.0-alpha.79 - 2026-04-29

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5270)中授予 trigger-release 作業 actions: write 權限

## 變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5283)中將發佈任務的預設分支設為 master，並更新變更日誌的措辭
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5282)中更新自動產生變更日誌的工作流程，以使用最新版本
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5280)中新增路徑篩選條件並移除已停用的工作流程，以提升工作流程效率
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5274)中更新文件，使範例連結參照 master 分支
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5272)中更新 v3 的文件與範例

## 修正

- 由 @AkagiYui 在[PR](https://github.com/wailsapp/wails/pull/5265)中為反向 Proxy 新增重試邏輯，並在開發環境強制使用 IPv4
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5281)中重寫未發佈變更日誌的觸發工作流程

## 移除

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5267)中移除用於各種測試用途的 Shell 測試指令碼
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5266)中刪除 v3-alpha 文件部署工作流程及 CNAME 記錄

### 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5196)中將「前端路由」項目新增至側邊欄導覽
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5185)中新增前端路由指南，提供各框架專屬的建議
- 新增對模態面板（macOS）的支援
- 由 @leaanthony 升級 ghw 版本，以更妥善支援 Apple 裝置 (#4977)
- 為 Dock 服務新增`GetBadge`方法
- 為`wails3 build`指令新增`-tags`旗標，以傳遞自訂 Go 建置標籤（例如`wails3 build -tags gtk4`）(#4957)
- 新增繫結產生器自動產生列舉的文件，包括專屬的「列舉」頁面及側邊欄導覽（#4972）
- 為`wails3 build`指令新增`-tags`旗標，以傳遞自訂 Go 建置標籤（例如`wails3 build -tags gtk4`）(#4957)
- 在`v3/examples/web-apis/`中新增 Web API 範例，展示41項瀏覽器 API，包括儲存空間（localStorage、sessionStorage、IndexedDB、Cache API）、網路（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、媒體（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、裝置（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、效能（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）等
- 新增 WebView API 相容性檢查器範例（`v3/examples/webview-api-check/`），用於跨平台測試200多項瀏覽器 API
- 新增`internal/libpath`套件，用於在 Linux 上尋找原生程式庫路徑，支援平行搜尋、快取及 Flatpak/Snap/Nix
- <strong>開發中：</strong>新增實驗性的 Linux WebKitGTK 6.0／GTK4 支援，可透過`-tags gtk4`使用（GTK3/WebKit2GTK 4.1仍為預設值）
- 注意：在並排式視窗管理員（例如 Hyprland、Sway）上，由於視窗幾何配置由視窗管理員控制，最小化／最大化操作可能不會如預期運作
- 由 @AbdelhadiSeddar 在<strong>在 JavaScript 中監聽事件</strong>的文件中新增如何使用<strong>一次性處理常式</strong>的說明
- 由 @leaanthony 為`WebviewWindowOptions`新增`UseApplicationMenu`選項，讓 Windows/Linux 上的視窗可繼承透過`app.Menu.Set()`設定的應用程式選單
- 由 @wimaha 新增使用`.icon`檔案（Apple Icon Composer 格式）產生 Liquid Glass 圖示和資產目錄（macOS）的支援（#4934）
- 新增用於無頭／Web 部署的實驗性伺服器模式（`-tags server`）。可在不依賴原生 GUI 的情況下，將 Wails 應用程式作為 HTTP 伺服器執行。請使用`wails3 task build:server`建置。詳情請參閱`examples/server`。
- 新增`internal/libpath`套件，用於在 Linux 上尋找原生程式庫路徑，具備平行搜尋與快取功能，並支援 Flatpak/Snap/Nix
- 由 @leaanthony 為`MacWindow`新增`CollectionBehavior`選項，用於控制視窗跨 macOS Spaces 及全螢幕時的行為（#4756）
- 由 @leaanthony 為 pkg/application 新增單元測試
- 由 @leaanthony 為 MSIX 封裝新增自訂通訊協定支援
- 新增 Linux 桌面環境偵測[PR #4797](https://github.com/wailsapp/wails/pull/4797)
- 由 @leaanthony 為 JavaScript 執行階段新增`Window.Print()`方法，以便從前端觸發列印對話方塊（#4290）
- 由 @leaanthony 在 Linux 的`wails3 doctor`輸出中新增`XDG_SESSION_TYPE`
- 由 @leaanthony 為 Linux 新增額外的 WebKit2 載入狀態變更事件：`WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished`（#3896）
- 由 @leaanthony 在 Linux 的`wails3 doctor`輸出中新增`XDG_SESSION_TYPE`
- 在 Linux 建置期間產生`.desktop`檔案，而不再僅於封裝時產生（#4575）
- 由 @leaanthony 新增 Linux 執行階段相依套件文件，包含各發行版專用的套件名稱與 nfpm 封裝範例（#4339）
- 由 @leaanthony 在 Linux 的`wails3 doctor`輸出中新增 NVIDIA 驅動程式版本資訊
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中為原始訊息處理常式新增來源資訊
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中新增 macOS 通用連結支援
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重構繫結傳輸層
- 由 @chinenual 在[PR](https://github.com/wailsapp/wails/pull/4760)中為 helloworld 範本新增 aria-label 識別碼，讓 Appium 測試用戶端可輕鬆測試範例應用程式
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中為原始訊息處理常式新增來源資訊
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中新增 macOS 通用連結支援
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重構繫結傳輸層
- 由 @fbbdev 與 @ianvs 在[#4633](https://github.com/wailsapp/wails/pull/4633)中新增具型別事件
- 新增`systray-clock`範例，示範具有即時工具提示更新功能的無頭系統匣（#4653）。
- 由 @Tolfx 在 #4510中新增適用於 Windows 的 NSIS Protocol 範本
- 由 @Tolfx 在 #4510中新增 build-assets 測試
- macOS：由 @nidib 在[#4588](https://github.com/wailsapp/wails/pull/4588)中於選單列顯示原生視窗控制項
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中新增 macOS Dock 服務，以在 Dock 中隱藏／顯示應用程式圖示
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中新增 macOS Dock 服務，以在 Dock 中隱藏／顯示應用程式圖示
- 由 @leaanthony 在[#4534](https://github.com/wailsapp/wails/pull/4534)中新增 macOS 原生 Liquid Glass 效果支援，使用 NSGlassEffectView（macOS 15.0+）並以 NSVisualEffectView 作為備援，且提供完整的材質自訂選項
- 由 @leaanthony 在[#4500](https://github.dev/wailsapp/wails/pull/4500)中新增瀏覽器 URL 清理功能。以 @APShenkin 的[#4484](https://github.com/wailsapp/wails/pull/4484)為基礎。
- 由[@leaanthony](https://github.com/leaanthony)新增 Windows/Mac 內容保護功能，以[@Taiterbase](https://github.com/Taiterbase)在此[PR](https://github.com/wailsapp/wails/pull/4241)中的原始工作為基礎
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4488)中新增支援，可透過`wails3 build`與`wails3 package`別名將 CLI 變數傳遞給 Task 命令（#4422）
- 由[@atterpac](https://github.com/atterpac)在[#4318](https://github.com/wailsapp/wails/pull/4318)中新增置放區支援，事件來源會提供被置放元素的資料
- 在[PR](https://github.com/wailsapp/wails/pull/4467)中，為`WindowsWindow`選項新增`AdditionalLaunchArgs`，以便將額外的命令列引數傳遞給 WebView2 瀏覽器。
- 由[@triadmoko](https://github.com/triadmoko)在[PR](https://github.com/wailsapp/wails/pull/4286)中新增在 wails init 後自動執行 go mod tidy 的功能
- 由 @leaanthony 在[PR](https://github.dev/wailsapp/wails/pull/4463)中新增 Windows Snap Assist 功能
- 在[PR](https://github.com/wailsapp/wails/pull/4467)中，為`WindowsWindow`選項新增`AdditionalLaunchArgs`，以便將額外的命令列引數傳遞給 WebView2 瀏覽器。
- 由[@triadmoko](https://github.com/triadmoko)在[PR](https://github.com/wailsapp/wails/pull/4286)中新增在 wails init 後自動執行 go mod tidy 的功能
- 由 @leaanthony 在[PR](https://github.dev/wailsapp/wails/pull/4463)中新增 Windows Snap Assist 功能
- 由[@almas-x](https://github.com/almas-x)在[PR](https://github.com/wailsapp/wails/pull/4427)中新增 Windows `getAccentColor`實作
- 由[@almas-x](https://github.com/almas-x)在[PR](https://github.com/wailsapp/wails/pull/4427)中新增 Windows `getAccentColor`實作
- Windows 深色主題選單與選單列。由 @leaanthony 在[a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)中實作
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4405)中重新命名內建服務，使 JS/TS 繫結更清楚
- 由[@etesam913](https://github.com/etesam913)新增`app.Env.GetAccentColor`，用於取得使用者系統的輔色。支援 MacOS。
- 由[@atterpac](https://github.com/atterpac)在[#4137](https://github.com/wailsapp/wails/pull/4137)中新增`window.ToggleFrameless()` API
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4345)中新增 Linux 各發行版專用的建置相依套件
- 由 @atterpac 在[PR](https://github.com/wailsapp/wails/pull/4404)中新增繫結指南
- **重新整理測試基礎架構**：由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中將 Docker 測試檔案移至專用的`test/docker/`目錄，並最佳化映像檔及提升建置可靠性
- **改進資源管理模式**：由[@leaanthony](https://github.com/leaanthony)在範例中加入正確的事件處理常式清理，以及可感知 context 的 goroutine 管理（[#4359](https://github.com/wailsapp/wails/pull/4359)）
- 由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中支援建置 aarch64 AppImage
- 由[@leaanthony](https://github.com/leaanthony)為`wails doctor`新增診斷章節
- 由[@leaanthony](https://github.com/leaanthony)在呼叫服務方法時，將視窗加入 context
- 由[@leaanthony](https://github.com/leaanthony)新增`window-call`範例，示範如何得知是哪個視窗正在呼叫服務
- 由[@leaanthony](https://github.com/leaanthony)新增選單指南
- 由[@leaanthony](https://github.com/leaanthony)改進 panic 處理
- 由[@leaanthony](https://github.com/leaanthony)新增選單指南
- 由[@fbbdev](https://github.com/fbbdev)在[#4024](https://github.com/wailsapp/wails/pull/4024)中新增 Service API 的文件註解
- 由[@leaanthony](https://github.com/leaanthony)在[#4024](https://github.com/wailsapp/wails/pull/4024)中新增`application.NewServiceWithOptions`函式，以使用額外設定初始化服務
- 由[@FalcoG](https://github.com/FalcoG)與[@leaanthony](https://github.com/leaanthony)在[#4031](https://github.com/wailsapp/wails/pull/4031)中改進選單控制
- 由[@leaanthony](https://github.com/leaanthony)新增更多文件
- 由[@leaanthony](https://github.com/leaanthony)支援在標準事件接聽器中取消事件
- 由[@leaanthony](https://github.com/leaanthony)新增系統匣`Hide`、`Show`及`Destroy`支援
- 由[@leaanthony](https://github.com/leaanthony)新增系統匣`SetTooltip`支援。原始構想來自[@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中，讓繫結產生器針對不支援型別的警告回報套件路徑
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中新增繫結產生器對泛型別名的支援
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中新增繫結產生器對`omitzero` JSON 旗標的支援
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中新增`//wails:ignore`指示詞，以避免為所選服務方法產生繫結
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中為服務與模型新增`//wails:internal`指示詞，以允許型別在 Go 中匯出、但不在 JS/TS 中匯出
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中新增繫結產生器對別名型別常數的支援，以允許弱型別列舉
- 由[@fbbdev](https://github.com/fbbdev)在[#4068](https://github.com/wailsapp/wails/pull/4068)中新增針對 Go 1.24功能的繫結產生器測試
- 在[#4065](https://github.com/wailsapp/wails/pull/4065)中為`OSInfo.Branding`新增對 macOS 15「Sequoia」的支援，以改進作業系統版本偵測
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中新增`PostShutdown`掛鉤，以在關閉程序完成後執行自訂程式碼
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中新增`FatalError`結構，以支援在自訂錯誤處理常式中偵測嚴重錯誤
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中將服務啟動與關閉順序標準化並加以記錄
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中新增應用程式啟動／關閉序列的測試框架，以及服務啟動／關閉測試
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中新增`RegisterService`方法，以在應用程式建立後註冊服務
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中，於應用程式與服務選項新增`MarshalError`欄位，以自訂繫結呼叫的錯誤處理
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中新增可取消的 Promise 包裝函式，可沿 Promise 鏈傳遞取消要求
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中新增將繫結呼叫的取消作業繫結至`AbortSignal`的功能
- 由[@leaanthony](https://github.com/leaanthony)讓 WML 除了一般的`wml-*`屬性外，也支援`data-wml-*`屬性
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為所有服務新增`Configure`方法，以支援延後設定／動態重新設定
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中，讓`fileserver`服務在未設定時傳送503「服務無法使用」回應
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中，讓`kvstore`服務在未設定時預設提供記憶體內鍵值儲存區
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`kvstore`服務新增`Load`方法，以在設定變更後從檔案重新載入資料
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`kvstore`服務新增`Clear`方法，以刪除所有鍵
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`log`服務新增`Level`型別，以提供 JS 端的記錄層級常數
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`log`服務新增`Log`方法，以動態指定記錄層級
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中，讓`sqlite`服務在未設定時預設提供記憶體內資料庫
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`sqlite`服務新增`Close`方法，以便手動關閉資料庫
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`sqlite`服務的查詢方法新增取消支援
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`sqlite`服務新增附帶 JS 繫結的預備陳述式支援
- 由[Lea Anthony](https://github.com/leaanthony)在[PR](https://github.com/wailsapp/wails/pull/3537)中新增 Gin 支援，基於[@AnalogJ](https://github.com/AnalogJ)在此[PR](https://github.com/wailsapp/wails/pull/3537)中的原始成果
- 由[@oSethoum](https://github.com/osethoum)在[#4134](https://github.com/wailsapp/wails/pull/4134)中修正自動儲存和密碼自動儲存一律啟用的問題
- 由[@leaanthony](https://github.com/leaanthony)為視窗新增`SetMenu()`，以便在視窗上設定選單
- 由[@popaprozac](https://github.com/popaprozac)在[#4098](https://github.com/wailsapp/wails/pull/4098)中新增通知支援
-  由[@wimaha](https://github.com/wimaha)在[#4177](https://github.com/wailsapp/wails/pull/4177)中新增 Mac 檔案關聯支援
- 由[@leaanthony](https://github.com/leaanthony)新增用於遞增語意版本的`wails3 tool version`
- 由[@popaprozac](https://github.com/popaprozac)在[#](https://github.com/wailsapp/wails/pull/4234)中新增 macOS 和 Windows 的徽章支援
- 由[@fbbdev](https://github.com/fbbdev)與[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中新增已註冊／強型別事件支援
- 由[@fbbdev](https://github.com/fbbdev)與[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中新增為自訂事件註冊掛鉤的功能
- 由[@Krzysztofz01](https://github.com/Krzysztofz01)、[@rcalixte](https://github.com/rcalixte)新增`app.OpenFileManager(path string, selectFile bool)`，以在系統檔案管理員中開啟路徑`path`，並可選擇透過`selectFile`醒目提示
- 由[@leaanthony](https://github.com/leaanthony)為`wails3 init`命令新增`-git`旗標
- 由[@leaanthony](https://github.com/leaanthony)新增`wails3 generate webview2bootstrapper`命令
- 由[@leaanthony](https://github.com/leaanthony)在執行階段中新增`init()`方法，以便手動初始化執行階段
- 由[@leaanthony](https://github.com/leaanthony)在 Window 的 WindowOptions 中新增`WindowDidMoveDebounceMS`選項
- 由[@leaanthony](https://github.com/leaanthony)新增單一執行個體功能。此功能以 @APshenkin 的[v2 PR](https://github.com/wailsapp/wails/pull/2951)為基礎。
- 由[@leaanthony](https://github.com/leaanthony)新增`wails3 generate template`命令
- 由[@leaanthony](https://github.com/leaanthony)新增`wails3 releasenotes`命令
- 由[@leaanthony](https://github.com/leaanthony)新增`wails3 update cli`命令
- 由[@leaanthony](https://github.com/leaanthony)為`wails3 generate bindings`命令新增`-clean`選項
- 由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中允許建置適用於 aarch64 (arm64) 的 Linux AppImage
- 由 @ansxuman 在[#3958](https://github.com/wailsapp/wails/pull/3958)中新增贊助者超連結
- 由以下人員新增 Linux 的 deb、rpm 和 Arch Linux 套件建置支援
- 由以下人員新增 Darwin 通用建置與套件支援
- 由以下人員將事件文件新增至網站
- 新增設定為非 SSR 開發的 sveltekit 和 sveltekit-ts 範本
- 由以下人員使用新的`wails3 update build-assets`命令更新建置資產
- 由以下人員新增用於測試 HTML 拖放 API 的範例
- 由[leaanthony](https://github.com/leaanthony)在以下項目中新增檔案關聯支援
- 由以下人員新增`wails3 generate runtime`命令
- 新增`InitialPosition`選項，用於指定視窗是否應置中，或
- 由以下人員在`application`套件中新增`Path`與`Paths`方法
- 新增`GeneralAutofillEnabled`和`PasswordAutosaveEnabled` Windows 選項
- 由以下人員新增擷取呼叫服務方法之視窗的功能
- 由以下人員新增 WebView2 的`EnabledFeatures`和`DisabledFeatures`選項
- ⊞ 由以下人員新增用於強化高 DPI 顯示器支援的 DIP 系統
- ⊞ 由[windom](https://github.com/windom/)在以下項目中新增視窗類別名稱選項
- 服務已擴充，可提供外掛程式功能。作者為
- 🐧 以下項目中的 WindowDidMove / WindowDidResize 事件
- ⊞ 以下項目中的 WindowDidResize 事件
-  新增 ApplicationShouldHandleReopen 事件，以便處理 Dock
-  由 @tmclane 在以下項目中將 getPrimaryScreen/getScreens 新增至實作
-  由以下人員新增在 macOS 全螢幕模式中顯示工具列的選項
- 🐧 新增 onKeyPress 邏輯，將 Linux 按鍵事件轉換為快速鍵
- 🐧 由以下人員新增工作`run:linux`
- 由[@almas-x](https://github.com/almas-x)在以下項目中匯出`SetIcon`方法
- 由[@almas-x](https://github.com/almas-x)在以下項目中改善`OnShutdown`
- 由以下人員在`Window`介面中還原`ToggleMaximise`方法
- 由 @leaanthony 在以下項目中為`Environment()`新增更多資訊
- 由以下人員在`Window`介面上公開`WebviewWindow.IsFocused`方法
- 由以下人員在 WML 系統中支援以空格分隔的多個觸發事件
- 由以下人員從隨附的 JS 執行階段指令碼新增 ESM 匯出
- 新增繫結產生器旗標，以使用隨附的 JS 執行階段指令碼，而非
- 由[@abichinger](https://github.com/abichinger)在 Linux 上實作`setIcon`
- 在 dev 命令中新增`-port`旗標，並支援環境變數
- 新增繫結方法呼叫的測試，由
- ⊞ 為已建立的視窗新增`SetIgnoreMouseEvents`，由
-  新增設定視窗堆疊層級（順序）的功能，由

### 修正

- 修正在 Retina Mac 上`Screen.Bounds`、`WorkArea`和`Size`減半的問題：將 NSScreen 的點數值轉換為`Physical*`欄位中的裝置像素，並填入頂層`Screen.X`／`Y`，使[PR](https://github.com/wailsapp/wails/pull/5168)中的多顯示器接觸偵測和工作區域放置正確，由 @wayneforrest 貢獻
- 修正 ScreenManager 中的資料競爭；該問題會在顯示器設定變更時（例如睡眠／喚醒期間熱插拔外接顯示器）造成 WebKit DisplayLink 死結
- 當 Assets.car 存在時，直接將 CFBundleIconName 設為 appicon，見[PR](https://github.com/wailsapp/wails/pull/5154)，由 @symball 貢獻
- 修正`wails3 doctor`在 Fedora、openSUSE、Arch 和 NixOS 上回報錯誤 WebKitGTK 套件的問題——由於 v3 在編譯時需要4.1 API，因此已移除4.0備援項目（#5071）
- 修正 openSUSE webkit2gtk doctor 套件名稱（`webkit2gtk4_1-devel` → `webkit2gtk3-devel`，後者才是正確的 openSUSE 套件名稱）（#5071）
- 修正在桌面開發模式中缺少`/wails/custom.js`時發生的`Unexpected token '<'`錯誤。為`/wails/custom.js`新增明確的404處理常式，並在`loadOptionalScript`中新增不區分大小寫的`Content-Type`驗證，以防止將 HTML SPA 備援內容注入為 JavaScript。（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- 修正 macOS 上系統匣選單的醒目提示狀態——現在開啟選單時，圖示會顯示為已選取狀態（#4910）
- 修正 macOS 上附加至系統匣的視窗出現在其他視窗後方的問題——現在會使用正確的快顯視窗層級（#4910）
- 修正文件中錯誤的`@wailsio/runtime`匯入範例（#4989）
- 修正在 darwin 上無法最小化無框視窗的問題（#4294）
- 將`node_modules/`排除於 go-task 的最新狀態檢查之外，以修正執行`wails3 build`和`wails3 dev`時長達20-30分鐘的停滯。先前`sources: "**/*"`萬用字元模式會使 go-task 列舉`node_modules/`中的每個檔案並計算其總和檢查碼（使用 MUI 等大型相依套件時，檔案數可達50000-100000以上），在 Windows/NTFS 上尤其緩慢（#4939）
- 修正 C 的`Screen` typedef 與 X11 Xlib.h 衝突所造成的 GTK4 建置失敗（#4957）
- 修正 macOS 上 Dock 標記方法的一致性
- 修正`InvisibleTitleBarHeight`套用至所有 macOS 視窗，而非僅套用至無框或透明標題列視窗的問題（#4960）
- 啟用`InvisibleTitleBarHeight`時，在視窗邊緣附近略過拖曳啟動，以修正從頂端角落調整大小時視窗抖動／晃動的問題（#4960）
- 修正 JS/TS 繫結中使用列舉鍵的映射型別產生問題（#4437），由 @fbbdev 貢獻
- 修正在 Windows 上，顯示縮放比例不是100% 時檔案拖放無法運作的問題
- 修正在 Windows 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 修正 Windows 上檔案放置座標使用錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正在 Linux 上搭配懸停效果時，檔案拖放無法可靠運作的問題
- 修正在 Linux 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 透過使用`gtk_window_present()`，修正 Linux/GTK4 上顯示／隱藏視窗有時會還原為最小化狀態的問題（#4957）
- 透過`XTranslateCoordinates`／`XMoveWindow`新增以 X11 為條件的支援，修正 Linux/GTK4 上取得／設定視窗位置總是傳回0,0的問題（#4957）
- 新增以訊號為基礎的大小限制來取代已移除的`gtk_window_set_geometry_hints`，修正 Linux/GTK4 上未強制執行視窗大小上限的問題（#4957）
- 透過`gdk_monitor_get_scale`（GTK 4.14+）實作正確的 PhysicalBounds 計算及非整數縮放支援，修正 Linux/GTK4 上的 DPI 縮放
- 修正在 Linux/GTK4 上建立新視窗時選單項目重複的問題
- 修正 JS/TS 繫結中使用列舉鍵的映射型別產生問題（#4437），由 @fbbdev 貢獻
- 修正在 Windows 上，顯示縮放比例不是100% 時檔案拖放無法運作的問題
- 修正在 Windows 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 修正 Windows 上檔案放置座標使用錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正在 Linux 上搭配懸停效果時，檔案拖放無法可靠運作的問題
- 修正在 Linux 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 透過`gdk_monitor_get_scale`（GTK 4.14+）實作正確的 PhysicalBounds 計算及非整數縮放支援，修正 Linux/GTK4 上的 DPI 縮放
- 修正在 Linux/GTK4 上建立新視窗時選單項目重複的問題
- 修正 JS/TS 繫結中使用列舉鍵的映射型別產生問題（#4437），由 @fbbdev 貢獻
- 修正 App.Window.Current() 未從主執行緒存取 AppKit API 而在 macOS 上造成的「幽靈視窗」問題（#4947），由 @wimaha 貢獻
- 透過實作 WKUIDelegate runOpenPanelWithParameters，修正 HTML `<input type="file">`在 macOS 上無法運作的問題（#4862）
- 修正在 macOS/Linux 上使用`@wailsio/runtime` npm 模組時，原生檔案拖放無法運作的問題（#4953），由 @leaanthony 貢獻
- 修正跨套件型別別名的繫結產生問題（#4578），由 @fbbdev 貢獻
- 修正 Linux 上因違反 GTK 執行緒安全規則而造成 OpenFileDialog 當機的問題（#3683），由 @ddmoney420 貢獻
- 修正在隱藏或已銷毀的視窗上呼叫`Focus()`時發生 SIGSEGV 當機的問題（#4890），由 @ddmoney420 貢獻
- 修正在 Linux 上設定空白圖示或點陣圖時可能發生 panic 的問題（#4923），由 @ddmoney420 貢獻
- 修正在 macOS 上從服務繫結呼叫 ErrorDialog 時發生當機的問題（#3631），由 @leaanthony 貢獻
- 使選單可在 Windows 作業系統的`v3\examples\dialogs`中顯示，由 @ndianabasi 貢獻
- 修正頁面重新載入期間因競爭條件而發生 TypeError 的問題（#4872），由 @ddmoney420 貢獻
- 移除`Collector.IsVoidAlias()`方法中的全域狀態，以修正繫結產生器測試的錯誤輸出（#4941），由 @fbbdev 貢獻
- 修正`<input type="file">`檔案選擇器在 macOS 上無法運作的問題（#4862），由 @leaanthony 貢獻
- 修正在 macOS 上`Position()`與`SetPosition()`使用不一致的座標系統，導致儲存／還原狀態時視窗位置偏移的問題（#4816），由 @leaanthony 修正
- 修正已透過應用程式資訊清單設定 DPI 感知時，SetProcessDpiAwarenessContext 發生「存取遭拒」錯誤的問題（#4803）
- 更新鍵盤快速鍵文件頁面，並修正`KeyBinding.Add`回呼參數的型別，由 @ndianabasi 完成
- 修正產生自訂繫結的相關文件；必須使用`-d String`，而非`-o String`
- 修正呼叫`menu.Update()`時選單未清除子項目的問題
- 修正文件中已過時的 Manager API 參照（更新31檔案，改用`app.Window.New()`、`app.Event.Emit()`等新模式），由 @leaanthony 完成
- 修正 WebKit 覆寫訊號處理常式，導致繫結至 JS 的 Go 方法發生 panic 時 Linux 當機的問題（#3965），由 @leaanthony 修正
- 修正 SaveFileDialog.SetFilename() 在 Linux 上無效的問題（#4841），由 @samstanier 修正
- 修正拖放範例中的放置座標顯示為 undefined 的問題
- 修正 APP_NAME 包含空格時無法建立 macOS 應用程式套件的問題（大括號展開問題）
- 修正在 Windows 上呼叫服務方法時發生索引超出範圍 panic 的問題（還原 goccy/go-json）
- 修正在 Windows 上使用非100% 顯示縮放比例時，檔案拖放無法運作的問題
- 修正在 Windows 上啟用檔案放置時，HTML5 內部拖放功能失效的問題
- 修正 Windows 上的檔案放置座標位於錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正在 Linux 上使用游標懸停效果時，檔案拖放無法可靠運作的問題
- 修正在 Linux 上啟用檔案放置時，HTML5 內部拖放功能失效的問題
- 更新所有作業系統之 Taskfile.yml 檔案中的全部命令，以支援`APP_NAME`等變數內的空格，由 @ndianabasi 完成
- 修正在 Linux 上執行 'build:universal:lipo:go' 工作時的命令引數錯誤，由 @wux1an 修正
- 修正在 Linux 上執行 'wails3 build GOOS=darwin GOARCH=arm64' 時發生的 Docker 錯誤「undefined symbol: **<em>ubsan</em>handle_xxxxxxx」，由 @wux1an 修正
- 整合自訂通訊協定文件並新增 Universal Links 章節，由 @leaanthony 完成
- 新增防止並行呼叫 TrackPopupMenuEx 的防護措施，修正在 Windows 上重複點按系統匣圖示時選單當機的問題（#4151），由 @leaanthony 修正
- 防止在 app.Run() 之前呼叫 systray.Run() 時應用程式當機，由 @leaanthony 完成
- 修正在 macOS 上啟用 ApplicationShouldTerminateAfterLastWindowClosed 時，透過 Hide()/Show() 切換視窗可見性會導致應用程式當機的問題（#4389），由 @leaanthony 修正
- 修正在 macOS 與 Windows 上重複開啟內容選單時發生記憶體洩漏的問題（#4012），由 @leaanthony 修正
- 修正 macOS 未重複使用內容選單的原生資源，導致每次顯示時都建立新選單的問題（#4012），由 @leaanthony 修正
- 修正 macOS 應用程式以`Hidden: true`啟動時，點按 Dock 圖示不會顯示隱藏視窗的問題（#4583），由 @leaanthony 修正
- 修正 CGO 呼叫使用錯誤的視窗指標型別，導致 macOS 列印對話方塊無法開啟的問題（#4290），由 @leaanthony 修正
- 修正 appmenu-gtk-module 存取尚未具現化的視窗，導致 Wayland 上的視窗選單當機的問題（#4769），由 @leaanthony 修正
- 修正應用程式名稱包含無效字元（空格、括號等）時 GTK 應用程式當機的問題，由 @leaanthony 修正
- 修正在 Windows 上初始化拖放功能時發生「記憶體不足」錯誤的問題（#4701），由 @overlordtm 修正
- 修正 URI 跳脫不正確，導致 Linux 上的檔案總管開啟錯誤目錄的問題（#4397），由 @leaanthony 修正
- 透過自動偵測`.relr.dyn` ELF 區段並停用剝除，修正在現代 Linux 發行版（Arch、Fedora 39+、Ubuntu 24.04+）上建置 AppImage 失敗的問題（#4642），由 @leaanthony 修正
- 修正`wails doctor`在 Fedora／以 DNF 為基礎的系統上錯誤回報已安裝 webkit 套件的問題（#4457），由 @leaanthony 修正
- 修正預設`config.yml`會使用正式環境組建執行`wails3 dev`的問題，由 @mbaklor 修正
- 修正 iOS 服務存根因匯入不存在的套件而導致建置失敗的問題，由 @leaanthony 修正
- 修正 debug/info 方法中的結構化記錄導致「no formatting directives」錯誤的問題，由 @leaanthony 修正
- 移除合併行動平台程式碼時意外納入的暫時偵錯列印陳述式，由 @leaanthony 完成
- 透過自動停用 DMA-BUF 轉譯器，修正在搭載 NVIDIA GPU 的 Wayland 上 WebKitGTK 當機的問題（錯誤 71 Protocol error），由 @leaanthony 修正
- 解決 Linux 上`application.WebviewWindowOptions.BackgroundColour`忽略 Alpha 值的問題（[#4722](https://github.com/wailsapp/wails/pull/4722)，@BradHacker）
- 修正未提供自訂圖示時，Windows 系統匣圖示不會預設使用應用程式圖示的問題（#4704）
- 追蹤`HICON`的擁有權，僅銷毀使用者建立的控制代碼，以防止 Explorer 回收時當機（#4653）。
- 在銷毀期間釋放 Windows 系統佈景主題監聽器與保留的系統匣圖示，以免持續洩漏 goroutine 與裝置內容（#4653）。
- 將系統匣工具提示截斷為127個 UTF-16單位，以免損壞代理對與多位元組字符（#4653）。
- 修正 Windows 套件工作失敗的問題（#4667）
- 修正 Linux taskfile 中的 Linux AppImage appicon 變數，[PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 修正 go-webview2 v1.0.22簽章變更所造成的 Windows 建置錯誤（#4513、#4645）
- 修正 Linux taskfile 中的 Linux AppImage appicon 變數，[PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 將`<.Info.Protocol>`改為`<.Protocol>`，修正 Linux desktop.tmpl 的通訊協定範圍，由 @Tolfx 於 #4510完成
- 修正 liquid glass 示範中的重複定義錯誤，[#4542](https://github.com/wailsapp/wails/pull/4542)，由 @Etesam913 修正
- 修正 Linux 上的系統匣選單更新，[#4604](https://github.com/wailsapp/wails/issues/4604)，由[@JackDoan](https://github.com/JackDoan)修正
- 修正在 Windows 上建立隱藏視窗時出現白色視窗的問題，由 @leaanthony 於[#4612](https://github.com/wailsapp/wails/pull/4612)修正
- 修正文件中的 notifications 套件匯入路徑，由 @rxliuli 於[#4617](https://github.com/wailsapp/wails/pull/4617)修正
- 修正在使用 npm 套件 @wailsio/runtime 時拖放功能無法運作的問題（#4489），由 @leaanthony 於 #4616 中修正
- Windows：修正啟動時視窗閃爍，以及隱藏視窗錯誤顯示的問題；由 @leaanthony 於 [PR](https://github.com/wailsapp/wails/pull/4600) 中修正。
- 修正 Wayland 視窗最大化時的尺寸問題（https://github.com/wailsapp/wails/issues/4429），由 [@samstanier](https://github.com/samstanier) 修正
- 修正 Wayland 視窗最大化時的尺寸問題（https://github.com/wailsapp/wails/issues/4429），由 [@samstanier](https://github.com/samstanier) 修正
- 修正 liquid glass 示範程式中的重複定義錯誤；由 @Etesam913 於 [#4542](https://github.com/wailsapp/wails/pull/4542) 中修正
- 修正 AssetServer 在 MacOS 上可能當機的問題；由 @jghiloni 於 [#4576](https://github.com/wailsapp/wails/pull/4576) 中修正
- 修正使用 NextJs 建置時的編譯問題。由 @rev42 於 [#4585](https://github.com/wailsapp/wails/pull/4585) 中修正
- 修正每夜版本的發布管線；由 @riadafridishibly 於 [#4597](https://github.com/wailsapp/wails/pull/4597) 中修正
- 修正 liquid glass 示範程式中的重複定義錯誤；由 @Etesam913 於 [#4542](https://github.com/wailsapp/wails/pull/4542) 中修正
- 修正 AssetServer 在 MacOS 上可能當機的問題；由 @jghiloni 於 [#4576](https://github.com/wailsapp/wails/pull/4576) 中修正
- 修正使用 NextJs 建置時的編譯問題。由 @rev42 於 [#4585](https://github.com/wailsapp/wails/pull/4585) 中修正
- 修正每夜版本的發布管線；由 @riadafridishibly 於 [#4597](https://github.com/wailsapp/wails/pull/4597) 中修正
- 修正 liquid glass 示範程式中的重複定義錯誤；由 @Etesam913 於 [#4542](https://github.com/wailsapp/wails/pull/4542) 中修正
- 修正 Windows 上的 SetBackgroundColour；由 @PPTGamer 於 [PR](https://github.com/wailsapp/wails/pull/4492) 中修正
- 更新文件以反映 Manager API 重構所帶來的變更；由 @yulesxoxo 於 [PR #4476](https://github.com/wailsapp/wails/pull/4476) 中更新
- 修正 Linux taskfile 中 Linux .desktop 檔案的 appicon 變數；見 [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- 更新文件以反映 Manager API 重構所帶來的變更；由 @yulesxoxo 於 [PR #4476](https://github.com/wailsapp/wails/pull/4476) 中更新
- 由 @leaanthony 在[#4460](https://github.com/wailsapp/wails/pull/4460)中修正[#4456](https://github.com/wailsapp/wails/issues/4456)所回報的 Windows nil 指標解參考錯誤
- 在 macOS WKWebView 中新增對 `allowsBackForwardNavigationGestures` 的支援，以啟用雙指滑動導覽手勢（#1857）
- 修正初始設為停用的選單項目無法使用 onClick 的問題；由 @leaanthony 於 [PR #4469](https://github.com/wailsapp/wails/pull/4469) 中修正。感謝 @IanVS 進行初步調查。
- 修正建置失敗時未清理 Vite 伺服器的問題（#4403）
- 修正在 windows 上關閉或取消 `SaveFileDialog` 時發生的 panic；由 @hkhere 於 [PR](https://github.com/wailsapp/wails/pull/4284) 中修正
- 修正 Windows 上 HTML 層級的拖放功能；由 [@mbaklor](https://github.com/mbaklor) 於 [#4259](https://github.com/wailsapp/wails/pull/4259) 中修正
- 在 macOS WKWebView 中新增對 `allowsBackForwardNavigationGestures` 的支援，以啟用雙指滑動導覽手勢（#1857）
- 修正初始設為停用的選單項目無法使用 onClick 的問題；由 @leaanthony 於 [PR #4469](https://github.com/wailsapp/wails/pull/4469) 中修正。感謝 @IanVS 進行初步調查。
- 修正建置失敗時未清理 Vite 伺服器的問題（#4403）
- 修正 Windows 上的通知剖析問題；由 @popaprozac 於 [PR](https://github.com/wailsapp/wails/pull/4450) 中修正
- 修正 doctor 命令，使其檢查 Windows SDK 相依項目；由 [@kodumulo](https://github.com/kodumulo) 於 [#4390](https://github.com/wailsapp/wails/issues/4390) 中修正
- 修正 Mac 上 processURLRequest 中的 nil 指標解參考錯誤；由 [@etesam913](https://github.com/etesam913) 於 [#4366](https://github.com/wailsapp/wails/pull/4366) 中修正
- 修正導致無法使用篩選對話方塊的 linux 錯誤；由 [@bh90210](https://github.com/bh90210) 於 [#4287](https://github.com/wailsapp/wails/pull/4287) 中修正
- 修正 Windows 與 Linux 上的「編輯」選單問題；由 [@leaanthony](https://github.com/leaanthony) 於 [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062) 中修正
- 將 macOS .plist 檔案中的最低系統版本由 10.13.0 更新為 10.15.0；由 [@AkshayKalose](https://github.com/AkshayKalose) 於 [#3981](https://github.com/wailsapp/wails/pull/3981) 中更新
- 修正視窗 ID 跳號問題；由 [@leaanthony](https://github.com/leaanthony) 修正
- 修正呼叫 RegisterContextMenu 時的 nil 選單問題；由 [@leaanthony](https://github.com/leaanthony) 修正
- 修正繫結產生器輸出中的相依性循環；由 [@fbbdev](https://github.com/fbbdev) 於 [#4001](https://github.com/wailsapp/wails/pull/4001) 中修正
- 修正繫結產生器輸出中先使用後定義的錯誤；由 [@fbbdev](https://github.com/fbbdev) 於 [#4001](https://github.com/wailsapp/wails/pull/4001) 中修正
- 將建置旗標傳遞給繫結產生器；由 [@fbbdev](https://github.com/fbbdev) 於 [#4023](https://github.com/wailsapp/wails/pull/4023) 中完成
- 將 windows Taskfile 中的路徑改為使用正斜線，以確保其可在非 Windows 平台上運作；由 [@leaanthony](https://github.com/leaanthony) 完成
- 現已修正 Mac 與 Mac JS 事件；由 [@leaanthony](https://github.com/leaanthony) 修正
- 修正 macOS 上的事件死結；由 [@leaanthony](https://github.com/leaanthony) 修正
- 修正 Windows 上初始化 Window 時，在提供 HTML 但未提供 JS 的情況下發生的 `Parameter incorrect` 錯誤；由 [@leaanthony](https://github.com/leaanthony) 修正
- 修正資產伺服器中用於偵測內容類型的回應前綴大小；由 [@fbbdev](https://github.com/fbbdev) 於 [#4049](https://github.com/wailsapp/wails/pull/4049) 中修正
- 修正資產伺服器在根索引路徑上對非 404 回應的處理；由 [@fbbdev](https://github.com/fbbdev) 於 [#4049](https://github.com/wailsapp/wails/pull/4049) 中修正
- 修正繫結產生器測試泛型型別的屬性時出現的未定義行為；由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中修正
- 修正當基礎型別與具名包裝型別的屬性不相同時，繫結產生器對模型產生的輸出；由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中修正
- 修正繫結產生器針對映射鍵型別與預處理所產生的輸出；由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中修正
- 修正繫結產生器針對實作 marshaler 介面的結構所產生的輸出；由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中修正
- 修正繫結產生器對涉及泛型型別之型別循環的偵測；由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中修正
- 修正繫結產生器輸出中對未匯出模型的無效參照，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 完成
- 將注入的程式碼移至服務檔案末尾，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 完成
- 修正繫結產生器處理檔案關閉作業錯誤的方式，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 完成
- 針對定義了生命週期或 HTTP 方法、但未定義其他繫結方法的服務隱藏警告，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 完成
- 修正非 React 範本在使用淺色系統配色時無法顯示 Hello World 頁尾的問題，由 [@marcus-crane](https://github.com/marcus-crane) 於 [#4056](https://github.com/wailsapp/wails/pull/4056) 完成
- 修正 macOS 上隱藏選單項目的問題，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修正訊息處理器中的錯誤處理與格式化，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 完成
-  修正結束應用程式時略過服務關閉程序的問題，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 完成
-  確保選單更新在主執行緒上執行，由 [@leaanthony](https://github.com/leaanthony) 完成
- 拖曳與調整大小機制現在更加穩健，且更符合各平台的預期行為，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 完成
- 修正 [#4097](https://github.com/wailsapp/wails/issues/4097) 中 Webpack/angular 捨棄執行階段初始化程式碼的問題，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 完成
- 修正初始為隱藏狀態的選單項目，由 [@IanVS](https://github.com/IanVS) 於 [#4116](https://github.com/wailsapp/wails/pull/4116) 完成
- 修正請求路徑沒有副檔名，且`[request]`不存在、但`[request].html`存在時，assetFileServer 未提供`.html`檔案的問題
- 修正圖示產生路徑，由 [@robin-samuel](https://github.com/robin-samuel) 於 [#4125](https://github.com/wailsapp/wails/pull/4125) 完成
- 修正未發出 `fullscreen`、`unfullscreen`、`unminimise` 與 `unmaximise` 事件的問題，由 [@oSethoum](https://github.com/osethoum) 於 [#4130](https://github.com/wailsapp/wails/pull/4130) 完成
- 修正設定中預設版本的前綴不正確所造成的 NSIS 錯誤，由[@robin-samuel](https://github.com/robin-samuel)於[#4126](https://github.com/wailsapp/wails/pull/4126)完成
- 修正 Windows 上 Dialogs 執行階段函式傳回已逸出的路徑，由 [TheGB0077](https://github.com/TheGB0077) 於 [#4188](https://github.com/wailsapp/wails/pull/4188) 完成
- 修正 HKCU 中的 Webview2 偵測路徑，由 [@leaanthony](https://github.com/leaanthony) 完成。
- 修正 macOS 上的輸入問題，由 [@leaanthony](https://github.com/leaanthony) 完成。
- 修正 Windows 圖示產生工作的檔案名稱，由 [@yulesxoxo](https://github.com/yulesxoxo) 於 [#4219](https://github.com/wailsapp/wails/pull/4219) 完成。
- 修正無框視窗的透明度問題，由 [@leaanthony](https://github.com/leaanthony) 根據 @kron 的工作完成。
- 修正視窗停用或最小化時的焦點呼叫，由 [@leaanthony](https://github.com/leaanthony) 根據 @kron 的工作完成。
- 修正工作列重新啟動後系統匣未顯示的問題，由 [@leaanthony](https://github.com/leaanthony) 根據 @kron 的工作完成。
- 修正 fallbackResponseWriter 未實作 Flush() 的問題，見 [#4245](https://github.com/wailsapp/wails/pull/4245)
- 修正 fallbackResponseWriter 未實作 Flush() 的問題，由 [@superDingda] 於 [#4236](https://github.com/wailsapp/wails/issues/4236) 完成
- 修正 macOS 視窗在非同步 Go 繫結函式呼叫仍待處理時關閉會導致當機的問題，由 [@joshhardy](https://github.com/joshhardy) 於 [#4354](https://github.com/wailsapp/wails/pull/4354) 完成
- 修正 Windows 效率模式啟動時的競爭條件，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修正 Windows 圖示控制代碼的清理，由 [@leaanthony](https://github.com/leaanthony) 完成。
- 修正 Windows 上的 `OpenFileManager`，由 [@PPTGamer](https://github.com/PPTGamer) 於 [#4375](https://github.com/wailsapp/wails/pull/4375) 完成。
- 修正 Linux 的最小／最大寬度選項，由 @atterpac 於 [#3979](https://github.com/wailsapp/wails/pull/3979) 完成
- 透過提高 npm 版本，修正 TypeScript 範本的型別定義，由 @atterpac 於 [#3966](https://github.com/wailsapp/wails/pull/3966) 完成
- 修正 SvelteKit 範本的 CSS 參照，由 @atterpac 於 [#3945](https://github.com/wailsapp/wails/pull/3945) 完成
- 確保 window run() 中的關鍵回呼在主執行緒上呼叫，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修正對話方塊目錄選擇器範例，由 [@leaanthony](https://github.com/leaanthony) 完成
- 新增 index.html 遺失時顯示的中文錯誤頁面，由 [@leaanthony](https://github.com/leaanthony) 完成
-  確保 `windowDidBecomeKey` 回呼在主執行緒上執行，由 [@leaanthony](https://github.com/leaanthony) 完成
-  支援無框視窗全螢幕顯示，由 [@leaanthony](https://github.com/leaanthony) 完成
-  改善視窗銷毀邏輯，由 [@leaanthony](https://github.com/leaanthony) 完成
-  修正視窗附加至系統匣時的位置邏輯，由 [@leaanthony](https://github.com/leaanthony) 完成
-  支援無框視窗全螢幕顯示，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修正事件處理，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修正視窗關閉邏輯，由 [@leaanthony](https://github.com/leaanthony) 完成
- 共用 taskfile 現在預設會為 TypeScript 範本產生 TypeScript 繫結，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修正沒有開啟任何視窗／僅有系統匣時，收到 WM_CLOSE 訊息無法關閉應用程式的問題，由[@mmalcek](https://github.com/mmalcek)於[#3990](https://github.com/wailsapp/wails/pull/3990)完成
- 修正 garble 建置，由 @5aaee9 於 [#3192](https://github.com/wailsapp/wails/pull/3192) 完成
- 修正 Windows NSIS 建置，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修正 Linux 多選對話方塊中因未關閉而造成的死結
- 修正 Windows 建置期間跨平台清理 .syso 檔案的問題，由
- 修正 amd64 AppImage 編譯，由 @atterpac 於
- 修正建置資產更新，由 @ansxuman 於
- 由 @atterpac 修正 Linux 系統匣的 `OnClick` 與 `OnRightClick` 實作
- 修正 `AlwaysOnTop` 在 Mac 上無法運作的問題，作者：
-  修正 `application.NewEditMenu` 包含重複項目的問題
- 🐧 修正 aarch64 編譯問題
- ⊞ 修正選項按鈕群組選單項目，作者：
- 修正在 MacOS 上建置可執行的 .app 時，'name' 與 'outputfilename' 所引發的錯誤
- 修正拖放範例中使用 customEventProcessor 的錯誤，作者：
- 🐧 修正新增 IgnoreMouseEvents 所導致的 Linux 編譯錯誤，作者：
- ⊞ 修正 syso 圖示檔案的產生錯誤，作者：
- 🐧 納入原生執行於 Wayland 的修正，來源：
- 不要繫結以下項目中的內部服務方法：
- ⊞ 修正以下項目中的系統匣啟動 panic：
- 不要繫結以下項目中的內部服務方法：
- ⊞ 修正以下項目中的系統匣啟動 panic：
- 大幅重構選單項目與事件處理。目前主要改善 macOS。作者：
- 修正外掛程式與事件重構後的測試，位於：
- ⊞ 修正 `Failed to unregister class Chrome_WidgetWin_0` 警告。作者：
- 模組問題
- 由 [atterpac](https://github.com/atterpac) 修正以下項目中的調整大小事件訊息：
- 🐧 修正 NixOS 上的佈景主題處理錯誤，作者：
- 修正 Windows 上跨磁碟區安裝專案的問題，作者：
- 修正 React 範本的 CSS 以顯示頁尾，作者：
- 更新至最新版 refresh，修正開發模式下的殭屍程序問題
- 由 [Atterpac](https://github.com/atterpac) 修正 AppImage 的 WebKit 檔案來源問題
- 由 [Atterpac](https://github.com/Atterpac) 修正以下項目中的 Doctor apt 套件驗證：
- 由 @5aaee9 修正應用程式結束時凍結的問題（Darwin），位於：
- 修正 Windows 上範例的背景色彩，作者：
- 由 [mmghv](https://github.com/mmghv) 修正以下項目中的預設內容功能表：
- 修正 Darwin 上方向鍵的十六進位值，作者：
- 使 Windows 上的拖放功能正常運作。新增者：
- 修正使用者沒有適當驅動程式時，Doctor 在 Linux 上發生的錯誤
- 修正啟動時的 DPI 縮放問題（Windows）。由 [@almas-x](https://github.com/almas-x) 變更，位於：
- 修正 `go.mod` 中的取代行，使其使用相對路徑。修正包含以下內容的 Windows 路徑：
- 修正 MacOS 系統匣未連結視窗時的點擊處理，作者：
- 修正未知選項導致 Windows 建置失敗的問題，作者：
- 修正 Windows 上沒有以下項目時，以滑鼠左鍵按一下系統匣圖示所造成的當機：
- 由 @5aaee9 在 PR 中修正開啟視窗兩次時 baseURL 錯誤的問題
- 修正 `WebviewWindow.Restore` 方法中 if 分支的順序，作者：
- 在多次叫用 `GetStartURL` 時正確計算 `startURL`，條件為：
- 修正 `Screen` 結構的 JS 型別，使其與對應的 Go 型別一致，作者：
- 修正 `WML.Reload` 方法，確保已註冊的事件能正確清理
- 修正自訂內容功能表在 Linux 上立即關閉的問題，作者：
- 修正繫結所產生之模型檔案的輸出路徑與副檔名
- 修正繫結所產生之 JS 程式碼中的模型檔案匯入路徑
- 修正部分 Linux 發行版上的拖放問題，作者：
- 修正使用 `wails3 task dev` 時 macOS 缺少工作項目的問題，作者：
- 修正註冊事件造成 nil map 賦值的問題，作者：
- 修正繫結方法參數的反序列化處理，作者：
- 修正繫結方法傳回多個值時的處理，作者：
- 修正 Doctor 對未透過系統套件管理員安裝之 npm 的偵測
- 修正缺少 MicrosoftEdgeWebview2Setup.exe 的問題。感謝：
- 由 @leaanthony 修正視窗 ID 處理導致 Linux 上隨機當機的問題。依據：
- 修正 systemTray.setIcon 在 Linux 上造成當機的問題，作者：
- 修正以下平台上 `setFrameless` 函式第一次呼叫時未確實套用視窗框架的問題：

### 變更

- **破壞性變更**：產生的 JS/TS 繫結現在會將 Map 索引鍵標示為選用，以準確反映 Go map 的語意。現在於 Typescript 中存取 Map 值時會傳回 `T | undefined`，而非 `T`，因此必須進行 null 檢查或使用斷言（#4943），作者：`@fbbdev`
- 依照 `@wailsio/runtime` 中的變更，將 `Event` 的使用方式改為 `Events`，並在 `Features/Events/Event System` 的文件中改用適當的函式呼叫，作者：@AbdelhadiSeddar
- 將 `EnabledFeatures`、`DisabledFeatures` 與 `AdditionalBrowserArgs` 從各視窗選項移至應用程式層級的 `Options.Windows`（#4559），作者：@leaanthony
- 更新 `Drag N Drop` 範例的 README，並特別指出該範例示範了 `Internal Drag and Drop`，作者：@ndianabasi
- 將多項偵錯記錄的層級從 Info 改為 Debug（作者：@mbaklor）
- <strong>破壞性變更：</strong>將視窗選項中的 `EnableDragAndDrop` 重新命名為 `EnableFileDrop`
- <strong>破壞性變更：</strong>將事件內容中的 `DropZoneDetails` 重新命名為 `DropTargetDetails`
- <strong>破壞性變更：</strong>將 `WindowEventContext` 上的 `DropZoneDetails()` 方法重新命名為 `DropTargetDetails()`
- <strong>破壞性變更：</strong>移除 `WindowDropZoneFilesDropped` 事件，改用 `WindowFilesDropped`
- <strong>不向後相容的變更：</strong>將 HTML 屬性從`data-wails-dropzone`變更為`data-file-drop-target`
- <strong>不向後相容的變更：</strong>將 CSS 懸停類別從`wails-dropzone-hover`變更為`file-drop-target-active`
- <strong>不向後相容的變更：</strong>從 Windows 移除`DragEffect`、`OnEnterEffect`和`OnOverEffect`選項（這些選項原屬於已移除的 IDropTarget）
- 所有執行階段 JSON 處理（方法繫結、事件、webview 請求、通知、kvstore）皆改用 goccy/go-json，效能提升21-63%，記憶體配置減少40-60%
- 最佳化 BoundMethod 結構配置並快取 isVariadic 旗標，以降低每次呼叫的額外負擔
- 對具有`<=8`個引數的方法使用堆疊配置的引數緩衝區，以避免堆積配置
- 最佳化方法呼叫中的結果收集，避免為單一傳回值配置切片
- MIME 類型快取改用 sync.Map，以提升並行效能
- 讀取 HTTP 傳輸請求本文時使用緩衝區集區
- 在內容類型偵測器中延遲配置 CloseNotify 通道，以減少每個請求的配置
- 從資產伺服器移除 CSS 偵錯記錄
- 擴充 MIME 類型副檔名對應表，以涵蓋50多種常見 Web 格式（字型、音訊、視訊等）
- 更新 Window `X/Y`選項的文件 @ruhuang2001
- 更新`Frontend Runtime`文件，加入更多產生前端繫結的選項，由 @ndianabasi 貢獻
- 更新 Wails v3 Asset Server 的文件頁面，由 @ndianabasi 貢獻
- **重大變更**：移除套件層級的對話方塊函式（`application.InfoDialog()`、`application.QuestionDialog()`等）。請改用`app.Dialog`管理器：`app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()`
- 更新對話方塊文件以符合實際 API：使用`app.Dialog.*`、搭配回呼使用的`AddButton()`（而非`SetButtons()`）、`SetDefaultButton(*Button)`（而非字串）、`AddFilter()`（而非`SetFilters()`）、`SetFilename()`（而非`SetDefaultFilename()`），並使用`app.Dialog.OpenFile().CanChooseDirectories(true)`選取資料夾
- **不向後相容的變更**：現在預設建立生產版本。若要建立開發版本，請在 Taskfile 中設定`DEV=true`。請產生新專案以查看範例。此變更由 @leaanthony 完成
- 發出具有零個或一個資料引數的自訂事件時，資料值會直接指派給 Data 欄位，不再包裝於切片中；由[@fbbdev](https://github.com/fbbdev)在[#4633](https://github.com/wailsapp/wails/pull/4633)中實作
- Windows 系統匣現在會透過切換`NIS_HIDDEN`來遵循`SystemTray.Show()`/`Hide()`，讓應用程式能真正消失並重新顯示（#4653）。
- 系統匣註冊會重複使用已解析的圖示、只設定一次`NOTIFYICON_VERSION_4`，並啟用`NIF_SHOWTIP`，使工具提示能在 Explorer 重新啟動後恢復（#4653）。
- macOS：視窗置中時使用`visibleFrame`而非`frame`，以排除選單列和 Dock 區域
- macOS：視窗置中時使用`visibleFrame`而非`frame`，以排除選單列和 Dock 區域
- 使用`-config`參數執行`wails3 update build-assets`時，透過`-product*`參數設定的值會
- `window.NativeWindowHandle()` -> `window.NativeWindow()`，由 @leaanthony 在[#4471](https://github.com/wailsapp/wails/pull/4471)中實作
- 重構內部視窗處理，由 @leaanthony 在[#4471](https://github.com/wailsapp/wails/pull/4471)中實作
- 移除`application.WindowIDKey`和`application.WindowNameKey`（由`application.WindowKey`取代），由[@leaanthony](https://github.com/leaanthony)實作
- ContextMenuData 現在傳回字串而非 any，由[@leaanthony](https://github.com/leaanthony)實作
- 在 JS/TS 繫結中，固定長度陣列類型的類別欄位現在會以預期長度初始化，而非初始化為空陣列；由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中實作
- ContextMenuData 現在傳回字串而非 any，由[@leaanthony](https://github.com/leaanthony)實作
- `application.NewService`不再接受選項作為可選參數（請改用`application.NewServiceWithOptions`），由[@leaanthony](https://github.com/leaanthony)在[#4024](https://github.com/wailsapp/wails/pull/4024)中實作
- 移除`nanoid`相依性，由[@leaanthony](https://github.com/leaanthony)實作
- 更新 Window 範例以示範 mica/acrylic/tabbed 視窗樣式，由[@leaanthony](https://github.com/leaanthony)實作
- 在 JS/TS 繫結中，`internal.js/ts`模型檔案已移除；現在所有模型都可在`models.js/ts`中找到。由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中實作
- 在 JS/TS 繫結中，具名類型絕不會呈現為其他具名類型的別名；舊行為現在僅限於別名。由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中實作
- 在 JS/TS 繫結的類別模式中，類型為型別參數的結構欄位會標記為可選，且絕不會自動初始化。由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中實作
- 從範本移除 ESLint，由[@IanVS](https://github.com/IanVS)在[#4059](https://github.com/wailsapp/wails/pull/4059)中實作
- 將著作權年份更新為2025，由[@IanVS](https://github.com/IanVS)在[#4037](https://github.com/wailsapp/wails/pull/4037)中實作
- 新增 event.Sender 文件，由[@IanVS](https://github.com/IanVS)在[#4075](https://github.com/wailsapp/wails/pull/4075)中實作
- 支援 Go 1.24，由[@leaanthony](https://github.com/leaanthony)實作
- `ServiceStartup`掛鉤現在會在呼叫`App.Run`時叫用，而非在`application.New`中叫用；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中實作
- `ServiceStartup`錯誤現在會從`App.Run`傳回，而非終止程序；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中實作
- 從 JS 發出的繫結與對話方塊呼叫，現在會以錯誤物件而非字串拒絕；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中實作
- 改善 Windows 上的系統匣選單定位，由[@leaanthony](https://github.com/leaanthony)實作
- JS 執行階段已移植至 TypeScript，由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中實作
- runtime 在匯入後會立即初始化，無須等待視窗載入，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 完成
- runtime 不再匯出 init 方法。可使用僅為產生副作用的匯入來初始化，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 完成
- 繫結方法現在會傳回 `CancellablePromise`；若取消，則會以 `CancelError` 拒絕。呼叫的實際結果會被捨棄，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 完成
- 內建服務型別現在統一稱為 `Service`，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 完成
- 帶有選項的內建服務建立函式現在統一稱為 `NewWithConfig`，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 完成
- `sqlite` 服務的 `Select` 方法現在命名為 `Query`，以便與 Go API 保持一致，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 完成
- 範本：將 runtime 移至「dependencies」，並整理 package.json 檔案，由 [@IanVS](https://github.com/IanVS) 於 [#4133](https://github.com/wailsapp/wails/pull/4133) 完成
- 在開發環境中建立應用程式套件並臨時簽署，以啟用特定 macOS API，由 [@popaprozac](https://github.com/popaprozac) 於 [#4171](https://github.com/wailsapp/wails/pull/4171) 完成
- 將建置資產移至各平台專用目錄，由 [@leaanthony](https://github.com/leaanthony) 完成
- 將 Taskfile 移至各平台專用目錄並重新命名，由 [@leaanthony](https://github.com/leaanthony) 完成
- 大幅改善缺少 `index.html` 時的使用體驗，由 [@leaanthony](https://github.com/leaanthony) 完成
- [Windows] 改善視窗最小化與還原的效能，由 [@leaanthony](https://github.com/leaanthony) 完成。以 [562589540](https://github.com/562589540) 原先的 [PR](https://github.com/wailsapp/wails/pull/3955) 為基礎
- 移除 `ShouldClose` 選項（請改為註冊 events.Common.WindowClosing 的掛鉤），由 [@leaanthony](https://github.com/leaanthony) 完成
- [Windows] 減少開啟視窗時的閃爍，由 [@leaanthony](https://github.com/leaanthony) 完成
- 移除 `Window.Destroy`，因為此函式原本就僅供內部使用，由 [@leaanthony](https://github.com/leaanthony) 完成
- 將 `WindowClose` 事件重新命名為 `WindowClosing`，由 [@leaanthony](https://github.com/leaanthony) 完成
- 前端建置現在會依建置類型使用 Vite 環境「development」或「production」，由 [@leaanthony](https://github.com/leaanthony) 完成
- 更新至 go-webview2 v1.19，由 [@leaanthony](https://github.com/leaanthony) 完成
- 確保使用 Taskfile 的分支版本，由 @leaanthony 完成
- 更新 Taskfile 的分支版本，以修正使用以下方式安裝時的版本問題
- 使用 Taskfile 的分支版本，以修正使用以下方式安裝時的版本問題
- `service.OnStartup` 現在會在發生錯誤時關閉應用程式，並執行
- 重構系統匣點擊訊息，使其更符合使用者互動方式，由
- 資產嵌入納入 `all:frontend/dist`，以支援會產生以下內容的框架
- Taskfile 重構，由 [leaanthony](https://github.com/leaanthony) 於
- 升級至 `go-webview2` v1.0.16，由
- 修正 `Screen` 型別，使其包含 `ID` 而非 `Id`，由
- 更新 `go.mod.tmpl` 的 Wails 版本以支援 `application.ServiceOptions`，由
- 修正服務名稱的判定方式，由 [windom](https://github.com/windom/) 於
- mkdocs serve 現在使用 Docker，由 [leaanthony](https://github.com/leaanthony) 完成
- 將開發設定整合至 `config.yml`，由
- 系統匣對話方塊現在預設使用應用程式圖示（若有；Windows），由
- 改善 macOS 的 GPU 與記憶體資訊回報，由
- 移除 `WebviewGpuIsDisabled` 和 `EnableFraudulentWebsiteWarnings`
- 事件 API 變更：`On`/`Emit` -> 使用者事件，`OnApplicationEvent` ->
- 修正 Linux 上的事件 API，由 [TheGB0077](https://github.com/TheGB0077) 於
- [CI] 改善 actions，並允許在分支儲存庫中也能執行 actions，以及
- 將 `AbsolutePosition()` 重新命名為 `Position()`，由
- 將 Linux WebKit 相依套件從 webkitgtk2-4.0 更新為 webkit2gtk-4.1，以
- 隨附的 JS runtime 指令碼現在是 ESM 模組：匯入該指令碼的 script 標籤
- `@wailsio/runtime` 套件不會在 `window.wails` 上公開其 API
- 視窗 API 模組 `@wailsio/runtime/src/window` 現在會公開其所屬的
- JS 視窗 API 已更新，以符合目前的 Go `WebviewWindow`
- 繫結產生器現在預設使用依 ID 呼叫。`-id` CLI 選項
- 新的繫結程式碼配置：輸出檔案先前依資料夾整理
- 結構欄位 `application.Options.Bind` 已重新命名為
- 繫結服務的新語法：服務執行個體現在必須包裝在
- 在非終端機或 CI 環境中停用旋轉指示器，由

### 已移除

- **破壞性變更**：從各視窗的 `WindowsWindow` 選項中移除 `EnabledFeatures`、`DisabledFeatures` 和 `AdditionalLaunchArgs`。請改用應用程式層級的 `Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures` 和 `Options.Windows.AdditionalBrowserArgs`。這些旗標會全域套用至共用的 WebView2 環境（#4559），由 @leaanthony 完成
- 移除 Windows 上原生的 `IDropTarget` 實作，改用以 JavaScript 為基礎的方法（與 v2 的行為一致）
- 移除 github.com/wailsapp/mimetype 相依套件，改用擴充的副檔名對應表與標準函式庫的 http.DetectContentType，使二進位檔大小減少約 1.2MB
- 透過為 Linux 檔案總管實作精簡的 .desktop 檔案剖析器，移除 gopkg.in/ini.v1 相依套件，節省約 45KB
- 改用 Go 1.21+ 標準函式庫的 slices 套件和最精簡的內部輔助函式，從執行階段程式碼中移除 samber/lo，節省約 310KB
- 從 Darwin URL scheme 處理常式中移除偵錯 printf 陳述式（#4834）
- **破壞性變更**：移除 `linux:WindowLoadChanged` 事件；請改用 `linux:WindowLoadFinished` 偵測 WebView 何時完成載入（#3896），由 @leaanthony 完成

### 破壞性變更

- **Manager API 重構**：將應用程式 API 從扁平結構重新組織為分類清楚的管理器，以改善程式碼組織及可探索性；由 [@leaanthony](https://github.com/leaanthony) 於 [#4359](https://github.com/wailsapp/wails/pull/4359) 完成
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- 重新命名 Service 方法：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`；由 [@leaanthony](https://github.com/leaanthony) 完成
- 將 `Path` 和 `Paths` 方法移至 `application` 套件；由 [@leaanthony](https://github.com/leaanthony) 完成
- 應用程式選單現在僅適用於 macOS；由 [@leaanthony](https://github.com/leaanthony) 完成

## v3.0.0-alpha.78 - 2026-04-21

## 新增

## 修正

## v3.0.0-alpha.77 - 2026-04-18

## 修正

## v3.0.0-alpha.76 - 2026-04-17

## 修正

## v3.0.0-alpha.75 - 2026-04-16

## 修正

## v3.0.0-alpha.74 - 2026-03-01

## 新增

## 修正

## v3.0.0-alpha.73 - 2026-02-27

## 修正

## v3.0.0-alpha.72 - 2026-02-16

## 修正

## v3.0.0-alpha.71 - 2026-02-10

## 新增

## 修正

## v3.0.0-alpha.70 - 2026-02-09

## 新增

## 修正

## v3.0.0-alpha.69 - 2026-02-08

## 新增

## 修正

## v3.0.0-alpha.68 - 2026-02-07

## 新增

## 變更

## 修正

## v3.0.0-alpha.67 - 2026-02-04

## 新增

## 變更

## 修正

## v3.0.0-alpha.66 - 2026-02-03

## 新增

## 變更

## 修正

## 移除

## v3.0.0-alpha.65 - 2026-02-01

## 新增

## v3.0.0-alpha.64 - 2026-01-26

## 新增

## v3.0.0-alpha.63 - 2026-01-25

## 修正

## v3.0.0-alpha.62 - 2026-01-22

## 修正

## v3.0.0-alpha.61 - 2026-01-20

## 修正

## v3.0.0-alpha.60 - 2026-01-14

## 修正

## v3.0.0-alpha.59 - 2026-01-11

## 變更

## v3.0.0-alpha.58 - 2026-01-09

## 修正

## v3.0.0-alpha.57 - 2026-01-05

## 變更

## 修正

## v3.0.0-alpha.56 - 2026-01-04

## 新增

## 變更

## 修正

## 移除

## v3.0.0-alpha.55 - 2026-01-02

## 變更

## 修正

## 移除

## v3.0.0-alpha.54 - 2025-12-29

## 新增

## 修正

## 移除

## v3.0.0-alpha.53 - 2025-12-27

## 新增

## 修正

## v3.0.0-alpha.52 - 2025-12-26

## 修正

## v3.0.0-alpha.51 - 2025-12-23

## 修正

## v3.0.0-alpha.50 - 2025-12-21

## 變更

## v3.0.0-alpha.49 - 2025-12-18

## 變更

## v3.0.0-alpha.48 - 2025-12-16

## 新增

## 變更

## 修正

## v3.0.0-alpha.47 - 2025-12-15

## 新增

## 修正

## v3.0.0-alpha.46 - 2025-12-14

## 新增

## 移除

## v3.0.0-alpha.45 - 2025-12-13

## 新增

## 修正

## v3.0.0-alpha.44 - 2025-12-12

## 新增

## 變更

## 修正

## v3.0.0-alpha.43 - 2025-12-11

## 新增

## v3.0.0-alpha.42 - 2025-12-10

## 新增

## v3.0.0-alpha.41 - 2025-11-23

## 修正

## v3.0.0-alpha.40 - 2025-11-13

## 修正

## v3.0.0-alpha.39 - 2025-11-12

## 新增

## 變更

## v3.0.0-alpha.38 - 2025-11-04

## 新增

## 變更

## 修正

## v3.0.0-alpha.37 - 2025-11-02

## 修正

## v3.0.0-alpha.36 - 2025-10-15

## 修正

## v3.0.0-alpha.35 - 2025-10-14

## 修正

## v3.0.0-alpha.34 - 2025-10-06

## 新增

## 修正

## v3.0.0-alpha.33 - 2025-10-04

## 已修正

## v3.0.0-alpha.32 - 2025-10-02

## 已修正

## v3.0.0-alpha.31 - 2025-09-27

## 已修正

## v3.0.0-alpha.30 - 2025-09-26

## 已修正

## v3.0.0-alpha.29 - 2025-09-25

## 已新增

## 已變更

## 已修正

## v3.0.0-alpha.29 - 2025-09-25

## 已新增

## 已變更

## 已修正

## v3.0.0-alpha.27 - 2025-09-07

## 已修正

## v3.0.0-alpha.26 - 2025-08-24

## 已新增

## v3.0.0-alpha.25 - 2025-08-16

## 已變更

不再被忽略，且會覆寫設定值。

## v3.0.0-alpha.24 - 2025-08-13

## 已新增

## v3.0.0-alpha.23 - 2025-08-11

## 已修正

## v3.0.0-alpha.22 - 2025-08-10

## 已新增

## 已變更

+ 修正範圍過廣的 Linux 套件相依性，並修正過時的 RPM 相依性。

## v3.0.0-alpha.21 - 2025-08-07

## 已修正

## v3.0.0-alpha.20 - 2025-08-06

## 已修正

## v3.0.0-alpha.19 - 2025-08-05

## 已新增

## 已修正

## v3.0.0-alpha.18 - 2025-08-03

## 已新增

## 已修正

## v3.0.0-alpha.17 - 2025-07-31

## 已修正

## v3.0.0-alpha.16 - 2025-07-25

## 已新增

## v3.0.0-alpha.15 - 2025-07-25

## 已新增

## v3.0.0-alpha.14 - 2025-07-25

## 已新增

## v3.0.0-alpha.12 - 2025-07-15

### 已新增

### 已修正

## v3.0.0-alpha.11 - 2025-07-12

## 已新增

## v3.0.0-alpha.10 - 2025-07-06

### 破壞性變更

### 已新增

### 已修正

### 已變更

## v3.0.0-alpha.9 - 2025-01-13

### 已新增

### 已修正

### 已變更

## v3.0.0-alpha.8.3 - 2024-12-07

### 已變更

## v3.0.0-alpha.8.2 - 2024-12-07

### 已變更

`go install`，由 @leaanthony 完成

## v3.0.0-alpha.8.1 - 2024-12-07

### 已變更

`go install`，由 @leaanthony 完成

## v3.0.0-alpha.8 - 2024-12-06

### 已新增

@atterpac 於[#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) 於   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) 於   [#3867](https://github.com/wailsapp/wails/pull/3867)   由[atterpac](https://github.com/atterpac)於   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) 於   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   由[leaanthony](https://github.com/leaanthony)於指定的 X/Y 位置定位，見   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman)與   [leaanthony](https://github.com/leaanthony) 於   [#3823](https://github.com/wailsapp/wails/pull/3823)   由[leaanthony](https://github.com/leaanthony)於   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) 於   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony)。 -

### 變更

`service.OnShutdown`適用於先前已啟動的任何服務，由 @atterpac 於下列項目中變更：   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac 於[#3907](https://github.com/wailsapp/wails/pull/3907)   子資料夾，由 @atterpac 於下列項目中變更：   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) 於   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) 於   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   （已由`EnabledFeatures`與`DisabledFeatures`選項取代），由   [leaanthony](https://github.com/leaanthony)

### 修正

由 @michael-freling 修正 channel 變數，見   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) 於   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   於[#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) 於   [#3841](https://github.com/wailsapp/wails/pull/3841)   由[johnmccabe](https://github.com/johnmccabe)修正 Darwin 上編輯選單中的`PasteAndMatchStyle`角色，見   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) 於   [#3854](https://github.com/wailsapp/wails/pull/3854)，由   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   兩者不同。由 @nickisworking 於   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### 新增

[mmghv](https://github.com/mmghv) 於   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac)與   [leaanthony](https://github.com/leaanthony) 於   [#3570](https://github.com/wailsapp/wails/pull/3570)

### 變更

將應用程式事件`OnWindowEvent`改歸為視窗事件，由   [leaanthony](https://github.com/leaanthony)完成，見   [#3734](https://github.com/wailsapp/wails/pull/3734)   由[stendler](https://github.com/stendler)處理以`v3/`或`v3-`為前綴的分支，見   [#3747](https://github.com/wailsapp/wails/pull/3747)

### 修正

[etesam913](https://github.com/etesam913) 於   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) 於   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) 於   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) 於   [#3614](https://github.com/wailsapp/wails/pull/3614)，由   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720)，由   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693)，由   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720)，由   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693)，由   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746)，由   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 修正

## v3.0.0-alpha.5 - 2024-07-30

### 新增

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   由 @5aaee9 處理圖示點擊，見[#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) 於   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   於[#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) 於   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) 於   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   以 @Mai-Lapyst 的[PR](https://github.com/wailsapp/wails/pull/2044)為基礎   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   由[@fbbdev](https://github.com/fbbdev)處理 npm 套件，見   [#3334](https://github.com/wailsapp/wails/pull/3334)   於[#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT`，由[@abichinger](https://github.com/abichinger)於   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) 於   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) 於   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) 於   [#3674](https://github.com/wailsapp/wails/pull/3674)

### 修正

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) 於   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) 於   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) 於   [#3477](https://github.com/wailsapp/wails/pull/3477)   由 [Atterpac](https://github.com/atterpac) 於   [#3320](https://github.com/wailsapp/wails/pull/3320)。   於 [#3306](https://github.com/wailsapp/wails/pull/3306)。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) 於   [#2750](https://github.com/wailsapp/wails/pull/2750)。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   [jaybeecave](https://github.com/jaybeecave) 於   [#3052](https://github.com/wailsapp/wails/pull/3052)。   [@pylotlight](https://github.com/pylotlight) 於   [PR](https://github.com/wailsapp/wails/pull/3039)   安裝。由 [@pylotlight](https://github.com/pylotlight) 於   [PR](https://github.com/wailsapp/wails/pull/3032) 新增   [PR](https://github.com/wailsapp/wails/pull/3145)   空格 — @leaanthony。   [thomas-senechal](https://github.com/thomas-senechal) 於 PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) 於 PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   附加的視窗 [tw1nk](https://github.com/tw1nk) 於 PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) 於   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` 存在。   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   監聽器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) 於   [#3330](https://github.com/wailsapp/wails/pull/3330)   產生器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3334](https://github.com/wailsapp/wails/pull/3334)   產生器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) 於   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) 於   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) 於   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) 於   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) 於   [#3431](https://github.com/wailsapp/wails/pull/3431)   由 [@pekim](https://github.com/pekim) 於   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel)。   PR [#3466](https://github.com/wailsapp/wails/pull/3622)，由   [@5aaee9](https://github.com/5aaee9) 提交。   [@windom](https://github.com/windom/) 於   [#3636](https://github.com/wailsapp/wails/pull/3636)。   Windows，由 [@bruxaodev](https://github.com/bruxaodev/) 於   [#3691](https://github.com/wailsapp/wails/pull/3691)。

### 變更

[mmghv](https://github.com/mmghv) 於   [#3611](https://github.com/wailsapp/wails/pull/3611)   支援 Ubuntu 24.04 LTS，由 [atterpac](https://github.com/atterpac) 於   [#3461](https://github.com/wailsapp/wails/pull/3461)   必須具有 `type="module"` 屬性。由   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   物件，且不會啟動 WML 系統。此變更旨在改善   封裝。如有需要，可呼叫新的 `WML.Enable` 方法手動啟動 WML 系統。   隨附的 JS 執行階段指令碼仍會自動執行這兩項操作。由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   將 window 物件作為預設匯出。現在已無法再透過 ESM 具名匯入或命名空間匯入語法   匯入個別方法。   API。部分方法已變更名稱或原型，具體而言：`Screen`   變為 `GetScreen`；`GetZoomLevel`/`SetZoomLevel` 變為 `GetZoom`/`SetZoom`；   `GetZoom`、`Width` 和 `Height` 現在會直接傳回值，不再將其包裝於   物件中。由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   已移除。使用 `-names` CLI 選項可切換回按名稱呼叫。   由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   過去以其所在套件命名；現在則使用完整的 Go 匯入路徑，   包括模組路徑。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   對 `application.NewService` 的呼叫。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) 於   [#3574](https://github.com/wailsapp/wails/pull/3574)
