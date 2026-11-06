---
title: "變更日誌"
description: "Wails v3的版本歷程與發行說明"
slug: "changelog"
sourcePath: "changelog.md"
---

圖例：

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- 此專案的所有重要變更都會記錄在此檔案中。

此格式以[維護變更日誌](https://keepachangelog.com/en/1.0.0/)為基礎，而本專案遵循[語意化版本](https://semver.org/spec/v2.0.0.html)規範。

- `Added`表示新增功能。
- `Changed`表示現有功能的變更。
- `Deprecated`表示即將移除的功能。
- `Removed`表示已移除的功能。
- `Fixed`表示任何錯誤修正。
- `Security`表示安全漏洞。

_/

/_   * 自動更新 *   發行說明會從已合併的 PR 產生。請勿在一般貢獻 PR 中編輯此檔案，也不要重新建立已移除的`v3/UNRELEASED_CHANGELOG.md`。如需維護變更日誌，請參閱`AGENTS.md`。 _/

## [尚未發行]

## v3.2.0-nightly.20261103 - 2026-11-03

## 變更
- [v3] 在 Menu.Update 時重新整理系統匣選單，詳見 @johnmaguire 的[PR](https://github.com/wailsapp/wails/pull/6314)
- [v3/Windows] 為視窗實作 SetIcon，詳見 @hypherionmc 的[PR](https://github.com/wailsapp/wails/pull/6502)

## v3.2.0-nightly.20261102 - 2026-11-02

## 變更
- 以 DisableOverscroll 取代 DisableBounce，詳見 @lethusang 的[PR](https://github.com/wailsapp/wails/pull/6077)

## v3.2.0-nightly.20261101 - 2026-11-01

## 新增
- 新增實驗性的無 CGO macOS 與 Linux 後端，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6521)
- 原生 YAML 建置、可攜式封裝與簽署，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6531)

## 修正
- 修正各種v3夜間建置問題並清理相依套件，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6519)
- 修復v3發行驗證與測試資料分類，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6533)
- 解決跨平台發行驗證失敗，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6538)

## v3.1.0 - 2026-11-01

## 新增
- 新增實驗性的無 CGO macOS 與 Linux 後端，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6521)
- 原生 YAML 建置、可攜式封裝與簽署，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6531)

## 變更
- 避免複製已捨棄的桌面中繼資料，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6391)
- 不使用 fmt 格式化內容長度，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6393)
- 不收集候選項目即可選取索引根目錄，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6404)
- 重複使用不可搜尋檔案的複製緩衝區，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6408)
- 避免配置用於探測前綴的讀取器，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6410)
- 對已知副檔名略過 MIME 探測讀取，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6414)
- 防止佇列增長與縮減反覆震盪，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6415)
- 避免配置隨後會被捨棄的請求內容，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6419)
- 略過已停用之偵錯日誌的準備工作，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6420)
- 延遲配置 MIME 探測緩衝區，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6428)
- 減少超連結格式化配置，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6443)
- 移除項目符號格式化的暫時配置，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6449)
- 將 UUID 字串轉換減至一次配置，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6453)
- 減少陳述式控制代碼的格式化配置，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6454)
- 消除以名稱為基礎之 UUID 摘要的暫時配置，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6455)
- 避免複製尚未完成的命令輸出行，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6460)
- 在套件暫存期間以串流方式複製可執行檔，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6472)
- 避免重新掃描尚未完成的輸出行，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6473)
- 重複使用已編譯的命令路由模式，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6475)
- 避免在執行階段探索期間查詢中繼資料，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6482)
- 減少本地化備援頁面的配置，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6485)
- 避免為每個處理器重新編譯模式，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6490)

## 修正
- 為系統匣實作 SetTooltip，詳見 @fan711 的[PR](https://github.com/wailsapp/wails/pull/6017)
- 從設定入口網站讀取深色模式，詳見 @TechHutTV 的[PR](https://github.com/wailsapp/wails/pull/6072)
- 允許 Window.SetZoom 低於1.0，詳見 @APshenkin 的[PR](https://github.com/wailsapp/wails/pull/6166)
- 依發出順序傳遞至每個 Go 接聽器，詳見 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/6275)
- 在 Emit 中於 hookLock 鎖定期間讀取事件掛鉤，詳見 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/6276)
- 一律公告系統匣的 Menu 屬性，確保在 Run 之後建立的系統匣不會被捨棄，詳見 @AlbinoGeek 的[PR](https://github.com/wailsapp/wails/pull/6284)
- 將 Run 之前傳送的主執行緒工作保留至 Run 啟動，詳見 @AlbinoGeek 的[PR](https://github.com/wailsapp/wails/pull/6287)
- 在 GitHub 更新程式提供者中使用正確的資產 URL，詳見 @gavidroselj 的[PR](https://github.com/wailsapp/wails/pull/6304)
- 當 WebView2 省略 Content-Length 時保留請求本文，詳見 @strraaa 的[PR](https://github.com/wailsapp/wails/pull/6326)
- 以 DIP 空間回報 macOS 與 GTK 螢幕（混合 DPI 邊界），詳見 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/6346)
- 不要將已關閉的通知回報為點擊，詳見 @bercedev 的[PR](https://github.com/wailsapp/wails/pull/6367)
- 終端機關閉時停止開發程序，詳見 @xmasdev 的[PR](https://github.com/wailsapp/wails/pull/6371)
- 在證據檢查中，將變更日誌執行結果對應至 squash merge 提交，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6373)
- 將工作列重新啟動的防彈跳延遲維持一秒，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6386)
- 保留從全螢幕狀態開始的視窗大小限制，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6387)
- 不要調整全螢幕或最大化視窗的大小以符合新限制，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6388)
- 不要在 liquid-glass 範例中的 Run 之前開啟對話方塊，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6389)
- 在 Run 之前要求開啟對話方塊時不執行任何動作，而非永遠等待，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6390)
- 桌面項目沒有命令時改用後援機制，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6392)
- 在不保留檔名的情況下偵測目前內容，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6394)
- 辨識大小寫混合的資產副檔名，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6395)
- 釋放已耗用的內容偵測器緩衝區，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6396)
- 釋放已移除的 ringqueue 參照，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6397)
- 在選擇後援方案之前遵循標頭刪除設定，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6398)
- 在清送之前提交後援回應標頭，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6399)
- 僅尋找名稱完全相符的索引檔案，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6400)
- 對不可定位的 HEAD 回應避免產生主體及完整讀取，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6401)
- 跨越多次短讀取填滿 MIME 偵測前綴，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6402)
- 提交成功但為空的不可定位回應，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6403)
- 避免重複已展開的 URI 目標，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6405)
- 在欄位展開期間保留路徑，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6406)
- 遵循 XDG 桌面檔案搜尋目錄，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6407)
- 空清送時提交偵測器標頭，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6409)
- 略過遮蔽桌面項目的目錄，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6411)
- 保留在可執行檔路徑內找到的目標，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6412)
- 計時器觸發後釋放回呼擷取的內容，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6413)
- 避免 ldconfig 程式庫誤比對，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6416)
- 限制檔案系統程式庫的後援比對，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6417)
- 消除快取失效處理中的資料競爭，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6418)
- 遵循選用的記錄器設定，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6421)
- 避免共用啟動 URL 發生突變及使用過時的連接埠，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6422)
- 取消要求時停止開發代理重試，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6423)
- 缺少內嵌資產時傳回404，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6424)
- 處理內嵌執行階段的 HEAD 要求，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6425)
- 路由服務時保留逸出的路徑，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6426)
- 正確排序超大數值的預發行識別碼，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6427)
- 儲存工具下載項目之前拒絕 HTTP 錯誤，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6429)
- 在搜尋桌面項目期間遵循啟動期限，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6430)
- 在開發代理重試期間保留要求主體，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6432)
- 在每一條傳回路徑上關閉 HTTP 回應主體，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6434)
- 在平行搜尋程式庫時保留取消錯誤，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6435)
- 解碼前拒絕失敗的 GitHub 標籤回應，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6436)
- 保留含方括號的常值範本路徑，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6437)
- 僅安裝 wails3 一次並保留安裝程式錯誤，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6438)
- 忽略註解區塊中的隱藏項目和標頭，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6439)
- 呈現穩定版查詢錯誤，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6440)
- 列印純文字訊息而不解讀格式，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6441)
- 載入 JSON null 後避免寫入操作發生 panic，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6442)
- 讀取整個儲存區時傳回快照，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6444)
- 在部分資料列反覆運算後傳播錯誤，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6446)
- 將解壓縮目的地路徑保留為常值，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6447)
- 關閉及重新開啟時釋放預備陳述式控制代碼，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6448)
- 處理預備操作與關閉或重新開啟之間的競爭，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6450)
- 遵循後援索引頁面的 HEAD 要求，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6451)
- 要求使用完整的十六進位陳述式控制代碼，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6452)
- 在並行處理下保留預設的記憶體內資料庫，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6456)
- 將 stderr 與剖析後的 HEAD 輸出分開，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6457)
- 從內嵌套件組合產生執行階段資產，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6458)
- 偵測 TODO、FIXME 和 TBD 預留位置，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6459)
- 驗證版本資訊 API 回應，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6461)
- 在狀態檢查中遵循工作目錄，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6462)
- 防止複製同一檔案時遺失資料，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6463)
- 將工作目錄納入快取識別資訊，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6464)
- 每次檢查僅探測連接埠一次，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6465)
- 避免快取項目解碼為 null 時發生 panic，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6466)
- 遞增數值預發行版本時不受整數範圍限制，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6467)
- 防止舊版檢查掩蓋快取未命中，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6468) 中完成
- 在執行器的基底目錄中執行預設任務，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6469) 中完成
- 保留檔案總管命令中的常值百分比符號，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6470) 中完成
- 統一前置條件與任務命令的目錄，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6471) 中完成
- 以產生的一般檔案取代符號連結，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6474) 中完成
- 建置失敗時保留現有的 DMG，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6476) 中完成
- 保留加上引號的正規表示式及路徑中的反斜線，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6477) 中完成
- 配置影像記憶體前先驗證產生的圖示尺寸，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6478) 中完成
- 透過 shell 處理命令鏈與重新導向，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6479) 中完成
- 建置失敗時保留先前的 Linux 套件，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6480) 中完成
- 使用平台專屬的 shell 備援執行方式，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6481) 中完成
- 支援提供路徑開頭含重複斜線的資產，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6483) 中完成
- 識別含建置中繼資料的穩定版本，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6484) 中完成
- 正確回報資源產生失敗，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6486) 中完成
- 防止取代作用中回報器時發生競爭，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6487) 中完成
- 將並行的快取檔案儲存作業序列化，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6488) 中完成
- 支援多個範本相等性候選項目，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6489) 中完成
- 正規化已儲存但為 null 的類別對應表，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6491) 中完成
- 顯示失敗的變數 shell 命令所產生的診斷資訊，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6492) 中完成
- 不將 shell 命令鏈納入隱含 Go 快取，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6493) 中完成
- 將 Linux 類別持久化作業序列化，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6494) 中完成
- 快取含加上引號或逸出空格的 Go 輸出，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6495) 中完成
- 編碼失敗時保留 Mac 圖示，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6496) 中完成
- 不再捨棄無效的 Go 測試次數旗標，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6497) 中完成
- 複製至新檔案時保留來源權限，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6498) 中完成
- 修正各種 v3 每夜建置問題並清理相依性，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6519) 中完成
- 修復 v3 發行驗證與測試資料分類，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6533) 中完成
- 解決跨平台發行驗證失敗，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6538) 中完成

## v3.1.0-nightly.20261028 - 2026-10-28

## 已變更
- 重複使用已編譯的正規表示式模式，以加快變更記錄產生速度，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6490) 中完成
- 提供本地化備援頁面時減少記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6485) 中完成
- 避免不必要的中繼資料查詢，以加快 AppImage 執行階段探索，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6482) 中完成
- 重複使用已編譯的命令路由模式，以改善 wake 效能，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6475) 中完成
- 避免重新掃描尚未完成的輸出行，以改善效能，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6473) 中完成
- 在 MSIX 套件暫存期間以串流方式複製可執行檔，以減少記憶體用量並加快建置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6472) 中完成
- 避免複製尚未完成的命令輸出行，以改善效能，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6460) 中完成
- 產生以名稱為基礎的 UUID 時減少記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6455) 中完成
- 格式化 SQLite 陳述式控制代碼時減少記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6454) 中完成
- 將 UUID 字串轉換減至僅配置一次記憶體，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6453) 中完成
- 格式化變更記錄項目符號時減少記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6449) 中完成
- 格式化終端機超連結時減少記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6443) 中完成
- 提供資產時延後探測 MIME 緩衝區，以減少記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6428) 中完成
- 停用偵錯記錄時，略過資產伺服器中的偵錯記錄準備作業，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6420) 中完成
- 避免建立會被捨棄的要求內容，以減少 WebView 要求處理時的記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6419) 中完成
- 對已知副檔名略過 MIME 探測，以加快資產提供速度，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6414) 中完成
- 提供資產時避免使用讀取器探測內容，以減少記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6410) 中完成
- 提供不可搜尋的資產檔案時減少記憶體配置，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6408) 中完成
- 避免收集候選項目，以加快資產伺服器選取索引根目錄的速度，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6404) 中完成
- 不使用 fmt 格式化內容長度，以改善資產伺服器效能，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6393) 中完成
- 讀取桌面檔案總管中繼資料時減少記憶體複製，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6391) 中完成

## 已修正
- 將來源檔案複製至新的目的地檔案時保留其權限，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6498) 中完成
- 最小或最大尺寸限制變更時，不調整全螢幕及最大化視窗的大小，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6388) 中完成
- 發布 wailsapp/wails 儲存庫的排程星號歷史記錄更新，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6500) 中完成
- 不再於 wake 中捨棄無效的 Go 測試次數旗標，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6497) 中完成
- 編碼失敗時保留 Mac 圖示，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6496) 中完成
- 修正含加上引號或逸出空格的 Go 輸出造成建置快取未命中的問題，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6495) 中完成
- 將 Linux 上的通知類別持久化作業序列化，以防止並行寫入競爭，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6494) 中完成
- 不將 shell 命令鏈納入隱含 Go 建置快取，由 @taliesin-ai 於 [PR](https://github.com/wailsapp/wails/pull/6493) 中完成
- 修正儲存的資料為 null 時無法載入通知類別對應表的問題，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6491)
- 顯示變數 shell 命令失敗時的診斷資訊，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6492)
- 修正存在多個相等候選項時的 wake 範本比對，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6489)
- 將並行的快取檔案儲存作業序列化，以防止快取寫入損毀，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6488)
- 防止取代目前的 reporter 時發生競爭條件，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6487)
- 在 CLI 中準確回報資源產生失敗，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6486)
- 讓 CLI 能辨識包含建置中繼資料的穩定版本，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6484)
- 提供請求路徑含有重複前置斜線的資產，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6483)
- 修正 wake 在各平台上的 shell 後援執行，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6481)
- 建置失敗時保留先前的 Linux 套件，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6480)
- 透過 shell 執行命令鏈與重新導向，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6479)
- 在解析 wake 命令時保留加上引號的規則運算式及路徑中的反斜線，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6477)
- 封裝建置失敗時保留現有的 DMG 檔案，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6476)
- 在繫結輸出中以產生的一般檔案取代符號連結，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6474)
- 在 Linux 檔案總管命令中保留字面百分比符號，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6470)
- 對齊前置條件與任務命令的目錄，使其解析結果一致，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6471)
- 在執行器的基底目錄中執行預設任務，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6469)
- 防止舊版檢查掩蓋快取未命中，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6468)
- 修正 CLI 的預發行版本編號，使其遞增不受整數溢位限制，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6467)
- 避免快取項目解碼為 null 時發生 panic，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6466)
- 讓 CLI 每次檢查僅探測一次連接埠，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6465)
- 將任務目錄納入快取識別資訊，以避免使用過時的 wake 結果，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6464)
- 防止封裝作業將檔案複製到自身時遺失資料，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6463)
- 讓狀態檢查遵循任務目錄，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6462)
- 偵測變更記錄項目中的 TODO、FIXME 與 TBD 預留位置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6459)
- 讓 CLI 從嵌入式套件產生執行階段資產，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6458)
- 在 git 命令中將 stderr 與解析後的 HEAD 輸出分開，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6457)
- 在並行存取時維持預設記憶體內 SQLite 資料庫，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6456)
- 對後援索引頁面的 HEAD 請求傳回正確回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6451)
- 處理 SQLite 陳述式準備作業與資料庫關閉或重新開啟之間的競爭條件，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6450)
- SQLite 資料庫關閉或重新開啟時釋放已準備陳述式的控制代碼，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6448)
- 使用 gosod 時將解壓縮目的地路徑保留為字面值，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6447)
- 在部分資料列反覆運算後傳播 SQLite 查詢錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6446)
- 讓 kvstore 的完整儲存區讀取傳回快照，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6444)
- 防止 kvstore 載入 JSON null 值後在寫入時發生 panic，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6442)
- 輸出純文字終端訊息，不解讀格式化動詞，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6441)
- 讓更新程式在查詢最新穩定版本時顯示錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6440)
- 忽略變更記錄註解區塊中的隱藏項目與標頭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6439)
- 僅安裝一次 wails3，並讓更新程式顯示安裝程式錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6438)
- 在 gosod 中保留含方括號的字面範本路徑，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6437)
- 讓更新程式在解碼 GitHub 標籤回應前拒絕失敗的回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6436)
- 在平行搜尋程式庫期間保留取消錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6435)
- 在每條傳回路徑上關閉 HTTP 回應本文，以防止資源洩漏，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6434)
- 開發代理重試請求時保留請求本文，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6432)
- 在探索桌面檔案總管期間遵循啟動期限，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6430)
- 儲存 AppImage 工具下載項目之前拒絕 HTTP 錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6429)
- 在 semver 比較中正確排序過大的數字預發行識別碼，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6427)
- 在資產伺服器中路由服務時保留逸出的路徑，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6426)
- 在資產伺服器中處理對隨附執行階段的 HEAD 請求，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6425)
- 缺少隨附資產時傳回404，而非不正確的回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6424)
- 請求取消時停止開發代理重試，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6423)
- 避免資產伺服器共用的起始 URL 發生異動及使用過時的連接埠，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6422)
- 讓資產伺服器遵循選用的記錄器設定，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6421)
- 消除程式庫路徑快取失效處理中的資料競爭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6418)
- 避免搜尋 ldconfig 輸出時誤判程式庫相符，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6416)
- 防止資產伺服器佇列在負載下增長並反覆縮減，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6415)
- 在防彈跳計時器觸發後釋放回呼所擷取的值，以免保留過時的參照，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6413)
- 保留在可執行檔路徑中找到的檔案總管目標，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6412)
- 在檔案總管中略過會遮蔽桌面項目的目錄，詳見 [PR](https://github.com/wailsapp/wails/pull/6411)，作者為 @taliesin-ai
- 修正回應本文為空時資產伺服器未傳送標頭的問題，詳見 [PR](https://github.com/wailsapp/wails/pull/6409)，作者為 @taliesin-ai
- 尋找應用程式圖示時遵循 XDG 桌面檔案搜尋目錄，詳見 [PR](https://github.com/wailsapp/wails/pull/6407)，作者為 @taliesin-ai
- 在檔案總管中展開欄位時保留檔案路徑，詳見 [PR](https://github.com/wailsapp/wails/pull/6406)，作者為 @taliesin-ai
- 避免在檔案總管中重複展開後的 URI 目標，詳見 [PR](https://github.com/wailsapp/wails/pull/6405)，作者為 @taliesin-ai
- 提交成功且內容為空的不可搜尋資產回應，而不再讓其保持待處理狀態，詳見 [PR](https://github.com/wailsapp/wails/pull/6403)，作者為 @taliesin-ai
- 修正資產伺服器在多次短讀取時填入 MIME 偵測前置資料的方式，詳見 [PR](https://github.com/wailsapp/wails/pull/6402)，作者為 @taliesin-ai
- 避免資產伺服器為不可搜尋的 HEAD 請求讀取完整回應本文，詳見 [PR](https://github.com/wailsapp/wails/pull/6401)，作者為 @taliesin-ai
- 資產伺服器解析目錄請求時僅提供完全相符的索引檔案，詳見 [PR](https://github.com/wailsapp/wails/pull/6400)，作者為 @taliesin-ai
- 修正在不支援動作的伺服器上，已關閉的 Linux 通知被回報為點擊的問題，詳見 [PR](https://github.com/wailsapp/wails/pull/6367)，作者為 @bercedev
- 終端機關閉時停止開發程序，詳見 [PR](https://github.com/wailsapp/wails/pull/6371)，作者為 @xmasdev
- Windows 工作列重新啟動的防彈跳機制會等待一秒，再清除旗標，詳見 [PR](https://github.com/wailsapp/wails/pull/6386)，作者為 @taliesin-ai
- 修正資產伺服器遺失後備回應標頭的問題，詳見 [PR](https://github.com/wailsapp/wails/pull/6399)，作者為 @taliesin-ai
- 資產伺服器會在選擇後備回應前遵循標頭刪除操作，詳見 [PR](https://github.com/wailsapp/wails/pull/6398)，作者為 @taliesin-ai
- 釋放資產伺服器中已移除的環形佇列參照，以避免記憶體洩漏，詳見 [PR](https://github.com/wailsapp/wails/pull/6397)，作者為 @taliesin-ai
- 釋放資產伺服器中已使用的內容偵測緩衝區，以減少記憶體用量，詳見 [PR](https://github.com/wailsapp/wails/pull/6396)，作者為 @taliesin-ai
- 提供資產時可辨識大小寫混合的資產副檔名，詳見 [PR](https://github.com/wailsapp/wails/pull/6395)，作者為 @taliesin-ai
- 資產伺服器無須保留檔名即可正確偵測內容類型，詳見 [PR](https://github.com/wailsapp/wails/pull/6394)，作者為 @taliesin-ai
- 桌面項目沒有命令時改用預設值，詳見 [PR](https://github.com/wailsapp/wails/pull/6392)，作者為 @taliesin-ai
- 防止 liquid-glass 範例在非 macOS 上於應用程式啟動前顯示對話方塊，詳見 [PR](https://github.com/wailsapp/wails/pull/6389)，作者為 @taliesin-ai
- 在 Windows 上從受限制的尺寸進入全螢幕時保留視窗尺寸限制，詳見 [PR](https://github.com/wailsapp/wails/pull/6387)，作者為 @taliesin-ai
- 修正在 Run 前呼叫對話方塊會永久封鎖而非傳回的問題，詳見 [PR](https://github.com/wailsapp/wails/pull/6390)，作者為 @taliesin-ai

## 安全性
- 在配置影像前驗證產生的圖示尺寸，以防止過度使用記憶體，詳見 [PR](https://github.com/wailsapp/wails/pull/6478)，作者為 @taliesin-ai
- 使用版本資訊 API 回應前先加以驗證，詳見 [PR](https://github.com/wailsapp/wails/pull/6461)，作者為 @taliesin-ai
- SQLite 外掛程式要求完整的十六進位陳述式控制代碼，詳見 [PR](https://github.com/wailsapp/wails/pull/6452)，作者為 @taliesin-ai
- 限制檔案系統程式庫的後備比對，以防載入非預期的程式庫，詳見 [PR](https://github.com/wailsapp/wails/pull/6417)，作者為 @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## 變更
- 將 404 頁面連結至 v2.wails.io 上對應的 Wails v2 文件頁面，詳見 [PR](https://github.com/wailsapp/wails/pull/6366)，作者為 @taliesin-ai
- 將 v2 文件網站移至 v2.wails.io，詳見 [PR](https://github.com/wailsapp/wails/pull/6365)，作者為 @taliesin-ai
- 從 wails.io 提供 v3 文件，並從 v2.wails.io 提供 v2 文件，詳見 [PR](https://github.com/wailsapp/wails/pull/6364)，作者為 @taliesin-ai
- 更新 v3.0.0 版本資訊並更正安全公告中的主張，詳見 [PR](https://github.com/wailsapp/wails/pull/6363)，作者為 @taliesin-ai

## 修正
- 依發出順序循序傳遞 Go 事件監聽器回呼，詳見 [PR](https://github.com/wailsapp/wails/pull/6275)，作者為 @DevLumuz
- 一律公開系統匣的 Menu 屬性，避免在 Linux 上捨棄 Run 後建立的系統匣，詳見 [PR](https://github.com/wailsapp/wails/pull/6284)，作者為 @AlbinoGeek
- 將 Run 前提交的主執行緒工作保留至應用程式啟動，詳見 [PR](https://github.com/wailsapp/wails/pull/6287)，作者為 @AlbinoGeek
- 下載 GitHub 更新程式的版本資產時使用正確的資產 API URL，詳見 [PR](https://github.com/wailsapp/wails/pull/6304)，作者為 @gavidroselj
- 修正混合 DPI 的 macOS 與 Linux 設定中的螢幕邊界和位置，詳見 [PR](https://github.com/wailsapp/wails/pull/6346)，作者為 @julianstorer
- 在 Windows 上，當 WebView2 省略 Content-Length 標頭時保留請求本文，詳見 [PR](https://github.com/wailsapp/wails/pull/6326)，作者為 @strraaa
- 在鎖定期間擷取鉤子的快照，修正發出事件鉤子時的資料競爭，詳見 [PR](https://github.com/wailsapp/wails/pull/6276)，作者為 @DevLumuz
- 修正 `SystemTray.SetTooltip()` 在 Linux 上毫無作用的問題：現在會設定 StatusNotifierItem 的工具提示描述（#6017）
- 在每週發行工作流程中執行變更日誌證據檢查，詳見 [PR](https://github.com/wailsapp/wails/pull/6374)，作者為 @taliesin-ai
- 修正版本變更日誌證據檢查，將執行記錄對應至壓縮合併提交，詳見 [PR](https://github.com/wailsapp/wails/pull/6373)，作者為 @taliesin-ai
- 修正 Linux 深色模式偵測在深色桌面上仍傳回淺色，直到配色變更才更新的問題，以及 `linux:SystemThemeChanged` 在 GTK4 後端永遠不會觸發的問題。兩個後端現在共用同一個實作，並以 `org.freedesktop.appearance` 入口網站命名空間為基礎（#6072）
- 在有框線視窗（Windows、macOS、Linux）上，`Window.SetZoom` 不再將低於 1.0 的倍率提高至 1.0：WebView2 `ZoomFactor`、WKWebView `magnification` 與 WebKitGTK 的縮放層級均接受任何正值，而 `Zoom` 視窗選項原本就會直接傳遞這類值。無框線視窗仍維持 1.0 的下限，因為其拖曳與調整大小區域不會依縮放比例調整（[#4590](https://github.com/wailsapp/wails/issues/4590)）；非正值、NaN 與 Inf 倍率現在會被忽略，不再提高至 1.0 — 作者為 @APshenkin

## v3.0.0 - 2026-10-26

## 變更
- 每夜預發行版會以 `vX.Y.0-nightly.YYYYMMDD` 發布，用來預覽下一個每週發行版；每週發行版則會遞增次要版本（`vX.Y.0`），因此 Go 的 `@latest` 一律解析為每週發行版或緊急修正版；`@wailsio/runtime` npm 套件版本一律與 Wails 發行版本一致，每夜版使用 `nightly` dist-tag
- 專案範本會將 `@wailsio/runtime` 固定為產生該範本的 CLI 版本
- Wails v3 需要 Go 1.26 或更新版本（Go 1.26.10 或 1.27.3 及更新版本包含 `crypto/x509` 安全性修正）

## 修正
- 修正 `wails3 update cli` 找不到 v3 發行版的問題：它只考慮 v2 標籤，且從錯誤的路徑安裝
- 修正部分 Windows 組建版本在繪製應用程式選單列時，Windows 應用程式因 `panic: SetTextColor failed` 而當機的問題：選單列現在會在有效的裝置內容上繪製，否則改用預設繪製方式；`UAHMENUITEM.Position` 現在與 C 的 `int` 相符，且選單繪製絕不會引發 panic（#6311）
- 在等待重新啟動期間確保暫存更新的安全：關機過程中啟動的下載或第二次重新啟動，不會刪除交換輔助程式正等待安裝的酬載，也不會與其爭用；由 @julianstorer 貢獻
- 修正 Windows 應用程式在 `applyDarkMode` 中當機的問題：`pkg/w32` 初始化時卸載了 `uxtheme.dll`，但其序號函式指標仍在使用中；該模組現在會在程序的整個生命週期內保持載入（#6356）
- 在交換前驗證暫存的酬載，並在捨棄備份前驗證已安裝的應用程式，使損壞的更新能夠回復，而不會取代正常運作的應用程式；由 @julianstorer 貢獻

## 已移除
- 在 [PR](https://github.com/wailsapp/wails/pull/5736) 完成最後一次翻譯同步後移除 Crowdin 整合；由 @github-actions[bot] 貢獻

## 安全性
- 透過在所有框架中注入防護指令碼，封鎖 Linux（GTK3 和 GTK4）上來自跨來源子框架的保留 `wails:` 訊息，以緩解 GHSA-2p5f-2gx8-mh8m
- 將 `github.com/klauspost/compress` 更新至 1.18.7（GHSA-259r-337f-4rfw），將 `golang.org/x/image` 更新至 0.45.0，並使用 Go 1.26.10 建置，以納入 `crypto/x509` 修正（GO-2026-6809）
- 在執行階段套件、設定精靈、範例及專案範本中，將 `vite` 更新至 8.2.4 或更新版本（GHSA-9jrq-w75r-8gcw、GHSA-rq7h-c2jc-7f22、GHSA-vfpm-58rq-9qcg、GHSA-r825-w8mm-5pxw、GHSA-587h-7wmj-h5pm、GHSA-698v-4r3x-3vcv）

### 已知限制
- **Linux**：子框架防護機制無法阻止*同源*子框架（例如`srcdoc`或`about:blank` iframe）透過其父框架存取原生訊息處理常式，因為 WebKitGTK 不會向該處理常式提供框架資訊。跨來源框架以及不含`allow-same-origin`的沙箱框架會遭到封鎖。請勿在 Wails 視窗中嵌入不受信任的同源內容。
- **Android（實驗性）：**`window.wails` JavaScript 橋接器會向 WebView 中的每個框架公開，且不會檢查是否為主框架或檢查來源。嵌入 iframe 中的內容（包括跨來源內容）可以呼叫繫結方法，並傳送保留的 `wails:` 訊息。請勿在應用程式的 WebView 中嵌入不受信任的第三方內容。這與 macOS 和 iOS 不同，後兩者會忽略來自子框架的 `wails:` 訊息。預計在 3.0.1 中修正此問題。
- <strong>所有平台：</strong>來自子框架的繫結方法呼叫會透過 HTTP 傳送至 `/wails/runtime`，而不是原生訊息處理常式，因此防護不會限制這類呼叫。請將您嵌入的任何 iframe 視為能夠存取已繫結的 Go 方法。

## v3.0.0-beta.34 - 2026-10-20

## 已修正
- 修正 Windows WebView 在視窗啟用時未重新取得焦點，以及 Tab 鍵受困於 WebView 內而無法移至原生視窗再返回的問題（#3783）
- 修正 `internal/webview2/pkg/webview2` COM vtable 結構（`ICoreWebView2_2`..`_27` 及其他衍生介面）遺漏繼承的方法鏈，導致呼叫落在錯誤 vtable 槽位的問題（#6129）
- 修正 Windows 詢問對話方塊忽略自訂按鈕標籤且從不觸發按鈕回呼的問題；自訂按鈕現在會以 TaskDialog 顯示（comctl32 v6 無法使用時，會退回使用 MessageBox，並依位置對應回呼）（#5794）
- 修正搭配 Vite 使用時，`wails3 dev` 會停滯 60 秒的問題：自動前端就緒閘門現在改用 TCP 探測，不再使用未傳送 `Accept` 標頭並收到 404 的 HTTP GET（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 更新開啟中的系統匣選單時，維持其正常運作；見 [PR](https://github.com/wailsapp/wails/pull/6301)，由 @taliesin-ai 貢獻
- 在執行階段操作中妥善處理暫時性 WebView2 COM 錯誤；見 [PR](https://github.com/wailsapp/wails/pull/6302)，由 @leaanthony 貢獻

## v3.0.0-beta.33 - 2026-10-19

## 已新增
- 繫結產生器現在會輸出 `bindings/events` 模組，為每個使用編譯期常數名稱向 `application.RegisterEvent` 註冊且可探索的事件匯出具型別的名稱常數（例如，為 `user:login` 匯出 `UserLogin`）；無須 Vite 外掛程式即可搭配 `Events.On`/`Emit` 使用（#5645）

## 已修正
- 修正 Windows 詢問對話方塊忽略自訂按鈕標籤且從不觸發按鈕回呼的問題；自訂按鈕現在會以 TaskDialog 顯示（comctl32 v6 無法使用時，會退回使用 MessageBox，並依位置對應回呼）（#5794）
- 修正搭配 Vite 使用時，`wails3 dev` 會停滯 60 秒的問題：自動前端就緒閘門現在改用 TCP 探測，不再使用未傳送 `Accept` 標頭並收到 404 的 HTTP GET（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 更新開啟中的系統匣選單時，維持其正常運作；見 [PR](https://github.com/wailsapp/wails/pull/6301)，由 @taliesin-ai 貢獻
- 在執行階段操作中妥善處理暫時性 WebView2 COM 錯誤；見 [PR](https://github.com/wailsapp/wails/pull/6302)，由 @leaanthony 貢獻

## v3.0.0-beta.32 - 2026-10-18

## 已新增
- 為 macOS 14+ 新增原生 macOS 選單區段標頭——請參閱[文件網站](https://v3.wails.io)（已更新 20 個頁面）；見 [PR](https://github.com/wailsapp/wails/pull/6297)，由 @leaanthony 貢獻

## 已修正
- 使 SetRelativePosition 在 macOS 上相對於螢幕工作區定位；見 [PR](https://github.com/wailsapp/wails/pull/6299)，由 @leaanthony 貢獻
- 父程序逾時或取代失敗後，更新程式輔助工具會重新啟動未變更的目標；見 [PR](https://github.com/wailsapp/wails/pull/6298)，由 @taliesin-ai 貢獻
- 將版本化文件使用的 pnpm 鎖定至 10 版；見 [PR](https://github.com/wailsapp/wails/pull/6293)，由 @taliesin-ai 貢獻

## v3.0.0-beta.31 - 2026-10-18

## 已新增
- 公開 `@wailsio/runtime/cancellable` 子路徑匯出，讓 `CancellablePromise` 無須求值根執行階段進入點即可匯入（#6153）

## 已變更
- 要求使用 Go 1.26 或更新版本：v3 模組現在宣告 `go 1.26.0`，由 `wails3 init` 和 `wails3 service init` 產生的專案也宣告相同的最低版本

## 已修正
- 修正使用 Go 1.27 建置時，繫結產生器輸出發生變化的問題；當時 `encoding/json` 型別（例如 `Marshaler`、`Unmarshaler` 和 `RawMessage`）成為 `encoding/json/v2` 和 `encoding/json/jsontext` 型別的別名：現在會再次將其呈現為 `any`，且不會為 `encoding/json/v2` 輸出模型（#6292）
- 複製前，驗證傳給 `wails3 init -t` 的遠端範本 URL：僅接受 `https://`、`http://`、`ssh://`、`git://` 和 `user@host:path` URL；值會在 `--` 之後傳給 git，且複製期間會將 git 限制為僅能使用這些傳輸方式
- 拒絕 HTTP 執行階段傳輸中超出範圍的 `x-wails-window-id` 標頭值，而不是無聲地將其截斷成另一個視窗的 ID
- 強化更新程式的封存檔解壓縮，以防範 zip-slip：現在會使用 `filepath.IsLocal` 檢查項目路徑，並對聯結後的路徑執行包含範圍檢查；測試涵蓋 zip 與 tar.gz 封存檔中的巢狀 `..`、反斜線、絕對路徑及透過符號連結寫入的項目
- 修正來自`wails3 generate appimage`的 AppImage 在 Debian 與 Ubuntu 以外的系統啟動時，因「Unable to spawn a new child process: … WebKitNetworkProcess」而當機的問題。隨附的 WebKit 現在會執行 AppImage 中提供的輔助程序。WebKit 的 bubblewrap 沙箱仍會保持啟用
- 修正`-tags server`組建在 Android 與 iOS 上編譯失敗的問題，並修復[#6271](https://github.com/wailsapp/wails/issues/6271)
- 修正更新程式在暫存目錄位於另一個檔案系統（例如 tmpfs `/tmp`）時無法於 Linux 上置換檔案的問題：現在會先在已安裝的二進位檔旁暫存檔案，再跨檔案系統複製
- 修正更新程式嘗試取代 AppImage 的唯讀掛載點，而非`.AppImage`檔案的問題
- 修正 macOS 與 iOS 上保留的`wails:`視窗控制訊息可從 iframe 接收的問題。現在只會接受來自視窗主框架的這類訊息；iframe 的其他訊息仍會連同來源資訊傳送至`RawMessageHandler`
- 停止在 macOS 上發送`WindowDidUpdate`與`WindowWillUpdate`：`NSWindowDidUpdateNotification`會在每次執行迴圈處理時觸發，因此頁面每秒會收到約120次空的`evaluateJavaScript:`呼叫；而在 macOS 上，27.0的每次呼叫都會耗用使用者點擊所產生的暫時啟用狀態，導致來自真實點擊的`navigator.clipboard.writeText()`及其他受啟用狀態限制的 API 以`NotAllowedError`拒絕執行（[#6170](https://github.com/wailsapp/wails/issues/6170)），由 @APshenkin 貢獻
- 將 Windows info.json 的預設語言修正為 en-US（0409），使檔案總管能顯示版本詳細資料（#6210）
- 修正配備 NVIDIA 獨立 GPU 與 AMD 或 Intel 整合式 GPU 的筆記型電腦會停用 Linux DMA-BUF 轉譯器的問題：在 Wayland 上，現在只有 WebKit 使用 NVIDIA 轉譯時（應用程式已卸載至獨立 GPU，或合成器在其上執行）才會停用，而不是只要載入 nvidia 模組就停用；後者會使整合式 GPU 上的轉譯慢到幾乎無法使用。X11 工作階段仍保留舊行為，因為 Xorg 的轉譯 GPU 是由其組態設定
- 在[PR](https://github.com/wailsapp/wails/pull/6135)中，為 Linux 更新程式的跨裝置重新命名失敗加入暫存複製的備援處理，由 @john-okeefe 貢獻
- 修正 Windows 上的`WebviewWindow.SetMenu`只有在設定選單列佈景主題時才會重新繪製選單列的問題，並修復[#6103](https://github.com/wailsapp/wails/issues/6103)
- 修正`WebviewWindow.SetMenu(nil)`在 Windows 上引發 panic，而未像 macOS 上一樣不執行任何操作的問題，並修復[#6104](https://github.com/wailsapp/wails/issues/6104)
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中加入重試機制，以在相互競爭的推送之間保留變更記錄項目，由 @leaanthony 貢獻
- 移除參照模組中從未提供之二進位檔的嵌入項目，修正`go mod vendor`在所有平台上皆因`pattern arm64/WebView2Loader.dll: no matching files found`而失敗的問題，並修復[#5782](https://github.com/wailsapp/wails/issues/5782)與[#5376](https://github.com/wailsapp/wails/issues/5376)；變更位於[PR](https://github.com/wailsapp/wails/pull/6031)，由 @Grantmartin2002 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6159)中防止 macOS 重複拆除視窗，由 @charlesbalderson 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6162)中處理系統匣選單的 WM_CONTEXTMENU，使其可按右鍵開啟，由 @wzk2025 貢獻
- 修正所有五種`InvokeSync*`變體在分派的回呼發生 panic 時永久阻塞呼叫端的問題，並修復[#6107](https://github.com/wailsapp/wails/issues/6107)；變更位於[PR](https://github.com/wailsapp/wails/pull/6108)，由 @RALIST 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6236)中安全地重建已銷毀的 GTK3 子選單，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6261)中實作 Windows worker 轉送墊片，以支援巢狀 worker，由 @taliesin-ai 貢獻
- 修正 Windows 上的`OpenFile()`在`CanChooseDirectories(true)`時忽略`PromptForMultipleSelection()`的問題：使用者現在可像在 macOS 與 Linux 上一樣，一次選取多個資料夾（#6263）
- 修正 windows/arm64 上以值傳遞給 WebView2 COM 方法的`double`引數問題；受影響的項目包括`ZoomFactor`、`RasterizationScale`、Cookie `Expires`，以及列印設定的縮放比例、頁面大小與邊界（列印設定 setter 先前會傳回`ErrDoubleArgUnsupported`）。受影響的呼叫現在會透過可載入`d0`的`runtime.cgocall`跳板執行，而 Go 的 syscall 路徑無法載入該值（golang.org/issue/62583），由 @APshenkin 貢獻

## 安全性
- 伺服器模式（`-tags server`）現在只接受以`POST`要求發出的`/wails/runtime`呼叫，且要求必須來自伺服器自身的來源，或`ServerOptions.WebSocketOriginPatterns`中列出的來源；其他方法會收到 HTTP 405，其他瀏覽器來源則會收到 HTTP 403。若部署環境位於會重寫`Host`標頭的反向 Proxy 後方，應將其公開來源加入`WebSocketOriginPatterns`
- 繫結至迴送位址（預設為`localhost`、`127.0.0.0/8`或`::1`）的伺服器模式（`-tags server`），現在會拒絕要求並回傳 HTTP 403，除非`Host`標頭為`localhost`、IP 常值、繫結主機，或符合`ServerOptions.WebSocketOriginPatterns`的主機；這可阻擋針對`/wails/runtime`與 WebSocket 端點的 DNS 重新繫結攻擊
- Windows 通知在 WinRT 快顯通知 API 失敗時不再退回使用 PowerShell，而是傳回錯誤

## v3.0.0-beta.30 - 2026-10-15

## 已修正
- 修正更新程式嘗試取代 AppImage 的唯讀掛載點，而非`.AppImage`檔案的問題
- 停止在 macOS 上發送`WindowDidUpdate`與`WindowWillUpdate`：`NSWindowDidUpdateNotification`會在每次執行迴圈處理時觸發，因此頁面每秒會收到約120次空的`evaluateJavaScript:`呼叫；而在 macOS 上，27.0的每次呼叫都會耗用使用者點擊所產生的暫時啟用狀態，導致來自真實點擊的`navigator.clipboard.writeText()`及其他受啟用狀態限制的 API 以`NotAllowedError`拒絕執行（[#6170](https://github.com/wailsapp/wails/issues/6170)），由 @APshenkin 貢獻
- 將 Windows info.json 的預設語言修正為 en-US（0409），使檔案總管能顯示版本詳細資料（#6210）
- 在[PR](https://github.com/wailsapp/wails/pull/6135)中，為 Linux 更新程式的跨裝置重新命名失敗加入暫存複製的備援處理，由 @john-okeefe 貢獻
- 修正 Windows 上的`WebviewWindow.SetMenu`只有在設定選單列佈景主題時才會重新繪製選單列的問題，並修復[#6103](https://github.com/wailsapp/wails/issues/6103)
- 修正`WebviewWindow.SetMenu(nil)`在 Windows 上引發 panic，而未像 macOS 上一樣不執行任何操作的問題，並修復[#6104](https://github.com/wailsapp/wails/issues/6104)
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中加入重試機制，以在相互競爭的推送之間保留變更記錄項目，由 @leaanthony 貢獻
- 移除參照模組中從未提供之二進位檔的嵌入項目，修正`go mod vendor`在所有平台上皆因`pattern arm64/WebView2Loader.dll: no matching files found`而失敗的問題，並修復[#5782](https://github.com/wailsapp/wails/issues/5782)與[#5376](https://github.com/wailsapp/wails/issues/5376)；變更位於[PR](https://github.com/wailsapp/wails/pull/6031)，由 @Grantmartin2002 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6159)中防止 macOS 重複拆除視窗，由 @charlesbalderson 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6162)中處理系統匣選單的 WM_CONTEXTMENU，使其可按右鍵開啟，由 @wzk2025 貢獻
- 修正所有五種`InvokeSync*`變體在分派的回呼發生 panic 時永久阻塞呼叫端的問題，並修復[#6107](https://github.com/wailsapp/wails/issues/6107)；變更位於[PR](https://github.com/wailsapp/wails/pull/6108)，由 @RALIST 貢獻

## v3.0.0-beta.29 - 2026-10-13

## 已新增
- macOS 上的全域快速鍵現在會考量鍵盤配置：在非 QWERTY 配置（AZERTY、QWERTZ、Dvorak 等）上，字母快速鍵會透過`UCKeyTranslate`繫結至標示該字母的按鍵，而非固定的 QWERTY 實體按鍵位置。

## 已修正
- 在[PR](https://github.com/wailsapp/wails/pull/6233)中，先初始化原生選單再於 Linux 上設定選單，由 @abhicodes-007 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6237)中解決 macOS 正式發行就緒度在輸出處理與設定精靈方面的退步問題，由 @taliesin-ai 貢獻
- 修正 Windows 上的 WebView2 COM `AddRef`/`Release`包裝函式將過時的執行緒最後錯誤回報為失敗的問題。它們現在會傳回參考計數，因此釋放作業不再會中止應用程式或導致資產要求失敗（#3931）
- 修正`ICoreWebView2Deferral.Release`呼叫`AddRef`插槽，該問題會在 Windows 上洩漏每個延後處理的資產請求之延後作業（#3931）

## v3.0.0-beta.28 - 2026-10-05

## 已修正
- 在[PR](https://github.com/wailsapp/wails/pull/6225)中將 macOS 應用程式輸出保留在檔案提供者之外，由 @taliesin-ai 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6224)中完成設定精靈並顯示系統詳細資料，由 @taliesin-ai 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6222)中改善 Windows 上的 WebView2 復原錯誤處理與資源清理，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6220)中改善 Windows 上的 WebView2 復原，以防止應用程式當機及控制器擱置，由 @taliesin-ai 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6221)中修正 macOS 媒體授權指引，由 @taliesin-ai 貢獻
- 修正以滑鼠右鍵按一下 Linux 系統匣圖示時，除了開啟選單還會執行點擊處理常式，導致附加的視窗在每次按右鍵時切換狀態的問題（#6018）
- 在[PR](https://github.com/wailsapp/wails/pull/6184)中讓 Windows 上的 WebView2 回呼處理常式保持存活，直到原生端釋放為止，由 @taliesin-ai 貢獻

## v3.0.0-beta.27 - 2026-10-01

## 已修正
- 修正 macOS 12及更新版本忽略`Permissions`選項的問題：現已實作`WKUIDelegate`媒體擷取方法，因此`PermissionAllow`與`PermissionDeny`會像在 Linux 和 Windows 上一樣套用至相機與麥克風請求。仍須提供`NSCameraUsageDescription`／`NSMicrophoneUsageDescription`；若在沙箱中執行，也仍須提供相符的裝置權利

## v3.0.0-beta.26 - 2026-09-25

## 已新增
- 監看器在[PR](https://github.com/wailsapp/wails/pull/6150)中對前端開發伺服器套用 HTTP 就緒檢查，由 @atterpac 貢獻

## 已修正
- 在[PR](https://github.com/wailsapp/wails/pull/6002)中改為從 WebView2 處理程序失敗中復原，而非留下空白視窗，由 @taliesin-ai 貢獻

## v3.0.0-beta.25 - 2026-09-22

## 已修正
- 主要應用程式結束時停止`wails3 dev`及其背景處理程序（#6048）

## v3.0.0-beta.24 - 2026-09-20

## 安全性
- 限制建置與發布工作流程中的 GITHUB_TOKEN 權限，以遵循最小權限原則
- 修正螢幕範例資產中介軟體中的路徑遊走及符號連結逸出問題
- 在將設定精靈的相依套件安裝限制為核准的套件管理器命令時，維持其正常運作

## v3.0.0-beta.23 - 2026-09-16

## 已新增

- 在[PR](https://github.com/wailsapp/wails/pull/5972)中允許覆寫 Linux 上的 GTK 應用程式 ID，由 @overlordtm 貢獻

## 已修正

- 在[PR](https://github.com/wailsapp/wails/pull/5973)中讓 Linux 單一執行個體以應用程式本身的 ID 宣告匯流排名稱，由 @overlordtm 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6126)中取消 Linux 上已中止的自訂配置請求，由 @leaanthony 貢獻

## v3.0.0-beta.22 - 2026-09-14

## 已新增

- 在[PR](https://github.com/wailsapp/wails/pull/6122)中完成翻譯並移植社群媒體內容——請參閱[文件網站](https://v3.wails.io)（已更新1812個頁面），由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6120)中記錄可嵌入且具備多種樣式與尺寸的 Wails 按鈕，由 @leaanthony 貢獻

## v3.0.0-beta.21 - 2026-09-13

## 已新增

- 在[PR](https://github.com/wailsapp/wails/pull/6116)中使用 M-Press 提供 Wails v3文件，由 @leaanthony 貢獻

## 已修正

- 在[PR](https://github.com/wailsapp/wails/pull/6118)中剖析 MPD 前置資料的 JSON slug 值，以產生變更記錄，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6080)中讓更新程式在備份失敗後清除輔助程式環境變數，並重新啟動原始目標，由 @cnmax 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6098)中於 App.Run 期間啟動預設訊號處理常式，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6112)中讓 Windows 選單處理 nil 選單、釋放被取代的資源，並重新繪製選單列，由 @taliesin-ai 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6115)中使用共用 YAML 設定，恢復新建專案的 MSIX 封裝功能，由 @leaanthony 貢獻
- 修正產生的 JavaScript 與 TypeScript 繫結在泛型模型建立器參照後續宣告的輔助程式時無法載入的問題，並防止建立彼此相依的泛型模型時發生堆疊溢位（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 已變更

- 在[PR](https://github.com/wailsapp/wails/pull/6082)中將 Clave 展示連結更新為目前的網站與儲存庫，由 @01xR4in 貢獻

## 已修正

- 取消已中止的 Windows 資產請求（包括工作執行緒請求），同時在導覽期間保留 keepalive 處理常式。在 Apple 平台上，透過應用程式包裝器轉送原生請求內容。（#5963、#5969）
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中透過重試機制，在競爭推送之間保留變更記錄項目，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6031)中移除參照模組從未隨附之二進位檔的嵌入項目，修正`go mod vendor`在所有平台上皆因`pattern arm64/WebView2Loader.dll: no matching files found`而失敗的問題，並解決[#5782](https://github.com/wailsapp/wails/issues/5782)與[#5376](https://github.com/wailsapp/wails/issues/5376)，由 @Grantmartin2002 貢獻

## 已移除

- 移除已由純 Go 載入器取代的原生 WebView2 載入器支援。這會移除內嵌的`WebView2Loader.dll`二進位檔與`github.com/jchv/go-winloader`相依套件。`native_webview2loader`建置標籤仍會被接受且不再引發錯誤，但對v3建置沒有作用；此變更見[PR](https://github.com/wailsapp/wails/pull/6031)，由 @Grantmartin2002 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6097)中從 macOS API 指南移除未使用的建置標籤與 FPS 選項，由 @leaanthony 貢獻

## v3.0.0-beta.19 - 2026-09-09

## 已新增

- 在[PR](https://github.com/wailsapp/wails/pull/6087)中以建置標籤限制私有 macOS API，使其僅能選擇性啟用——請參閱[文件](https://v3.wails.io/features/browser/integration)、[文件](https://v3.wails.io/features/environment/info)、[文件](https://v3.wails.io/features/windows/basics)、[文件](https://v3.wails.io/features/windows/frameless)、[文件](https://v3.wails.io/features/windows/notch-windows)、[文件](https://v3.wails.io/features/windows/options)、[文件](https://v3.wails.io/guides/build/macos)、[文件](https://v3.wails.io/guides/build/private-macos-apis)及[文件](https://v3.wails.io/reference/overview)，由 @leaanthony 貢獻

## 已修正

- 在[PR](https://github.com/wailsapp/wails/pull/6091)中以 HTTP 413拒絕超過64 MiB 的執行階段請求，由 @leaanthony 貢獻

## 安全性

- 在[PR](https://github.com/wailsapp/wails/pull/6092)中透過權杖驗證強化 MCP 來源與遠端存取安全性，由 @leaanthony 貢獻

## v3.0.0-beta.18 - 2026-09-08

## 已修正

- 透過使用指標接收器，修正 Linux 與 Darwin 上的 Calloc 記憶體洩漏，由 @4RH1T3CT0R7 於[PR](https://github.com/wailsapp/wails/pull/6083)中完成

## v3.0.0-beta.17 - 2026-09-06

## 已修正

- Windows：WebResourceRequested 處理常式中的`GetRequest`失敗或為 nil 時，不再終止處理程序（`log.Fatal`／nil 解參照 panic）；現在會捨棄該要求並記錄日誌，由 @midagedev 於[PR](https://github.com/wailsapp/wails/pull/6006)中完成

## v3.0.0-beta.16 - 2026-08-29

## 已變更

- 在新的終端機視窗中提示輸入公證密碼，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/6029)中完成

## 已修正

- 正確處理 macOS 上的系統匣點按類型，由 @ChewbaccaCookie 於[PR](https://github.com/wailsapp/wails/pull/5919)中完成
- CI 會在更新前移除未使用的 Microsoft apt 軟體庫，由 @Grantmartin2002 於[PR](https://github.com/wailsapp/wails/pull/6041)中完成

## v3.0.0-beta.15 - 2026-08-27

## 已修正

- 將 WebView2 嵌入逾時時間增加至60秒，由 @Grantmartin2002 於[PR](https://github.com/wailsapp/wails/pull/6043)中完成

## v3.0.0-beta.14 - 2026-08-26

## 已修正

- 在 macOS 上正確命名 Control-字母按鍵輸入，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6032)中完成
- 修正 Windows 上的 ICO 系統匣圖示，並使其跟隨工作列佈景主題，由 @nik9play 於[PR](https://github.com/wailsapp/wails/pull/6016)中完成

## v3.0.0-beta.13 - 2026-08-25

## 已修正

- 在 macOS 執行模態迴圈期間持續處理主執行緒工作，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/6026)中完成
- 讓行動裝置安全儲存操作可回報失敗，並在失敗時預設拒絕存取，由 @mortenolsrud 於[PR](https://github.com/wailsapp/wails/pull/5923)中完成
- 即使未註冊監聽器，也會執行應用程式事件掛鉤，由 @archy-rock3t-cloud 於[PR](https://github.com/wailsapp/wails/pull/5999)中完成
- 修正註解與本地化文件中的錯字，由 @haoku123 於[PR](https://github.com/wailsapp/wails/pull/6023)中完成
- 移除提交於`v3/examples`下的預先編譯 macOS 二進位檔，由 @4RH1T3CT0R7 於[PR](https://github.com/wailsapp/wails/pull/6025)中完成

## v3.0.0-beta.12 - 2026-08-21

## 已新增

- 新增具備生命週期和遙測範例的 macOS 瀏海通知視窗——請參閱[文件](https://v3.wails.io/features/windows/notch-windows)，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/6010)中完成
- 新增 macOS NSPanel 視窗支援，包括新選項與原生整合——請參閱[文件](https://v3.wails.io/features/windows/options)，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/6008)中完成

## 已修正

- 防止 SQLite Prepare 在並行呼叫下停滯，由 @archy-rock3t-cloud 於[PR](https://github.com/wailsapp/wails/pull/5998)中完成

## v3.0.0-beta.11 - 2026-08-20

## 已移除

- 從文件移除已過時的實作追蹤器，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/6005)中完成

## v3.0.0-beta.10 - 2026-08-19

## 已修正

- 修正 GTK4 Linux 主機遺漏自訂通訊協定與檔案關聯啟動引數的問題，由 @midagedev 於[PR](https://github.com/wailsapp/wails/pull/6000)中完成
- 在變更日誌驗證中正確處理已刪除的行與同一來源的修正，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5993)中完成

## v3.0.0-beta.9 - 2026-08-16

## 已新增

- 新增安全的 wails3 mcp 伺服器，用於代理程式輔助的專案管理，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5896)中完成
- 新增繫結中模型的文件——請參閱[文件](https://v3.wails.io/features/bindings/models)，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5988)中完成
- 支援在原子化 Linux 系統上透過 rpm-ostree 安裝，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5987)中完成
- 新增原生的每週星標歷史圖表產生與發布功能——請參閱[文件](https://v3.wails.io/credits)、[文件](https://v3.wails.io/de/credits)、[文件](https://v3.wails.io/fr/credits)、[文件](https://v3.wails.io/id/credits)、[文件](https://v3.wails.io/ja/credits)、[文件](https://v3.wails.io/ko/credits)、[文件](https://v3.wails.io/pt/credits)、[文件](https://v3.wails.io/ru/credits)、[文件](https://v3.wails.io/zh-cn/credits)及[文件](https://v3.wails.io/zh-tw/credits)，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5986)中完成
- 新增僅適用於 Darwin、用來解析應用程式套件資源的 mac 套件——請參閱[文件](https://v3.wails.io/guides/build/macos)，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5965)中完成
- 新增 Condui 展示頁面與索引項目——請參閱[文件](https://v3.wails.io/community/showcase/condui)及[文件](https://v3.wails.io/community/showcase)，由 @mgueregath 於[PR](https://github.com/wailsapp/wails/pull/5962)中完成
- 新增包含螢幕擷取畫面與專案連結的 Redis Viewer 展示頁面——請參閱[文件](https://v3.wails.io/community/showcase)及[文件](https://v3.wails.io/community/showcase/redisviewer)，由 @redisviewer 於[PR](https://github.com/wailsapp/wails/pull/5984)中完成

## 已變更

- 將 Linux 的 GTK 應用程式旗標更新為 G<em>APPLICATION</em>NON_UNIQUE，由 @overlordtm 於[PR](https://github.com/wailsapp/wails/pull/5971)中完成
- 以偵錯層級記錄遺漏的視窗事件，而非警告層級，由 @julianstorer 於[PR](https://github.com/wailsapp/wails/pull/5914)中完成

## 已修正

- 允許已註冊的 macOS 鍵盤快速鍵優先於 WebView 處理，由 @julianstorer 於[PR](https://github.com/wailsapp/wails/pull/5902)中完成
- 修復文件側邊欄中失效的連結，由 @northes 於[PR](https://github.com/wailsapp/wails/pull/5937)中完成
- 當 WebKit 中止相符的自訂配置任務時，取消 macOS 與 iOS 的資產要求情境 (#5963)
- 處理 WindowSetFullscreenButtonEnabled 訊息，由 @archy-rock3t-cloud 於[PR](https://github.com/wailsapp/wails/pull/5976)中完成
- 在 preact-ts 範本中匯入 Fragment 以解決建置失敗，由 @haoku123 於[PR](https://github.com/wailsapp/wails/pull/5979)中完成
- 防止舊版 GTK3 純服務應用程式在螢幕探索早於可用的作用中視窗或顯示器執行時崩潰 (#5966)
- 當未發布的變更日誌為空時，允許明確指定版本的發布執行繼續進行 (#5977)

## 安全性

- 將網站的 nanoid 鎖定檔更新至已修補的3.3.18，以解決安全性公告指出的問題，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5985)中完成

## v3.0.0-beta.8 - 2026-08-12

## 已新增

- 為自動變更日誌項目新增文件 URL 產生功能，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5957)中完成
- 新增 Streams：採用 WebSocket 程式設計模型，在 Go 與 JavaScript 之間提供雙向位元組串流，且無須監聽通訊端。在 Go 中使用`app.HandleStream(name, handler)`宣告串流，並從前端使用`Stream(name)`連線；此方法會傳回外形如`WebSocket`的物件。Go→JS 透過資產伺服器為每個視窗保留一個輪詢來傳輸，JS→Go 則使用一般 POST；不會繫結任何 TCP 連接埠，也不會有任何資料經過`evaluateJavaScript`。在伺服器組建（`-tags server`）中，同一個處理常式改由真正的 WebSocket 提供，因此應用程式碼在各種組建中皆相同。由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5935)中將信箱變更日誌項目移至「未發布」，由 @leaanthony 貢獻

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5938)中更新文件側邊欄的自動產生機制與部落格作者型別的衍生方式，由 @leaanthony 貢獻

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5952)中讓 WebView2 初始化使用期限與訊息幫浦，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5951)中讓 WebView2 Cookie 測試在 CI 中除非明確選用否則跳過，並將執行鎖定於目前的作業系統執行緒，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5944)中讓 Windows 選單建構器恢復子選單父項目的命令 ID，由 @gilad-ch 貢獻
- 使官方交叉編譯映像與 GTK 4.14+ Linux 支援基準一致（#5928）
- 在[PR](https://github.com/wailsapp/wails/pull/5915)中設定 iOS Xcode 專案保留繼承的連結器旗標並加入 -ObjC，由 @mortenolsrud 貢獻
- 修正大型前端中`wails3 dev`資產 Proxy 過度頻繁建立與中斷 TCP 連線的問題；此問題可能耗盡主機的暫時連接埠，並導致不相關的處理程序因`EADDRNOTAVAIL`而失敗
- 在[PR](https://github.com/wailsapp/wails/pull/5934)中依視窗將事件 JavaScript 排入佇列，以確保依序分派並提供背壓，由 @leaanthony 貢獻

## 移除

- 移除桌面二進位檔發布管線：v3發布僅包含標籤，而`wails3`CLI 則使用`go install`安裝。在[PR](https://github.com/wailsapp/wails/pull/5946)中刪除`release-v3.yml`以及分派該管線的夜間步驟，由 @leaanthony 貢獻

## v3.0.0-beta.7 - 2026-08-11

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5512)中新增 macOS 自動播放偏好設定，可停用媒體播放需要使用者操作的要求，由 @Eyalm321 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5935)中將信箱變更日誌項目移至「未發布」，由 @leaanthony 貢獻

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5945)中讓 macOS 縮放動畫使用 CADisplayLink 或 NSTimer，以提供更流暢的效能，由 @savely-krasovsky 貢獻

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5915)中設定 iOS Xcode 專案保留繼承的連結器旗標並加入 -ObjC，由 @mortenolsrud 貢獻
- 修正大型前端中`wails3 dev`資產 Proxy 過度頻繁建立與中斷 TCP 連線的問題；此問題可能耗盡主機的暫時連接埠，並導致不相關的處理程序因`EADDRNOTAVAIL`而失敗
- 在[PR](https://github.com/wailsapp/wails/pull/5934)中依視窗將事件 JavaScript 排入佇列，以確保依序分派並提供背壓，由 @leaanthony 貢獻

### 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5851)中實作通用非同步 FIFO 信箱，以依序傳遞事件，由 @savely-krasovsky 與 @DevLumuz 貢獻

## v3.0.0-beta.6 - 2026-08-09

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5930)中實作有界的主機端儲存空間來存放過大的事件，並依序傳遞 JavaScript，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5921)中實作 macOS Dock 圖示彈跳來提示視窗，由 @julianstorer 貢獻

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5931)中讓資產伺服器在排清期間保留內容類型偵測錯誤及尚未寫入的前置資料，由 @leaanthony 貢獻
- 防止 macOS 應用程式從 Wails 回呼取代應用程式選單時當機
- 修正 Windows 10 1809 / Windows Server 2019（組建17763）上無法辨讀的原生選單。深色模式的 uxtheme 匯出項目受限於組建18334，因此應用程式層級的深色模式選用設定從未在這些主機上執行：選單背景繪製為深色，但 Windows 仍以淺色佈景主題繪製選單文字，導致深色背景上出現深色文字。這些序號從17763起便已存在，因此現在的條件限制已與其相符。
- 修正`w32.GetStockObject`呼叫`GetDeviceCaps`而非`GetStockObject`的問題；此問題會使其對每個預設物件都傳回0。
- 在[PR](https://github.com/wailsapp/wails/pull/5924)中改善 WebView2 啟動載入器的下載錯誤處理與回報，由 @jannskiee 貢獻

## v3.0.0-beta.5 - 2026-08-07

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5897)中讓 macOS 應用程式啟用僅對一般應用程式遵循啟用原則，由 @julianstorer 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5898)中防護 Linux 組建內尚未初始化的 GTK 視窗，由 @julianstorer 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5899)中於載入 URL 前，為 Linux WebKit 視窗設定明確的不透明背景色彩，由 @julianstorer 貢獻

## v3.0.0-beta.4 - 2026-08-05

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5890)中讓 Android 組建任務預設使用 arm64，並讓 deploy-emulator 選取主機架構，由 @mortenolsrud 貢獻

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5900)中於拖曳期間保留 macOS 視窗的縮放狀態並減少動態效果，由 @leaanthony 貢獻
- 將`!server`加入`webview_window_windows_nonclient.go`組建條件約束，以修正 Windows 伺服器模式組建

## v3.0.0-beta.3 - 2026-08-03

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5881)中，於實作詳細資料內記錄第10階段 Beta 驗證已完成，由 @leaanthony 貢獻

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5877)中將視窗控制代碼傳給 Windows 深色模式 API 並驗證引數，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5870)中集中處理 macOS 無框線視窗的標題列按鈕狀態解析，由 @taliesin-ai 貢獻
- 防止 Windows 應用程式要求使用深色模式、但 Windows 應用程式佈景主題為淺色時，原生選單文字無法辨讀。在 Windows 能夠呈現深色選單文字前，選單現在會使用與其相符的淺色原生背景。
- 修正 Windows 10 1809 / Windows Server 2019（組建17763）上無法辨讀的原生選單。深色模式的 uxtheme 匯出項目受限於組建18334，因此應用程式層級的深色模式選用設定從未在這些主機上執行：選單背景繪製為深色，但 Windows 仍以淺色佈景主題繪製選單文字，導致深色背景上出現深色文字。這些序號從17763起便已存在，因此現在的條件限制已與其相符。

## v3.0.0-beta.2 - 2026-08-02

## 變更

- 將v3從 Alpha 升級為 Beta
- 記錄系統匣智慧型預設值與快顯視窗自動隱藏行為，並新增點擊處理常式選取的迴歸測試涵蓋範圍（#5840）。
- 在[PR](https://github.com/wailsapp/wails/pull/5861)中讓 GitHub 更新程式預設排除 Windows 安裝程式資產，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/5866)中新增對 macOS 無框架視窗圓角、方角及自訂圓角半徑的支援，作者為 @leaanthony

## 修正

- 回報即時 GTK4 視窗大小，並從設定的表面發出調整大小、最大化、最小化及全螢幕狀態事件（#5830）。
- 在[PR](https://github.com/wailsapp/wails/pull/5854)中修正於 fetch 要求中傳送 Blob 或 FormData 時造成 Linux WebKit 當機的問題，作者為 @taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5865)中讓 fetch shim 在缺少 Blob/FormData 標頭時傳遞 undefined，作者為 @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## 新增

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5866)中新增對 macOS 無框架視窗圓角、方角及自訂圓角半徑的支援，作者為 @leaanthony

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5865)中讓 fetch shim 在缺少 Blob/FormData 標頭時傳遞 undefined，作者為 @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5857)中新增 macOS DMG 封裝支援，以及新的選項與建置工作，作者為 @leaanthony

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5861)中讓 GitHub 更新程式預設排除 Windows 安裝程式資產，作者為 @leaanthony

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5854)中修正於 fetch 要求中傳送 Blob 或 FormData 時造成 Linux WebKit 當機的問題，作者為 @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5853)中實作 macOS 標題列按兩下以最大化或最小化的動作，作者為 @taliesin-ai

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5849)中更新 QR 服務教學以使用 NewServiceWithOptions，並加入間距，作者為 @jeongkyu

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5856)中讓 WKWebView 在 macOS 縮放期間保持回應，作者為 @leaanthony
- 修正 GTK4 視窗大小查詢，並從設定的`GdkSurface`發出調整大小、最大化、最小化及全螢幕狀態事件。

## v3.0.0-alpha2.119 - 2026-07-27

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5833)中更新多語言文件以納入架構圖，作者為 @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5825)中提供圖示產生作業輸入與輸出的預設路徑，作者為 @taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5797)中將來源進入點模組加入執行階段 package.json 的 sideEffects，作者為 @savely-krasovsky
- 在[PR](https://github.com/wailsapp/wails/pull/5816)中於貢獻指南加入授權條款與來源追溯章節，作者為 @taliesin-ai

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5800)中套用範圍限定的 GTK4 無框架 CSS 以移除圓角半徑，作者為 @savely-krasovsky
- 在[PR](https://github.com/wailsapp/wails/pull/5789)中妥善處理 Windows 上快顯功能表與螢幕列舉的游標位置取得失敗，作者為 @wayneforrest
- 在[PR](https://github.com/wailsapp/wails/pull/5678)中讓 macOS 開啟檔案對話方塊正確篩選副檔名，並依後綴驗證允許的檔案，作者為 @phergul
- 在[PR](https://github.com/wailsapp/wails/pull/5793)中防止 Windows 深色模式初始化呼叫 nil API，作者為 @roachadam
- 修正更新程式的32位元建置失敗：在`GOARCH=386`上將`maxArchiveTotalSize`常數（2 GiB）傳給`fmt.Errorf`時，該常數會溢位平台的`int`。現在已明確將其型別指定為`int64`。
- 修正啟動時發生的 nil 指標 panic：當視窗使用深色（或系統深色）標題列，而 Windows 組建未載入深色模式 uxtheme API（例如 Windows 10 1809 / Windows Server 2019（組建17763））時，便會發生此問題。視窗佈景主題設定中的`AllowDarkModeForWindow`呼叫現在會檢查 nil，與`w32.SetMenuTheme`中既有的防護一致。

## v3.0.0-alpha2.117 - 2026-07-08

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5462)中為 Windows 的非用戶端區域實作自訂命中測試邏輯，作者為 @savely-krasovsky

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5761)中依據 UseVisualHosting 設定 WebView2 顯示器縮放比例偵測，作者為 @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5763)中更新 FAQ 文件，使其聚焦於 Wails v3的功能與指引，作者為 @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## 修正

- 修正`Menu.Update()`未在 GTK4 Linux 上重新建置原生功能表的問題（#5659；@puneetdixit200 另行診斷並於 #5539中修正）
- 透過複製螢幕 ID/名稱字串並建立數量快照，修正在顯示器變更時列舉 macOS 螢幕所造成的當機（#5565；@x-haose 另行診斷並於 #5584中修正）
- 修正 Windows 上的當機：在最小化/還原轉換期間`WM_ERASEBKGND`繪製純色背景，而`GetClientRect`傳回 nil 時會發生此問題（防護措施由 @sinspired 於 #5636回報）
- 修正前端繫結錯誤一律解析為文字的問題，由 @mbaklor 於 #5690中完成
- 修正使用`server`建置標籤時 Windows 建置失敗的問題；原因是 Windows GUI 檔案缺少 macOS 與 Linux 對應檔案已有的`!server`建置條件約束（#5680）

## v3.0.0-alpha2.114 - 2026-07-05

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5720)中實作更新資訊清單通訊協定與端點提供者，作者為 @taliesin-ai

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5711)中將`webview2`繫結以`v3/internal/webview2`的形式整合至v3模組，並移除獨立模組、其每夜發行/同步工作流程，以及 go.mod 版本調整流程（v3是其唯一的使用端），作者為 @taliesin-ai

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5750)中將 WebView2 顯示器縮放比例偵測與 DPI 變更時的主機重新同步修正移至「尚未發行」章節，作者為 @taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5741)中更新 WebView2 對 float64 與 BOOL 參數的 COM 封送處理，作者為 @wayneforrest
- 在[PR](https://github.com/wailsapp/wails/pull/5703)中防止 Windows 系統匣圖示更新與銷毀期間發生 panic 及 nil 解參考，作者為 @wayneforrest
- 在[PR](https://github.com/wailsapp/wails/pull/5743)中修正 Windows 上隱藏的視窗未正確再次隱藏的問題，作者為 @wayneforrest
- 在[PR](https://github.com/wailsapp/wails/pull/5742)中使 WebView2 控制器的可見性與視窗的最小化、最大化及還原同步，作者為 @wayneforrest

### 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5734)中重新啟用 WebView2 顯示器縮放比例偵測，並僅在 DPI 變更時重新同步宿主；由 @taliesin-ai 提交，依據 @randalmurphal 驗證的修正，並由 @eleclin 驗證根本原因、@qq540491950 進行硬體測試

## v3.0.0-alpha2.113 - 2026-07-04

## 新增

- 在未設定`ANDROID_KEYSTORE_FILE`的情況下建置發行版 AAB 時顯示警告（Google Play 會拒絕以偵錯金鑰簽署的套件），並在[PR](https://github.com/wailsapp/wails/pull/5730)中記錄 App Bundle 封裝與簽署方式；由 @taliesin-ai 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5739)中新增多語言的「為何選擇 Wails」文件；由 @taliesin-ai 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5398)中支援在繫結中將 Go time.Time 對應至 JS Date 或字串；由 @fbbdev 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5728)中新增用於提交至 Play Store 的 Android App Bundle（AAB）封裝工作（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）；APK 工作仍用於本機／模擬器測試；由 @mortenolsrud 提交（修正[#5726](https://github.com/wailsapp/wails/issues/5726)）
- 在[PR](https://github.com/wailsapp/wails/pull/5735)中新增 Android 實體裝置工作目標，並恢復相機／位置權限；由 @taliesin-ai 提交

## 變更

- 將`webview2`升級至v1.0.28（[發行說明](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- 在[PR](https://github.com/wailsapp/wails/pull/5730)中，將 Android 範本的`compileSdk`／`targetSdk`從34升級至35，這是 Google Play 對新應用程式提交的要求；由 @taliesin-ai 提交

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5745)中修正 sponsorkit 中預先烘焙的頭像遮罩；由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5730)中修正 Android AVD 自動建立因依字典順序排列版本而選到錯誤系統映像或 cmdline-tools 版本的問題；由 @taliesin-ai 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5730)中修正設定精靈建議過時 Android NDK 版本的問題（現為26.3.11579264，與文件記載的需求一致）；由 @taliesin-ai 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5744)中更新關於 SvelteKit 與選項的法文文件；由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5516)中修正顯示器變更期間於 macOS 列舉螢幕時發生的 SIGSEGV；由 @flofreud 提交

## v3.0.0-alpha2.112 - 2026-07-03

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5724)中新增以 Go 編寫的貢獻者 SVG 產生器，並更新 docs／website 的鳴謝頁面；由 @taliesin-ai 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5398)中支援在繫結中將 Go time.Time 對應至 JS Date 或字串；由 @fbbdev 提交

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5719)中以 Go 產生器取代以 Node 為基礎的贊助者圖片處理流程；由 @taliesin-ai 提交

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5729)中修正 Android 建置資產相依項目的安裝指令碼；由 @taliesin-ai 提交
- 在`ValidateAndSanitizeURL`中拒絕 U+0085（NEXT LINE）控制字元，完成 URL 驗證器對空白字元的涵蓋
- 在[PR](https://github.com/wailsapp/wails/pull/4785)中，於無框線視窗的 DPI 變更時重新計算 DWM 框架；由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/4632)中修正 Windows 使用非100% 縮放比例時 DnD 放置區域偵測失敗的問題；由 @yulesxoxo 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5714)中，為 Darwin 的對話方塊、選單、系統匣與通知所使用的 Cocoa 物件新增明確的 Objective-C 記憶體管理；由 @taliesin-ai 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5718)中修正 Linux CGO 後端錯誤與系統匣問題；由 @taliesin-ai 提交

## v3.0.0-alpha2.111 - 2026-07-01

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5061)中將 HappyTools 新增至社群展示；由 @Aliuyanfeng 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5643)中新增印尼語地區設定支援與完整文件；由 @triadmoko 提交
- 在[PR](https://github.com/wailsapp/wails/pull/4813)中為 WindowsWindow 新增 DisableMenu 選項；由 @leaanthony 提交

## 變更

- 在[PR](https://github.com/wailsapp/wails/pull/5617)中更新 Taskfile 範本與 CLI，以使用 GOOS 和 ARCH 分派建置／封裝工作；由 @leaanthony 提交

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/5708)中修正 Mac 視窗分頁問題；由 @taliesin-ai 提交

## 移除

- 在[PR](https://github.com/wailsapp/wails/pull/5702)中移除貢獻、功能與指南的德文翻譯 MDX 檔案；由 @taliesin-ai 提交

## v3.0.0-alpha2.110 - 2026-06-30

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5129)中實作 macOS WebView 重新載入與強制重新載入，並新增 WebContent 處理程序終止後的復原機制；由 @wayneforrest 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5396)中新增涵蓋貢獻、功能與指南的完整德文文件；由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5333)中強化通知功能，新增音效、附件、排程與更新 API；由 @popaprozac 提交

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/4785)中，於無框線視窗的 DPI 變更時重新計算 DWM 框架；由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/4632)中修正 Windows 使用非100% 縮放比例時 DnD 放置區域偵測失敗的問題；由 @yulesxoxo 提交

## v3.0.0-alpha2.109 - 2026-06-29

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5026)中為 EventsEmit 文件新增程式碼範例；由 @iamhabbeboy 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5380)中新增 Windows WebView2 視覺化裝載選項；由 @MerIijn 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5536)中將 Klustr 新增至社群展示文件；由 @SametKUM 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5685)中將 Kira 新增至社群展示，並加入新頁面與變更日誌項目；由 @thiennguyen93 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5694)中為 MCP 服務指南新增意見回饋章節；由 @taliesin-ai 提交

## 變更

- 伺服器模式現在具備正式支援的生產環境建置方式，並與桌面版建置工作保持一致（#5693）。`task build:server`預設會建置生產環境二進位檔（`-tags server,production`、`-trimpath`、移除符號），並接受`DEV=true`（開發伺服器）、`OBFUSCATED=true`（garble）及`EXTRA_TAGS`。`task run:server`會執行開發伺服器。`Dockerfile.server`／`task build:docker`會先建置生產環境伺服器（`-tags server,production`）與生產環境前端；映像預設是在 distroless/static 上建置的純 Go 靜態版本，並將`CGO_ENABLED`、`GO_IMAGE`及`RUNTIME_IMAGE`公開為可覆寫的建置引數，以供 CGO 應用程式使用。

## 修正

- 在[PR](https://github.com/wailsapp/wails/pull/4435)中防止關閉仍有待處理非同步呼叫的視窗時發生當機；由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5249)中防止在 Windows 上開啟隱藏的應用程式時啟用視窗；由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5668)中確保 WebKit 要求中繼資料、回應完成處理與本文串流處理皆在 GTK 主執行緒上執行；由 @taliesin-ai 提交
- 修正`Menu.Update()`未在 GTK4 Linux 上重新建置原生選單的問題（#5659；@puneetdixit200 亦獨立診斷並於 #5539 中修正）
- 修正顯示器變更時列舉 macOS 螢幕所發生的當機，方法是複製螢幕 ID／名稱字串並建立數量快照（#5565；@x-haose 亦獨立診斷並於 #5584 中修正）
- 修正 Windows 上將視窗拖曳穿越不同 DPI 的顯示器後，WebView2 內容先縮小再消失的問題：在`WM_DPICHANGED`處理常式中重新設定控制器邊界，做法與取消最小化時的 DPI 重新同步一致（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 新增

- 透過`app.GlobalShortcut`新增全域（系統範圍）鍵盤快捷鍵（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）。即使應用程式未取得焦點，快捷鍵仍會觸發。各平台皆以原生方式實作，且不含第三方相依套件：macOS 使用 Carbon 快速鍵、Windows 使用`RegisterHotKey`、X11 使用`XGrabKey`，Wayland 則使用 XDG Desktop Portal 全域快捷鍵介面。
- 新增內建 MCP 伺服器：當應用程式使用`mcp`標籤建置時，Model Context Protocol 伺服器會自動啟動，讓 LLM 代理程式能測試及控制執行中的 Wails 應用程式，包括控制視窗、檢查 DOM、評估 JavaScript、呼叫繫結方法、處理事件，以及模擬滑鼠／鍵盤輸入並以畫面上的動畫游標呈現。無須撰寫任何使用者程式碼：設定`WAILS_MCP=1`後，`wails3 build`／`wails3 dev`會自動加入`mcp`標籤。所有設定皆透過環境變數完成（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）。

## 已修正

- 修正`Menu.Update()`未在 GTK4 Linux 上重新建置原生選單的問題（#5659；@puneetdixit200 亦獨立診斷並於 #5539 中修正）
- 修正顯示器變更時列舉 macOS 螢幕所發生的當機，方法是複製螢幕 ID／名稱字串並建立數量快照（#5565；@x-haose 亦獨立診斷並於 #5584 中修正）

## v3.0.0-alpha2.107 - 2026-06-27

## 新增

- @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5613)中新增含側邊欄導覽的實驗性 Wake 文件

## v3.0.0-alpha2.106 - 2026-06-24

## 變更

- 將`webview2`升級至v1.0.27。
  - ci(webview2)：修正發行版本建置（交叉編譯 Windows 並補齊 go.sum）（#5671）\

  **完整差異：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5672)中，從 webview2 發行工作流程的交叉編譯程序移除 go vet
- @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5670)中，將自動變更日誌使用的 OpenRouter 模型更新為 google/gemini-2.5-flash-lite
- 將`webview2`升級至v1.0.26。

### 修正

- **從暫時性執行階段 COM 錯誤復原，而非結束程式**（#5658、#5580）。`Chromium.errorCallback`先前遇到*任何* COM 錯誤都會呼叫`os.Exit(1)`，因此啟動後一次可復原的短暫異常就會終止整個應用程式。執行階段路徑（`Resize`／`GetClientRect`、`Navigate`／`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）現在會記錄錯誤並復原。尤其是`MessageReceived`中的格式錯誤／不受信任 Web 訊息，現在會被捨棄，而不會導致處理程序終止。這解決了跨越不同 DPI 顯示器時發生的一類當機（#5544、#5650）。建立環境／控制器的路徑仍會造成致命錯誤。\

**完整差異：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 已修正

- @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5671)中修正 release-webview2 工作流程，使其能正確處理 go.sum 檔案
- @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5659)中，透過清除並重新建置原生選單來修正 Linux GTK4 選單更新

## v3.0.0-alpha2.105 - 2026-06-21

## 新增

- 新增`application.System`，供共用程式碼在執行階段偵測平台：`System.IsMobile()`（iOS／Android）、`System.IsDesktop()`（macOS／Windows／Linux）、`System.IsServer()`（`server`建置標籤），以及用於直接測試單一目標的`System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`。它可在每個目標上編譯，因此無須建置標籤即可依平台分支。`@wailsio/runtime`中也提供對應的前端輔助函式（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）
- 新增「使用其他前端框架」指南，說明如何將自己的 Vite 專案直接放入`frontend/`（涵蓋 Solid、Preact、Lit、SvelteKit、Qwik、Angular 等）
- `wails3 setup`精靈現在會檢查行動平台（iOS／Android）工具鏈，包括 Xcode 與 iOS Simulator 執行階段、JDK、Android SDK／NDK 和模擬器；適用時並提供一鍵安裝及可複製的 shell 設定修正
- 產生的專案會隨附`frontend/.npmrc`，其中設定7天的`minimum-release-age`，以降低使用剛發佈（可能已遭入侵）套件的風險（pnpm 與 bun 會遵循此設定；npm 則會忽略，且不會造成影響）

## 變更

- 以全新的霓虹山景主視覺重新設計所有內建入門範本（Web、iOS 和 Android）
- <strong>TypeScript 現已成為入門範本的預設語言，並使用不含語言後綴的範本名稱。</strong>使用`wails3 init`（不含`-t`）會建立 TypeScript 專案；`-t vanilla`、`-t react`、`-t vue`和`-t svelte`皆為 TypeScript，而 JavaScript 版本則位於`-t vanilla-js`、`-t react-js`、`-t vue-js`和`-t svelte-js`。內建範本會在`template.yaml`中透過`typescript:`宣告其語言；使用`-ts`後綴的社群範本仍可作為備援方式繼續運作
- 以霓虹「digital Wails」主題重新設計`wails3 setup`精靈（以山景為背景的霧面玻璃鮮明效果）

## 已修正

- 修正 Windows 上還原應用程式時發生的當機：該應用程式最小化的時間長到足以讓 WebView2 暫停，或讓其轉譯／GPU 處理程序被回收。最小化／還原時的 DPI 重新同步（#5544）現在只會在視窗 DPI 確實變更時操作 WebView2 控制器，避免在常見的相同 DPI 還原情況下，對已暫停的控制器進行致命的 COM 呼叫（#5605）
- 修正長時間執行且頻繁載入資產／媒體的 Linux 應用程式反覆發生原生`SIGABRT`／`SIGSEGV`當機的問題（通常發生於 GTK 主迴圈期間的`g_object_unref`內）。資產伺服器先前會從工作 goroutine 完成`WebKitURISchemeRequest`，因此在 GTK 主執行緒之外呼叫執行緒不安全的 WebKit2GTK 函式；現在完成作業（`webkit_uri_scheme_request_finish_with_response`／`finish_error`）會在主執行緒上執行。這補全了 #5566中的部分修正。GTK3 與 GTK4／WebKitGTK 的6.0建置皆受影響（#5631、#5557）
- 修正 Linux／GTK3 上`setupSignalHandlers`中間歇出現`fatal error: invalid pointer found on stack`的問題。作為訊號`user_data`傳遞的視窗 ID 儲存在 Go 的`unsafe.Pointer`區域變數中，因此垃圾回收器在複製堆疊期間掃描此（非指標）值時會中止。現在 Go 端會將該 ID 保持為整數型別（`uintptr_t`），亦即將 #4958對 GTK4 路徑採用的相同修正向後移植至舊版 GTK3 路徑（該修正將 C 訊號函式改為`uintptr_t`，以消除`-race`／checkptr 錯誤）（#5631）

## 移除

- 移除`react-swc`、`preact`、`lit`、`solid`、`qwik`和`sveltekit`入門範本（以及其`-ts`變體）。目前支援的內建範本為`vanilla`、`react`、`vue`和`svelte`；每個範本預設皆使用 TypeScript，另提供`-js` JavaScript 變體。任何其他框架仍可透過[自備前端](https://v3.wails.io/guides/dev/frontend-frameworks)或自訂範本使用

## v3.0.0-alpha2.104 - 2026-06-18

## 已修正

- 修正已綁定的 Go 服務方法傳回空字串時導致 iOS 當機（SIGABRT）的問題。iOS 資產回應寫入器原本以`buf != nil`而非主體長度來防護主體指標，因此長度為零的主體會導致`&buf[0]`發生 panic；現在改為依長度進行防護，與桌面端寫入器一致

## v3.0.0-alpha2.103 - 2026-06-15

## 變更

- 將 iOS 和 Android 原生功能移至平台管理器：改為透過`application.IOS.*`和`application.Android.*`呼叫（例如`application.IOS.Haptic("medium")`、`application.Android.Share(payload)`），不再使用舊的`application.IOS*`／`application.Android*`自由函式（#5602）
- 重新命名行動裝置橋接事件：跨平台事件現在使用`common:*`前綴（例如`common:haptic`、`common:location`），平台專屬事件則使用`ios:*`／`android:*`（例如`ios:backgroundTask`、`android:foregroundService`）；不再使用`native:*`前綴（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 新增

- 新增實驗性的`wails3 setup`精靈，用於互動式專案設定及相依性檢查
- 為`wails3 doctor`新增`--json`旗標，以提供機器可讀的輸出
- 在`wails3 doctor`命令中新增簽署狀態區段

## 修正

- 修正 Linux 上的 npm 偵測，除套件管理器外也會檢查 PATH

## v3.0.0-alpha.101 - 2026-06-13

## 新增

- iOS：原生訊息對話方塊（UIAlertController），以及開啟檔案／多個檔案／目錄的對話方塊（UIDocumentPickerViewController）；儲存對話方塊會傳回明確的錯誤
- iOS：透過 UIPasteboard 支援剪貼簿
- iOS：透過 UIScreen 取得實際螢幕指標（點、像素、縮放比例、安全區域工作範圍）
- iOS：裝置建置（`IOS_PLATFORM=device`）、程式碼簽署身分／佈建描述檔／權利設定支援、`.ipa`封裝，以及透過 devicectl 執行`deploy-device`
- iOS：可設定最低 iOS 版本（build/config.yml 中的`ios.minIOSVersion`）
- iOS：`wails3 doctor`會回報 macOS 上 Xcode 與 iOS SDK 的可用性
- iOS：系統事件——電池、網路、主題、螢幕鎖定及記憶體不足會呈現為`events.IOS.*`和平台中立的`events.Common.*`應用程式事件
- iOS：原生行動功能橋接（匯出的`application.IOS*`）——分享面板、開啟 URL、保持喚醒、手電筒、安全區域內距、亮度、應用程式資訊、方向鎖定、狀態列、生物辨識（Face ID／Touch ID）、本機通知及 Keychain 安全儲存空間
- iOS：感測器與硬體——觸覺回饋、單次地理定位、加速度計、距離感測器、文字轉語音、儲存空間資訊、電源／電池狀態、網路狀態、鍵盤內距及螢幕擷取偵測
- iOS：文件（IOS.md 及文件網站指南）
- Android：原生訊息對話方塊（AlertDialog），以及開啟檔案／多個檔案的對話方塊（Storage Access Framework，匯入為快取副本）；開啟目錄與儲存對話方塊會傳回明確的錯誤
- Android：透過 ClipboardManager 支援剪貼簿
- Android：透過 WindowMetrics／DisplayMetrics 取得實際螢幕指標（dp、像素、縮放比例、系統列工作範圍）
- Android：觸覺回饋（`Android.Haptics.Vibrate`）、裝置資訊（`Android.Device.Info`）及 Toast（`Android.Toast.Show`）執行階段方法
- Android：具型別的生命週期事件（`events.Android.*`，由 events.txt 產生），並將`ActivityCreated`對應至`Common.ApplicationStarted`
- Android：建置管線會產生可安裝的偵錯版與發行版 APK（`android:run`、`android:package`、`android:package:fat`）；發行版預設使用偵錯金鑰庫簽署，也可透過`ANDROID_KEYSTORE_*`環境變數使用正式金鑰庫
- Android：`wails3 doctor`會回報 Android SDK、NDK 及 JDK
- Android：系統事件——電池、網路、主題、螢幕鎖定及記憶體不足會呈現為`events.Android.*`和平台中立的`events.Common.*`應用程式事件
- Android：原生行動功能橋接（匯出的`application.Android*`）——分享、開啟 URL、保持喚醒、手電筒、安全區域內距、亮度、應用程式資訊、方向鎖定、狀態列、生物辨識（BiometricPrompt）、本機通知及 EncryptedSharedPreferences 安全儲存空間
- Android：感測器與硬體——觸覺回饋、單次地理定位、加速度計、距離感測器、文字轉語音、儲存空間資訊、電源／電池狀態、網路狀態、鍵盤內距及透過 FLAG_SECURE 阻擋螢幕擷取
- Android：文件（ANDROID.md 及文件網站指南）
- 範例：`mobile`綜合範例新增 Mobile 和 Hardware 分頁，展示跨 iOS 與 Android 的原生功能橋接（膠囊式分頁會換行成多列）
- 行動裝置：電池——應用程式進入背景時，會暫停加速度計、距離感測器、手電筒及範例中的週期性時鐘，返回前景時再恢復（Android 會讓處理程序繼續在背景執行，而手電筒在 iOS 上是會持續保留的硬體狀態）；此外，Android 系統事件接收器只會在應用程式位於前景時註冊
- iOS：相機擷取——`application.IOSCapturePhoto`／`IOSCaptureVideo`（UIImagePickerController → 帶有 base64 縮圖的`native:capture`事件）
- iOS：背景執行——`application.IOSBeginBackgroundTask`／`IOSEndBackgroundTask`（UIApplication 背景工作時段），以及可設定的`ios.backgroundModes`（build/config.yml），用來將`UIBackgroundModes`套入產生的 Info.plist
- Android：相機擷取——`application.AndroidCapturePhoto`／`AndroidCaptureVideo`（透過 FileProvider 使用系統相機 → `native:capture`事件）
- Android：前景服務——`application.AndroidStartForegroundService`／`AndroidStopForegroundService`（帶有持續通知的`WailsForegroundService`會讓處理程序保持運作，以執行長時間的背景工作）
- 範例：新增 Camera 分頁，展示相片／影片擷取及背景執行（Android 使用前景服務，iOS 使用背景工作時段）

## 修正

- 修正`getUserMedia`在 Linux 上總是因`NotAllowedError`而失敗的問題：WebKitGTK 會拒絕無人處理的權限要求，而`permission-request`訊號先前並未連接。現在會依照新的跨平台`WebviewWindowOptions.Permissions`對應表（`map[PermissionType]Permission`）處理相機／麥克風權限，Linux（WebKitGTK）與 Windows（WebView2）皆會遵循此設定。Linux 沒有原生提示，因此相機／麥克風預設為允許（恢復`getUserMedia`），並可使用`PermissionDeny`停用（#5552）
- iOS：`GOOS=ios`現在可再次編譯（匯出的`events.IOS`、行動裝置方法名稱存根），帶有 production 標籤的建置也可編譯（修正 pkg/application 及數個服務中的建置標籤）
- iOS：Go→JS 事件和 ExecJS 現在可正常運作——頁面啟動時不再載入兩次，且`wails:runtime:ready`交握也不再可能遺失
- iOS：`ApplicationDidFinishLaunching`／`ApplicationStarted`不再與應用程式啟動產生競爭條件；已移除固定2秒的啟動等待
- iOS：修正每次執行 Go→JS JavaScript 時發生的 C 字串洩漏
- iOS：`hasListeners`現在會反映實際的監聽器註冊狀態
- iOS：正式環境建置會在編譯時排除框架偵錯記錄
- Android：`GOOS=android`現在可再次編譯——已定義`events.Android`、移除超出界限的`events_android.go`監聽器陣列、新增行動裝置方法名稱存根，並防止桌面 Linux 檔案（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）混入 Android 建置
- Android：JS→Go 綁定現在可正常運作——WebView 無法將`fetch()`POST 主體傳遞至`shouldInterceptRequest`，因此執行階段呼叫改經由 JavascriptInterface 傳輸（`nativeHandleRuntimeCall`），不再因要求主體為 nil 而當機
- Android：`Screens.*`執行階段呼叫現在會傳回實際資料——ScreenManager 現在會在啟動時填入資料（先前從未接線，因此`GetAll`會傳回 nil）
- Android：正式環境組建不會編譯框架偵錯記錄；偵錯組建則會以`Wails`標籤將記錄傳送至 logcat
- Android：實作真正的`hasListeners`登錄檔、JNI 參照／例外處理，以及僅載入一次的頁面生命週期（不會重複導覽）
- 修正在 Vite 開發伺服器執行時，`wails3 generate bindings`在 Windows 上因「存取遭拒」而失敗的問題：改為將產生的檔案同步至輸出目錄，而非透過重新命名覆蓋該目錄（#5515）
- 修正在顯示器變更後讀取螢幕資訊時，macOS 偶發嚴重當機的問題：螢幕 ID 與名稱所儲存的指標指向自動釋放的`UTF8String`緩衝區，而這些緩衝區可能在 Go 複製前即遭釋放（釋放後使用）。現在會對字串執行`strdup`，並在轉換後釋放；螢幕列舉也會在明確的自動釋放集區中執行，因此從 Go goroutine 呼叫時不再造成洩漏（#5556）
- 修正 assetserver 關閉`WebKitURISchemeRequest`時，Linux 上偶發 SIGSEGV 的問題：最後一次`g_object_unref`在 assetserver goroutine 上執行，導致 WebKit GObject 在 GTK 主執行緒之外完成終結。現在會透過`g_main_context_invoke`將 unref 封送至 GTK 主內容（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 新增

- 擴充`MacWebviewPreferences`，加入其他 WKWebView 組態選項：`EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize`及`ApplicationNameForUserAgent`（#5549）

## 修正

- 修正在 Vite 開發伺服器執行時，`wails3 generate bindings`在 Windows 上因「存取遭拒」而失敗的問題：改為將產生的檔案同步至輸出目錄，而非透過重新命名覆蓋該目錄（#5561）
- 修正 Linux 上無框視窗不會觸發 JS 調整大小事件的問題；並修正無框視窗的捲軸邊緣偵測（#5368）
- 修正當暫存目錄與安裝目錄位於不同磁碟區時，Windows 上的更新程式因「無效的跨裝置連結」而失敗的問題（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 修正

- 修正在 Vite 開發伺服器執行時，`wails3 generate bindings`在 Windows 上因「存取遭拒」而失敗的問題：改為將產生的檔案同步至輸出目錄，而非透過重新命名覆蓋該目錄（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 修正

- 修正 Linux 上 WebKit 閒置時（例如開啟檢查器時）UI 凍結的問題：不再對`SIGUSR1`強制套用`SA_ONSTACK`，因為這會破壞 JavaScriptCore 的 GC 執行緒同步（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 新增

- 新增偵錯頁面及使用`runtime/trace`的內容

## 變更

- 移除一些不必要的`_ "embed"`匯入，稍微整理程式碼

## 修正

- 修正 Windows 上視窗取消最大化後未強制套用最小寬度／高度限制的問題（#4593）
- 修正在 Frameless + Transparent 視窗選項下，全螢幕模式會發生滑鼠點擊穿透的問題（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 新增

- 新增 Garble 混淆支援（[#4563](https://github.com/wailsapp/wails/issues/4563)）：穩定的繫結方法 ID、組建／Taskfile 串接（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`），以及所有面向執行階段之承載資料上的 JSON 結構標籤（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`），使線上格式不受 Garble 重新命名匯出欄位的影響。

## v3.0.0-alpha.95 - 2026-05-20

## 新增

- 新增缺少的專案結構頁面

## 變更

- 文件：將架構頁面上的幾張圖改為循序圖，以呈現得更清楚
- 文件：加入執行前必須先安裝 D2 的說明

## 修正

- 修正 GTK4 預設組態上的`wails3 generate appimage`：打包工具現在會先從二進位檔偵測 GTK 技術堆疊，再搜尋執行階段檔案，因此 GTK4 組建會選取`libwebkitgtkinjectedbundle.so`（位於`webkitgtk-6.0/`下），`-tags gtk3`組建則會選取`libwebkit2gtkinjectedbundle.so`（位於`webkit2gtk-4.1/`下）。`.relr.dyn`探測也會檢查`libgtk-4.so.1`，因此無論使用哪種技術堆疊，現代工具鏈都能正確停用移除符號。（#5475）
- 修正以相對`-builddir`叫用`wails3 generate appimage`時失敗的問題：打包工具現在會預先將`-binary`、`-icon`、`-desktopfile`、`-builddir`及`-outputdir`解析為絕對路徑，因此流程中途的`s.CD`不會再破壞 AppRun 下載 goroutine 或複製後的`ldd`探測。
- 修正當桌面檔案的`Name=`欄位與二進位檔基本名稱不符時，`wails3 generate appimage`無法將最終 AppImage 移至`-outputdir`的問題：打包工具現在會透過`OUTPUT`環境變數，強制 linuxdeploy 的 appimage 外掛程式將 AppImage 寫入`<binary>-<arch>.AppImage`，而非使用衍生自桌面檔案的名稱。
- 修正 GTK4 + WebKitGTK 6.0技術堆疊於 alpha.93升為預設值後，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`及`Common.SystemDidWake`無法在 Linux 上觸發的問題。新的預設`application_linux.go` `run()`未呼叫`setupCommonEvents()`（此函式會將`Linux.*`事件轉送至其對應的`Common.*`事件）或`monitorPowerEvents()`。現在 GTK3 與 GTK4 組建路徑會透過`application_linux_dbus.go`共用 DBus 電源監控輔助程式。（#5474）

## v3.0.0-alpha.94 - 2026-05-19

## 修正

- 修正 GTK4 + WebKitGTK 6.0技術堆疊於 alpha.93升為預設值後，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`及`Common.SystemDidWake`無法在 Linux 上觸發的問題。新的預設`application_linux.go` `run()`未呼叫`setupCommonEvents()`（此函式會將`Linux.*`事件轉送至其對應的`Common.*`事件）或`monitorPowerEvents()`。現在 GTK3 與 GTK4 組建路徑會透過`application_linux_dbus.go`共用 DBus 電源監控輔助程式。（#5474）

## v3.0.0-alpha.93 - 2026-05-17

## 新增

- 由 @leaanthony 在 Linux 上將`XDG_SESSION_TYPE`加入`wails3 doctor`輸出

## 修正

- 由 @leaanthony 修正 appmenu-gtk-module 存取尚未實現的視窗，導致 Wayland 上視窗選單當機的問題（#4769）
- 由 @leaanthony 修正應用程式名稱包含無效字元（空格、括號等）時，GTK 應用程式當機的問題
- 由 @overlordtm 修正 Windows 上初始化拖放時發生「記憶體不足」錯誤的問題（#4701）
- 由 @leaanthony 修正 mainthread 回呼儲存區在刪除 map 項目時錯誤使用 RLock 所造成的競爭條件（Linux、macOS、iOS）（#4424）
- 修正將命令列引數傳遞至工作時的變數處理。現在會正確初始化以 KEY=VALUE 配對指定的 CLI 變數，並在整個工作執行期間傳播這些變數。
- 修正 macOS 上的 NSWindowZoomButton 衝突：`MaximiseButtonState`與`FullscreenButtonState`現在會在啟動和執行階段套用限制較嚴格的狀態；任一設定函式都無法再悄然覆寫另一個設定（#5319）
- 修正舊版 GTK3 建置路徑（`-tags gtk3`）中由 CodeRabbit 在 #5463 發現的一系列既有錯誤：透過檔案關聯啟動時不再略過啟動處理常式；`getTheme` 現已具備邊界與型別安全性；`appName` 不再釋放由 GLib 擁有的記憶體；`clipboardGet` 不再洩漏 GTK 傳回的 `gchar*`；`Calloc` 現在使用指標接收者（且 `NewCalloc` 會傳回 `*Calloc`），因此集區能確實追蹤配置；`zoomOut` 改用 `zoomInFactor` 的倒數，不再使用會被限制為 1.0 的負乘數；`execJS` 會重複使用預先配置的空白 world-name，不再於每次呼叫時洩漏一個 `C.CString("")`；並已從 `menuItem.setAccelerator` 移除開發期間使用的 `fmt.Println`。解決 #5465。
- 修正預設 GTK4 建置路徑（`linux_cgo.go`）中相同的 `Calloc` 值接收者洩漏：改用指標接收者與 `NewCalloc() *Calloc`，使每個視窗的 `c.String(...)` 配置都能確實受到追蹤並釋放。

## v3.0.0-alpha.92 - 2026-05-15

## 新增

- 修改 Taskfile，允許透過 `PACKAGE_MANAGER` 選項控制所使用的前端套件管理器
- 在範本資料中加入 `{{.Opn}}` 與 `{{.Cls}}`，使 Taskfile 範本的撰寫更可預期

## 變更

- 修改部分現有的 Taskfile，以使用 `{{.Opn}} and {{.Cls}}`

## 修正

- 修正面板讀取系統匣選單時更新該選單，會在 `linuxSystemTray` 中引發 `concurrent map read and map write` 執行階段嚴重錯誤的問題。
- WebView2 的錯誤與堆疊追蹤輸出改用 `log`，而非 `fmt`，避免應用程式在 Windows 上未連接主控台執行時遺失訊息。

## v3.0.0-alpha.91 - 2026-05-12

## 變更

- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5414) 中更新贊助者 SVG
- <strong>破壞性變更（macOS）：</strong>統一 macOS 座標系統，使 `GetScreens`、`Position` 與 `SetPosition` 全部使用相同的空間：以邏輯點為單位、Y 軸向下，並以主螢幕左上角的 `(0,0)` 為原點。這與 Windows、GTK，以及 Electron 和 Web 的公開 API 一致。實體位置在主螢幕上方的螢幕現在會回報負的 `Bounds.Y`（先前為正值），而 `Position()`/`SetPosition()` 值現在以邏輯點為單位，不再使用 `points × primaryScale`。仍可保證 `Position()` → `SetPosition()` 的往返轉換；若使用較早 alpha 版本記錄的絕對值，或手動計算的因應方式（例如乘以 `primaryScale`，或依螢幕高度反轉 Y 軸），則需要更新。解決 [#5117](https://github.com/wailsapp/wails/issues/5117)。

## 修正

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5416) 中加入防禦性驗證 DBus 訊號名稱與主體長度，避免發生 panic
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5363) 中修正 Linux 上 GTK 選單處理的記憶體安全性問題
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5295) 中偵測 NVIDIA GPU，並在 Linux 上停用 DMA-BUF 轉譯器
- 修正 macOS 上 `SetPosition` 的跨螢幕 Y 座標轉換：使用主螢幕高度作為全域基準，使視窗在垂直位置偏離主顯示器的螢幕上也能落在正確位置，詳見 [#5117](https://github.com/wailsapp/wails/issues/5117)
- 由 @wayneforrest 在 [PR](https://github.com/wailsapp/wails/pull/5109) 中修正 Git PR 範本，使其指向正確的意見回饋 URL
- 修正一系列 Windows 系統匣 `SetMenu` 當機問題。問題源自已損壞的 `DestroyMenu` 系統呼叫傳入四個引數而非一個，導致每次呼叫都傳回 FALSE，且未釋放任何資源。另外，重建選單時也會釋放 HMENU 與 HBITMAP 控制代碼（包括執行階段透過 `MenuItem.SetBitmap` 配置的控制代碼）、在 `Win32Menu.Update` 中重設過期的核取方塊／選項按鈕對應表，並移除 `systemtray.updateMenu` 中造成配置量加倍的多餘 `Update()` 呼叫。長時間執行的系統匣應用程式不再於每次重建選單時洩漏 GDI/USER 物件。

## v3.0.0-alpha.90 - 2026-05-11

## 新增

- 由 @vinhvoit225 在 [PR](https://github.com/wailsapp/wails/pull/5261) 中新增可設定的應用程式名稱，供 macOS 上的 WKWebView User-Agent 使用
- 由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5400) 中為 gin-service 範例新增間接相依套件 github.com/coder/websocket
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5402) 中為建置資產測試新增深度相等比較支援

## 變更

- 由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5401) 中將建置輸出整合至 assets 目錄
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5399) 中更新贊助者 SVG

## 修正

- 由 @overlordtm 在 [PR](https://github.com/wailsapp/wails/pull/5289) 中改用通知物件傳送 macOS 單一執行個體訊息
- 由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5383) 中批次處理 Windows 回呼，避免高負載下遺失 promise

## v3.0.0-alpha.89 - 2026-05-10

## 新增

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5316) 中新增 go<em>test</em>results 工作，以彙總 Go 測試結果

## 變更

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5369) 中依條件將大型 RPC 承載資料拆分為分塊 POST 請求
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5386) 中將所有前端範本的 Vite 從 5.x.x 升級至 8.0.0
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5365) 中將 Vite 開發伺服器的連接埠設定遷移至環境變數
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5361) 中設定所有範本的 Vite 開發伺服器繫結至 127.0.0.1
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5384) 中更新贊助者 SVG

## 修正

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5312) 中於 build-assets 更新期間清理 Info.plist 範本存根
- 在主執行緒同步套用選單項目修改方法（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`），修正 macOS 選單的過期狀態，並消除快速重新開啟選單時仍轉譯先前狀態的 `dispatch_async` 競爭條件（#5002）
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5203) 中讓開發模式忽略 `*_test.go` 檔案，以避免不必要的重新建置
- 由 @wucm667 在 [PR](https://github.com/wailsapp/wails/pull/5291) 中防止應用程式未執行時 Menu.Update() 發生記憶體區段錯誤
- 由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5382) 中使用 lastSizeWParam 控制 Windows 上選單列的重新繪製

## v3.0.0-alpha.88 - 2026-05-09

## 變更

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5371) 中將 HiddenOnTaskbar 改為使用 WS<em>EX</em>TOOLWINDOW
- 由 @atterpac 在 [PR](https://github.com/wailsapp/wails/pull/5370) 中重新排列相依套件，並移除 go.mod 中的 webview2 replace 指令
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5358) 中更新贊助者 SVG

## 修正

- 由 @fbbdev 在 [PR](https://github.com/wailsapp/wails/pull/5331) 中移除泛型間接別名，並整合對應表鍵值型別

## 移除

- 刪除 PR-master 工作流程，移除文件、Go 測試，以及 [PR](https://github.com/wailsapp/wails/pull/5377) 中由 @leaanthony 加入的略過測試

## v3.0.0-alpha.87 - 2026-05-07

## 新增

- 在 [PR](https://github.com/wailsapp/wails/pull/5352) 中由 @leaanthony 新增 Wails v3 的韓文文件
- 在 [PR](https://github.com/wailsapp/wails/pull/5354) 中由 @leaanthony 新增安裝與快速入門的法文文件
- 在 [PR](https://github.com/wailsapp/wails/pull/5355) 中由 @leaanthony 新增快速入門、概念及社群的葡萄牙文文件

## v3.0.0-alpha.86 - 2026-05-06

## 新增

- 在 [PR](https://github.com/wailsapp/wails/pull/5328) 中由 @leaanthony 新增法文文件本地化
- 在 [PR](https://github.com/wailsapp/wails/pull/5343) 中由 @leaanthony 將德文語系新增至文件網站

## 變更

- 在 [PR](https://github.com/wailsapp/wails/pull/5347) 中由 @leaanthony 將全部 8 個翻譯語系註冊至文件設定
- 在 [PR](https://github.com/wailsapp/wails/pull/5317) 中由 @leaanthony 更新多個與 WebView2 相關的 Windows 檔案

## 修正

- 在 [PR](https://github.com/wailsapp/wails/pull/5340) 中由 @leaanthony 將 Linux 的對話方塊分派拆分為 GTK3 與 GTK4 兩套處理
- 確保對話方塊回呼在 GTK 執行緒上執行，並在 [PR](https://github.com/wailsapp/wails/pull/5339) 中由 @leaanthony 修正記憶體區段錯誤

## v3.0.0-alpha.85 - 2026-05-05

## 新增

- 在 [PR](https://github.com/wailsapp/wails/pull/5179) 中由 @leaanthony 將 PR 範本 URL 新增至儲存庫
- 在 [PR](https://github.com/wailsapp/wails/pull/5330) 中由 @leaanthony 新增 Wails v3 的德文文件

## v3.0.0-alpha.84 - 2026-05-03

## 新增

- 在 [PR](https://github.com/wailsapp/wails/pull/5307) 中由 @leaanthony 新增可停用 macOS 上以 Escape 鍵退出全螢幕的選項
- 在 [PR](https://github.com/wailsapp/wails/pull/5310) 中由 @leaanthony 新增可停用 macOS 上以 Escape 鍵退出全螢幕的選項
- 在 [PR](https://github.com/wailsapp/wails/pull/5288) 中由 @yuseferi 新增 Pausa 社群展示文件

## 變更

- 在 [PR](https://github.com/wailsapp/wails/pull/5308) 中由 `@github-actions[bot]` 更新贊助者 SVG
- 在 [PR](https://github.com/wailsapp/wails/pull/5309) 中由 @leaanthony 更新圖示產生命令，以處理不支援的平台
- 以三態 ButtonState 取代布林值全螢幕 API，並在 [PR](https://github.com/wailsapp/wails/pull/5224) 中由 @leaanthony 實作平台繫結

## 修正

- 在 [PR](https://github.com/wailsapp/wails/pull/5315) 中由 @leaanthony 防止 WebView2 在控制器狀態為 nil 時執行焦點操作
- 在 [PR](https://github.com/wailsapp/wails/pull/5313) 中由 @leaanthony 更新 GitHub Actions 工作流程，使其正確參照 PR 的基底分支
- 在開發模式中忽略 `*_test.go` 檔案，以避免不必要的重新建置；此變更由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5203) 中完成
- 在 [PR](https://github.com/wailsapp/wails/pull/5291) 中由 @wucm667 防止應用程式未執行時 Menu.Update() 發生記憶體區段錯誤

## v3.0.0-alpha.83 - 2026-05-02

## 新增

- 在 [PR](https://github.com/wailsapp/wails/pull/5094) 中由 @symball 新增 InstallScope 旗標，以及機器／使用者安裝的建置選項
- 在 [PR](https://github.com/wailsapp/wails/pull/5294) 中由 @leaanthony 為 BrowserWindow 新增無操作的 SetScreen 方法，以滿足 Window 介面

## 修正

- 在 [PR](https://github.com/wailsapp/wails/pull/5295) 中由 @leaanthony 偵測 NVIDIA GPU，並在 Linux 上停用 DMA-BUF 算繪器
- 在 [PR](https://github.com/wailsapp/wails/pull/5109) 中由 @wayneforrest 修正 Git PR 範本，使其指向正確的意見回饋 URL
- 修正一系列 Windows 系統匣 `SetMenu` 當機問題。這些問題源於損壞的 `DestroyMenu` 系統呼叫傳入四個引數而非一個，導致每次呼叫都傳回 FALSE，且未釋放任何資源。此外，在選單重建時釋放 HMENU 與 HBITMAP 控制代碼（包括執行階段透過 `MenuItem.SetBitmap` 配置的控制代碼）、重設 `Win32Menu.Update` 中過期的核取方塊／選項按鈕對應，並移除 `systemtray.updateMenu` 中會使配置量加倍的多餘 `Update()` 呼叫。長時間執行的系統匣應用程式不再於每次選單重建時洩漏 GDI／USER 物件。

## v3.0.0-alpha.82 - 2026-05-01

## 修正

- 在 [PR](https://github.com/wailsapp/wails/pull/5232) 中由 @leaanthony 修正桌面檔案的產生方式，使其正確處理桌面名稱

## v3.0.0-alpha.81 - 2026-04-30

## 變更

- 在 [PR](https://github.com/wailsapp/wails/pull/5286) 中由 @leaanthony 將每日建置版發行排程調整為 UTC 15:00

## 修正

- 修正 Retina Mac 上 Screen Bounds、WorkArea 與 Size 的值減半問題 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 變更

- 在 [PR](https://github.com/wailsapp/wails/pull/5285) 中由 @leaanthony 更新文件相依套件與內容集合載入器

## v3.0.0-alpha.79 - 2026-04-29

## 新增

- 在 [PR](https://github.com/wailsapp/wails/pull/5270) 中由 @leaanthony 授予 trigger-release 作業 actions: write 權限

## 變更

- 在 [PR](https://github.com/wailsapp/wails/pull/5283) 中由 @leaanthony 將發行工作預設為 master 分支，並更新變更記錄措辭
- 在 [PR](https://github.com/wailsapp/wails/pull/5282) 中由 @leaanthony 更新 auto-changelog 工作流程以使用最新版本
- 在 [PR](https://github.com/wailsapp/wails/pull/5280) 中由 @leaanthony 新增路徑篩選條件並移除已失效的工作流程，以提升工作流程效率
- 在 [PR](https://github.com/wailsapp/wails/pull/5274) 中由 @leaanthony 更新文件，使範例連結參照 master 分支
- 在 [PR](https://github.com/wailsapp/wails/pull/5272) 中由 @leaanthony 更新 v3 的文件與範例

## 修正

- 在 [PR](https://github.com/wailsapp/wails/pull/5265) 中由 @AkagiYui 為開發環境的反向代理新增重試邏輯並強制使用 IPv4
- 在 [PR](https://github.com/wailsapp/wails/pull/5281) 中由 @leaanthony 重寫未發行變更記錄的觸發工作流程

## 移除

- 在 [PR](https://github.com/wailsapp/wails/pull/5267) 中由 @leaanthony 移除用於各種測試用途的 Shell 測試指令碼
- 在 [PR](https://github.com/wailsapp/wails/pull/5266) 中由 @leaanthony 刪除 v3-alpha 文件部署工作流程與 CNAME 記錄

### 新增

- 在 [PR](https://github.com/wailsapp/wails/pull/5196) 中由 @leaanthony 將「前端路由」項目新增至側邊欄導覽
- 在 [PR](https://github.com/wailsapp/wails/pull/5185) 中由 @leaanthony 新增含各框架專屬建議的前端路由指南
- 新增對模態工作表（macOS）的支援
- 由 @leaanthony 升級 ghw 版本，以改善對 Apple 裝置的支援（#4977）
- 在 Dock 服務中新增 `GetBadge` 方法
- 在 `wails3 build` 命令中新增 `-tags` 旗標，以傳遞自訂 Go 建置標籤（例如 `wails3 build -tags gtk4`）（#4957）
- 新增繫結產生器自動產生列舉的文件，包括專用的 Enums 頁面與側邊欄導覽（#4972）
- 在 `wails3 build` 命令中新增 `-tags` 旗標，以傳遞自訂 Go 建置標籤（例如 `wails3 build -tags gtk4`）（#4957）
- 在 `v3/examples/web-apis/` 中新增 Web API 範例，示範 41 個以上的瀏覽器 API，包括儲存（localStorage、sessionStorage、IndexedDB、Cache API）、網路（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、媒體（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、裝置（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、效能（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）等
- 新增 WebView API 相容性檢查工具範例（`v3/examples/webview-api-check/`），用於跨平台測試 200 個以上的瀏覽器 API
- 新增 `internal/libpath` 套件，用於在 Linux 上尋找原生程式庫路徑，支援平行搜尋、快取及 Flatpak/Snap/Nix
- <strong>開發中：</strong>新增實驗性的 Linux WebKitGTK 6.0 / GTK4 支援，可透過 `-tags gtk4` 使用（GTK3/WebKit2GTK 4.1 仍為預設值）
- 注意：在平鋪式視窗管理員（例如 Hyprland、Sway）上，由於視窗幾何由視窗管理員控制，最小化／最大化操作可能不會如預期運作
- 由 @AbdelhadiSeddar 在 **在 JavaScript 中監聽事件**文件中新增如何使用**單次處理常式**的說明
- 由 @leaanthony 在 `WebviewWindowOptions` 中新增 `UseApplicationMenu` 選項，讓 Windows/Linux 上的視窗可繼承透過 `app.Menu.Set()` 設定的應用程式選單
- 由 @wimaha 新增使用 `.icon` 檔案（Apple Icon Composer 格式）產生 Liquid Glass 圖示與資產目錄（macOS）的支援（#4934）
- 新增適用於無頭／Web 部署的實驗性伺服器模式（`-tags server`）。此模式可讓 Wails 應用程式在不依賴原生 GUI 的情況下作為 HTTP 伺服器執行。請使用 `wails3 task build:server` 建置。詳情請參閱 `examples/server`。
- 新增 `internal/libpath` 套件，用於在 Linux 上尋找原生程式庫路徑，支援平行搜尋、快取及 Flatpak/Snap/Nix
- 由 @leaanthony 在 `MacWindow` 中新增 `CollectionBehavior` 選項，用於控制視窗在 macOS Spaces 與全螢幕模式下的行為（#4756）
- 由 @leaanthony 為 pkg/application 新增單元測試
- 由 @leaanthony 為 MSIX 封裝新增自訂通訊協定支援
- 新增 Linux 桌面環境偵測 [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- 由 @leaanthony 在 JavaScript 執行階段新增 `Window.Print()` 方法，以便從前端觸發列印對話方塊（#4290）
- 由 @leaanthony 在 Linux 的 `wails3 doctor` 輸出中新增 `XDG_SESSION_TYPE`
- 由 @leaanthony 為 Linux 新增 WebKit2 載入狀態變更事件：`WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished`（#3896）
- 由 @leaanthony 在 Linux 的 `wails3 doctor` 輸出中新增 `XDG_SESSION_TYPE`
- 在 Linux 建置期間產生 `.desktop` 檔案，而不只是在封裝期間產生（#4575）
- 由 @leaanthony 新增 Linux 執行階段相依性文件，包含各發行版特定的套件名稱與 nfpm 封裝範例（#4339）
- 由 @leaanthony 在 Linux 的 `wails3 doctor` 輸出中新增 NVIDIA 驅動程式版本資訊
- 由 @APshenkin 在 [PR](https://github.com/wailsapp/wails/pull/4710) 中為原始訊息處理常式新增來源
- 由 @APshenkin 在 [PR](https://github.com/wailsapp/wails/pull/4712) 中新增 macOS 通用連結支援
- 由 @APshenkin 在 [PR](https://github.com/wailsapp/wails/pull/4702) 中重構繫結傳輸層
- 由 @chinenual 在 [PR](https://github.com/wailsapp/wails/pull/4760) 中為 helloworld 範本新增 aria-label 識別碼，讓 Appium 測試用戶端能輕鬆測試範例應用程式
- 由 @APshenkin 在 [PR](https://github.com/wailsapp/wails/pull/4710) 中為原始訊息處理常式新增來源
- 由 @APshenkin 在 [PR](https://github.com/wailsapp/wails/pull/4712) 中新增 macOS 通用連結支援
- 由 @APshenkin 在 [PR](https://github.com/wailsapp/wails/pull/4702) 中重構繫結傳輸層
- 由 @fbbdev 與 @ianvs 在 [#4633](https://github.com/wailsapp/wails/pull/4633) 中新增型別化事件
- 新增 `systray-clock` 範例，示範可即時更新工具提示的無頭系統匣（#4653）。
- 由 @Tolfx 在 #4510 中新增適用於 Windows 的 NSIS Protocol 範本
- 由 @Tolfx 在 #4510 中新增 build-assets 測試
- macOS：由 @nidib 在 [#4588](https://github.com/wailsapp/wails/pull/4588) 中讓選單列顯示原生視窗控制項
- 由 @popaprozac 在 [PR](https://github.com/wailsapp/wails/pull/4451) 中新增 macOS Dock 服務，以隱藏／顯示 Dock 中的應用程式圖示
- 由 @popaprozac 在 [PR](https://github.com/wailsapp/wails/pull/4451) 中新增 macOS Dock 服務，以隱藏／顯示 Dock 中的應用程式圖示
- 由 @leaanthony 在 [#4534](https://github.com/wailsapp/wails/pull/4534) 中新增 macOS 原生 Liquid Glass 效果支援，使用 NSGlassEffectView（macOS 15.0+）並以 NSVisualEffectView 作為後備方案，且包含完整的材質自訂選項
- 由 @leaanthony 在 [#4500](https://github.dev/wailsapp/wails/pull/4500) 中新增瀏覽器 URL 清理。以 @APShenkin 的 [#4484](https://github.com/wailsapp/wails/pull/4484) 為基礎。
- 由 [@leaanthony](https://github.com/leaanthony) 在此 [PR](https://github.com/wailsapp/wails/pull/4241) 中新增 Windows/Mac 內容保護，並以 [@Taiterbase](https://github.com/Taiterbase) 的原始工作為基礎
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/4488) 中新增支援，可透過 `wails3 build` 與 `wails3 package` 別名將 CLI 變數傳遞給 Task 命令（#4422）
- 由 [@atterpac](https://github.com/atterpac) 在 [#4318](https://github.com/wailsapp/wails/pull/4318) 中新增置放區支援，其事件來源會包含被放置元素的資料
- 在 [PR](https://github.com/wailsapp/wails/pull/4467) 中，於 `WindowsWindow` 選項新增 `AdditionalLaunchArgs`，以便將額外的命令列引數傳遞給 WebView2 瀏覽器。
- 由 [@triadmoko](https://github.com/triadmoko) 在 [PR](https://github.com/wailsapp/wails/pull/4286) 中新增在 wails init 後自動執行 go mod tidy 的功能
- 由 @leaanthony 在 [PR](https://github.dev/wailsapp/wails/pull/4463) 中新增 Windows Snap Assist 功能
- 在 [PR](https://github.com/wailsapp/wails/pull/4467) 中，於 `WindowsWindow` 選項新增 `AdditionalLaunchArgs`，以便將額外的命令列引數傳遞給 WebView2 瀏覽器。
- 由 [@triadmoko](https://github.com/triadmoko) 在 [PR](https://github.com/wailsapp/wails/pull/4286) 中新增在 wails init 後自動執行 go mod tidy 的功能
- 由 @leaanthony 在 [PR](https://github.dev/wailsapp/wails/pull/4463) 中新增 Windows Snap Assist 功能
- 由 [@almas-x](https://github.com/almas-x) 在 [PR](https://github.com/wailsapp/wails/pull/4427) 中新增 Windows `getAccentColor` 實作
- 由 [@almas-x](https://github.com/almas-x) 在 [PR](https://github.com/wailsapp/wails/pull/4427) 中新增 Windows `getAccentColor` 實作
- Windows 深色主題選單與選單列。由 @leaanthony 於 [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5) 中實作
- 重新命名內建服務，讓 JS/TS 繫結更清楚，由 @popaprozac 於 [PR](https://github.com/wailsapp/wails/pull/4405) 中實作
- `app.Env.GetAccentColor`可取得使用者系統的強調色。支援 macOS。由 [@etesam913](https://github.com/etesam913) 實作
- 新增 `window.ToggleFrameless()` API，由 [@atterpac](https://github.com/atterpac) 於 [#4137](https://github.com/wailsapp/wails/pull/4137) 中實作
- 新增各 Linux 發行版專用的建置相依套件，由 @leaanthony 於 [PR](https://github.com/wailsapp/wails/pull/4345) 中實作
- 新增繫結指南，由 @atterpac 於 [PR](https://github.com/wailsapp/wails/pull/4404) 中實作
- **整理測試基礎架構**：將 Docker 測試檔案移至專用的 `test/docker/` 目錄，並採用最佳化映像檔及提升建置可靠性，由 [@leaanthony](https://github.com/leaanthony) 於 [#4359](https://github.com/wailsapp/wails/pull/4359) 中實作
- **改進資源管理模式**：在範例中加入正確的事件處理常式清理，以及可感知 context 的 goroutine 管理，由 [@leaanthony](https://github.com/leaanthony) 於 [#4359](https://github.com/wailsapp/wails/pull/4359) 中實作
- 支援建置 aarch64 AppImage，由 [@AkshayKalose](https://github.com/AkshayKalose) 於 [#3981](https://github.com/wailsapp/wails/pull/3981) 中實作
- 在 `wails doctor` 中新增診斷章節，由 [@leaanthony](https://github.com/leaanthony) 實作
- 呼叫服務方法時，將視窗加入 context，由 [@leaanthony](https://github.com/leaanthony) 實作
- 新增 `window-call` 範例，示範如何得知是哪個視窗正在呼叫服務，由 [@leaanthony](https://github.com/leaanthony) 實作
- 新增選單指南，由 [@leaanthony](https://github.com/leaanthony) 撰寫
- 改進 panic 處理，由 [@leaanthony](https://github.com/leaanthony) 實作
- 新增選單指南，由 [@leaanthony](https://github.com/leaanthony) 撰寫
- 為 Service API 新增文件註解，由 [@fbbdev](https://github.com/fbbdev) 於 [#4024](https://github.com/wailsapp/wails/pull/4024) 中實作
- 新增函式 `application.NewServiceWithOptions`，以使用額外設定初始化服務，由 [@leaanthony](https://github.com/leaanthony) 於 [#4024](https://github.com/wailsapp/wails/pull/4024) 中實作
- 改進選單控制，由 [@FalcoG](https://github.com/FalcoG) 與 [@leaanthony](https://github.com/leaanthony) 於 [#4031](https://github.com/wailsapp/wails/pull/4031) 中實作
- 新增更多文件，由 [@leaanthony](https://github.com/leaanthony) 撰寫
- 支援在標準事件接聽器中取消事件，由 [@leaanthony](https://github.com/leaanthony) 實作
- 支援系統匣的 `Hide`、`Show` 與 `Destroy`，由 [@leaanthony](https://github.com/leaanthony) 實作
- 支援系統匣的 `SetTooltip`，由 [@leaanthony](https://github.com/leaanthony) 實作。原始構想來自 [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- 在繫結產生器針對不支援型別的警告中回報套件路徑，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中實作
- 新增繫結產生器對泛型別名的支援，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中實作
- 新增繫結產生器對 `omitzero` JSON 旗標的支援，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中實作
- 新增 `//wails:ignore` 指示詞，以防止為選定的服務方法產生繫結，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中實作
- 在服務與模型上新增 `//wails:internal` 指示詞，允許使用在 Go 中匯出、但不在 JS/TS 中匯出的型別，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中實作
- 新增繫結產生器對別名型別常數的支援，以允許弱型別列舉，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 中實作
- 新增針對 Go 1.24 功能的繫結產生器測試，由 [@fbbdev](https://github.com/fbbdev) 於 [#4068](https://github.com/wailsapp/wails/pull/4068) 中實作
- 在 `OSInfo.Branding` 中新增對 macOS 15「Sequoia」的支援，以改進作業系統版本偵測，於 [#4065](https://github.com/wailsapp/wails/pull/4065) 中實作
- 新增 `PostShutdown` 掛鉤，以在關閉程序完成後執行自訂程式碼，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 中實作
- 新增 `FatalError` 結構，以支援在自訂錯誤處理常式中偵測嚴重錯誤，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 中實作
- 標準化服務啟動與關閉順序並加以記錄，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 中實作
- 新增應用程式啟動／關閉流程的測試框架，以及服務啟動／關閉測試，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 中實作
- 新增 `RegisterService` 方法，以在應用程式建立後註冊服務，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 中實作
- 在應用程式與服務選項中新增 `MarshalError` 欄位，以自訂繫結呼叫的錯誤處理，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 中實作
- 新增可取消的 Promise 包裝函式，可透過 Promise 鏈傳遞取消要求，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 中實作
- 新增將繫結呼叫的取消操作連結至 `AbortSignal` 的能力，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 中實作
- WML 除了支援一般的 `wml-*` 屬性外，也支援 `data-wml-*` 屬性，由 [@leaanthony](https://github.com/leaanthony) 實作
- 在所有服務上新增 `Configure` 方法，以支援延後設定與動態重新設定，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 中實作
- `fileserver` 服務在未設定時會傳送 503 Service Unavailable 回應，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 中實作
- `kvstore` 服務在未設定時預設提供記憶體內鍵值存放區，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 中實作
- 在 `kvstore` 服務上新增 `Load` 方法，以在設定變更後從檔案重新載入資料，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 中實作
- 在 `kvstore` 服務上新增 `Clear` 方法，以刪除所有鍵，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 中實作
- 在 `log` 服務中新增型別 `Level`，以提供 JS 端的日誌層級常數，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 中實作
- 在 `log` 服務上新增 `Log` 方法，以動態指定日誌層級，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 中實作
- `sqlite`服務在未設定時，預設會提供記憶體內資料庫；由[@fbbdev](https://github.com/fbbdev)於[#4067](https://github.com/wailsapp/wails/pull/4067)中實作
- 在`sqlite`服務中新增`Close`方法，以便手動關閉資料庫；由[@fbbdev](https://github.com/fbbdev)於[#4067](https://github.com/wailsapp/wails/pull/4067)中實作
- 為`sqlite`服務的查詢方法新增取消支援；由[@fbbdev](https://github.com/fbbdev)於[#4067](https://github.com/wailsapp/wails/pull/4067)中實作
- 為`sqlite`服務新增具備 JS 繫結的預備陳述式支援；由[@fbbdev](https://github.com/fbbdev)於[#4067](https://github.com/wailsapp/wails/pull/4067)中實作
- 新增 Gin 支援；由[Lea Anthony](https://github.com/leaanthony)於[PR](https://github.com/wailsapp/wails/pull/3537)中實作，並以[@AnalogJ](https://github.com/AnalogJ)在此[PR](https://github.com/wailsapp/wails/pull/3537)中的原始工作為基礎
- 修正自動儲存與密碼自動儲存一律啟用的問題；由[@oSethoum](https://github.com/osethoum)於[#4134](https://github.com/wailsapp/wails/pull/4134)中實作
- 在視窗上新增`SetMenu()`，以便為視窗設定選單；由[@leaanthony](https://github.com/leaanthony)實作
- 新增通知支援；由[@popaprozac](https://github.com/popaprozac)於[#4098](https://github.com/wailsapp/wails/pull/4098)中實作
-  新增 macOS 檔案關聯支援；由[@wimaha](https://github.com/wimaha)於[#4177](https://github.com/wailsapp/wails/pull/4177)中實作
- 新增`wails3 tool version`以遞增語意化版本；由[@leaanthony](https://github.com/leaanthony)實作
- 新增 macOS 與 Windows 的徽章支援；由[@popaprozac](https://github.com/popaprozac)於[#](https://github.com/wailsapp/wails/pull/4234)中實作
- 新增已註冊／強型別事件支援；由[@fbbdev](https://github.com/fbbdev)與[@IanVS](https://github.com/IanVS)於[#4161](https://github.com/wailsapp/wails/pull/4161)中實作
- 新增為自訂事件註冊掛鉤的功能；由[@fbbdev](https://github.com/fbbdev)與[@IanVS](https://github.com/IanVS)於[#4161](https://github.com/wailsapp/wails/pull/4161)中實作
- 新增`app.OpenFileManager(path string, selectFile bool)`，用於在系統檔案管理員中開啟路徑`path`，並可選擇透過`selectFile`醒目顯示；由[@Krzysztofz01](https://github.com/Krzysztofz01)與[@rcalixte](https://github.com/rcalixte)實作
- 為`wails3 init`命令新增`-git`旗標；由[@leaanthony](https://github.com/leaanthony)實作
- 新增`wails3 generate webview2bootstrapper`命令；由[@leaanthony](https://github.com/leaanthony)實作
- 在執行階段中新增`init()`方法，以允許手動初始化執行階段；由[@leaanthony](https://github.com/leaanthony)實作
- 在 Windows 的 WindowOptions 中新增`WindowDidMoveDebounceMS`選項；由[@leaanthony](https://github.com/leaanthony)實作
- 新增單一執行個體功能；由[@leaanthony](https://github.com/leaanthony)實作。此功能以 @APshenkin 的[v2 PR](https://github.com/wailsapp/wails/pull/2951)為基礎。
- `wails3 generate template`命令；由[@leaanthony](https://github.com/leaanthony)實作
- `wails3 releasenotes`命令；由[@leaanthony](https://github.com/leaanthony)實作
- `wails3 update cli`命令；由[@leaanthony](https://github.com/leaanthony)實作
- `wails3 generate bindings`命令的`-clean`選項；由[@leaanthony](https://github.com/leaanthony)實作
- 允許建置 aarch64 (arm64) AppImage Linux 版本；由[@AkshayKalose](https://github.com/AkshayKalose)於[#3981](https://github.com/wailsapp/wails/pull/3981)中實作
- 新增贊助者超連結；由 @ansxuman 於[#3958](https://github.com/wailsapp/wails/pull/3958)中實作
- 支援 Linux 的 deb、rpm 與 Arch Linux 套件建置；由
- 新增 Darwin 通用版本與套件的建置支援；由
- 將事件文件新增至網站；由
- 為非 SSR 開發設定的 sveltekit 與 sveltekit-ts 範本
- 使用新的`wails3 update build-assets`命令更新建置資產；由
- 用於測試 HTML 拖放 API 的範例；由
- 檔案關聯支援；由[leaanthony](https://github.com/leaanthony)於
- 新增`wails3 generate runtime`命令；由
- 新增`InitialPosition`選項，用於指定視窗是否應置中，或
- 在`application`套件中新增`Path`與`Paths`方法；由
- 新增`GeneralAutofillEnabled`與`PasswordAutosaveEnabled` Windows 選項
- 新增擷取呼叫服務方法之視窗的功能；由
- 新增 WebView2 的`EnabledFeatures`與`DisabledFeatures`選項；由
- ⊞ 新增 DIP 系統，以強化高 DPI 顯示器支援；由
- ⊞ 視窗類別名稱選項；由[windom](https://github.com/windom/)於
- 擴充服務以提供外掛程式功能。由
- 🐧 WindowDidMove / WindowDidResize 事件，於
- ⊞ WindowDidResize 事件，於
-  新增 ApplicationShouldHandleReopen 事件，以便處理 Dock
-  將 getPrimaryScreen/getScreens 新增至實作；由 @tmclane 於
-  新增在 macOS 全螢幕模式下顯示工具列的選項；由
- 🐧 新增 onKeyPress 邏輯，將 Linux 按鍵轉換為快速鍵
- 🐧 新增`run:linux`任務；由
- 匯出`SetIcon`方法；由[@almas-x](https://github.com/almas-x)於
- 改善`OnShutdown`；由[@almas-x](https://github.com/almas-x)於
- 在`Window`介面中還原`ToggleMaximise`方法；由
- 為`Environment()`新增更多資訊。由 @leaanthony 於
- 在`Window`介面上公開`WebviewWindow.IsFocused`方法；由
- 在 WML 系統中支援以空格分隔的多個觸發事件；由
- 從隨附的 JS 執行階段指令碼匯出 ESM；由
- 新增繫結產生器旗標，以使用隨附的 JS 執行階段指令碼，而不是
- 在 Linux 上實作`setIcon`；由[@abichinger](https://github.com/abichinger)實作
- 在 dev 命令中新增`-port`旗標，並支援環境變數
- 新增繫結方法呼叫的測試；由
- ⊞ 為已建立的視窗新增`SetIgnoreMouseEvents`；由
-  新增設定視窗堆疊層級（順序）的功能；由

### 修正

- 修正 Retina Mac 上`Screen.Bounds`、`WorkArea`與`Size`減半的問題：將 NSScreen 點數值轉換為`Physical*`欄位中的裝置像素，並填入頂層`Screen.X`/`Y`，使多顯示器接觸偵測與工作區域配置正確；由 @wayneforrest 於[PR](https://github.com/wailsapp/wails/pull/5168)中修正
- 修正 ScreenManager 中的資料競爭；此問題會在顯示器設定變更時（例如睡眠／喚醒期間熱插拔外接顯示器）導致 WebKit DisplayLink 死結
- 當 Assets.car 存在時，直接將 CFBundleIconName 設為 appicon，由 @symball 於[PR](https://github.com/wailsapp/wails/pull/5154)中完成
- 修正在 Fedora、openSUSE、Arch 與 NixOS 上回報錯誤 WebKitGTK 套件的`wails3 doctor`問題——由於v3在編譯時需要4.1 API，因此已移除4.0備援項目（#5071）
- 修正 openSUSE webkit2gtk doctor 套件名稱（`webkit2gtk4_1-devel` → `webkit2gtk3-devel`，後者才是正確的 openSUSE 套件名稱）（#5071）
- 修正在桌面開發模式中缺少`/wails/custom.js`時發生的`Unexpected token '<'`錯誤。為`/wails/custom.js`新增明確的404處理常式，並在`loadOptionalScript`中加入不區分大小寫的`Content-Type`驗證，以防止將 HTML SPA 備援內容注入為 JavaScript。（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- 修正 macOS 上系統匣選單的醒目提示狀態——選單開啟時，圖示現在會顯示為已選取狀態（#4910）
- 修正 macOS 上系統匣附加視窗出現在其他視窗後方的問題——現在會使用正確的快顯視窗層級（#4910）
- 修正文件中不正確的`@wailsio/runtime`匯入範例（#4989）
- 修正 darwin 上無框線視窗無法最小化的問題（#4294）
- 將`node_modules/`排除於 go-task 的最新狀態檢查之外，修正`wails3 build`和`wails3 dev`期間長達20-30分鐘的停滯。先前，`sources: "**/*"` glob 會導致 go-task 列舉`node_modules/`中的每個檔案並計算其總和檢查碼（超過50000-100000個檔案，且包含 MUI 等龐大相依套件），在 Windows/NTFS 上尤其緩慢（#4939）
- 修正 C 的`Screen` typedef 與 X11 Xlib.h 衝突所造成的 GTK4 建置失敗（#4957）
- 修正 macOS 上 Dock 徽章方法的一致性
- 修正`InvisibleTitleBarHeight`套用至所有 macOS 視窗，而非僅套用至無框線或透明標題列視窗的問題（#4960）
- 啟用`InvisibleTitleBarHeight`時，在視窗邊緣附近略過拖曳起始動作，以修正從頂端角落調整視窗大小時發生晃動／抖動的問題（#4960）
- 修正 JS/TS 繫結中以列舉鍵產生對應型別的問題（#4437），由 @fbbdev 完成
- 修正 Windows 上檔案拖放在非100%顯示縮放比例下無法運作的問題
- 修正在 Windows 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 修正 Windows 上檔案放置座標位於錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正 Linux 上使用暫留效果時，檔案拖放無法可靠運作的問題
- 修正在 Linux 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 透過使用`gtk_window_present()`，修正 Linux/GTK4 上顯示／隱藏視窗時有時會還原為最小化狀態的問題（#4957）
- 透過`XTranslateCoordinates`／`XMoveWindow`新增以 X11 為條件的支援，修正 Linux/GTK4 上取得／設定視窗位置時一律傳回0,0的問題（#4957）
- 新增以訊號為基礎的大小限制來取代已移除的`gtk_window_set_geometry_hints`，修正 Linux/GTK4 上未強制執行視窗大小上限的問題（#4957）
- 透過`gdk_monitor_get_scale`（GTK 4.14+）實作正確的 PhysicalBounds 計算與小數縮放支援，修正 Linux/GTK4 上的 DPI 縮放問題
- 修正在 Linux/GTK4 上建立新視窗時選單項目重複的問題
- 修正 JS/TS 繫結中以列舉鍵產生對應型別的問題（#4437），由 @fbbdev 完成
- 修正 Windows 上檔案拖放在非100%顯示縮放比例下無法運作的問題
- 修正在 Windows 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 修正 Windows 上檔案放置座標位於錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正 Linux 上使用暫留效果時，檔案拖放無法可靠運作的問題
- 修正在 Linux 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 透過`gdk_monitor_get_scale`（GTK 4.14+）實作正確的 PhysicalBounds 計算與小數縮放支援，修正 Linux/GTK4 上的 DPI 縮放問題
- 修正在 Linux/GTK4 上建立新視窗時選單項目重複的問題
- 修正 JS/TS 繫結中以列舉鍵產生對應型別的問題（#4437），由 @fbbdev 完成
- 修正 App.Window.Current() 未從主執行緒存取 AppKit API 而在 macOS 上造成的「幽靈視窗」問題（#4947），由 @wimaha 完成
- 透過實作 WKUIDelegate runOpenPanelWithParameters，修正 HTML `<input type="file">`在 macOS 上無法運作的問題（#4862）
- 修正在 macOS/Linux 上使用`@wailsio/runtime` npm 模組時，原生檔案拖放無法運作的問題（#4953），由 @leaanthony 完成
- 修正跨套件型別別名的繫結產生問題（#4578），由 @fbbdev 完成
- 修正 Linux 上因違反 GTK 執行緒安全性而導致 OpenFileDialog 當機的問題（#3683），由 @ddmoney420 完成
- 修正在已隱藏或已銷毀的視窗上呼叫`Focus()`時發生 SIGSEGV 當機的問題（#4890），由 @ddmoney420 完成
- 修正在 Linux 上設定空白圖示或點陣圖時可能發生 panic 的問題（#4923），由 @ddmoney420 完成
- 修正在 macOS 上從服務繫結呼叫 ErrorDialog 時發生當機的問題（#3631），由 @leaanthony 完成
- 讓選單可在 Windows 作業系統上的`v3\examples\dialogs`中顯示，由 @ndianabasi 完成
- 修正頁面重新載入期間造成 TypeError 的競爭條件（#4872），由 @ddmoney420 完成
- 移除`Collector.IsVoidAlias()`方法中的全域狀態，以修正繫結產生器測試輸出不正確的問題（#4941），由 @fbbdev 完成
- 修正`<input type="file">`檔案選擇器在 macOS 上無法運作的問題（#4862），由 @leaanthony 完成
- 修正`Position()`與`SetPosition()`在 macOS 上使用不一致的座標系統，導致儲存／還原狀態時視窗位置偏移的問題（#4816），由 @leaanthony 完成
- 修正已透過應用程式資訊清單設定 DPI 感知時，SetProcessDpiAwarenessContext 發生「存取遭拒」錯誤的問題（#4803）
- 更新鍵盤快速鍵文件頁面，並修正`KeyBinding.Add`回呼參數的型別，由 @ndianabasi 完成
- 修正產生自訂繫結的相關文件，必須使用`-d String`，而非`-o String`
- 修正在`menu.Update()`時選單未清除子項目的問題
- 修正文件中過時的 Manager API 參照（更新31檔案，以使用`app.Window.New()`、`app.Event.Emit()`等新模式），由 @leaanthony 完成
- 修正 WebKit 覆寫訊號處理常式，導致 Linux 上與 JS 繫結的 Go 方法發生 panic 時程式當機的問題（#3965），由 @leaanthony 完成
- 修正 SaveFileDialog.SetFilename() 在 Linux 上沒有作用的問題（#4841），由 @samstanier 完成
- 修正拖放範例中的放置座標顯示為 undefined 的問題
- 修正 APP_NAME 包含空格時，macOS 應用程式套件建立失敗的問題（大括號展開問題）
- 修正在 Windows 上呼叫服務方法時發生索引超出範圍 panic 的問題（還原 goccy/go-json）
- 修正 Windows 上檔案拖放在非100%顯示縮放比例下無法運作的問題
- 修正在 Windows 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 修正 Windows 上檔案放置座標位於錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正 Linux 上啟用懸停效果時，檔案拖放功能無法穩定運作的問題
- 修正 Linux 上啟用檔案放置功能時，HTML5 內部拖放功能失效的問題
- 由 @ndianabasi 更新所有作業系統之 Taskfile.yml 檔案中的全部命令，以支援 `APP_NAME` 等變數中的空格
- 由 @wux1an 修正在 Linux 上執行 'build:universal:lipo:go' 任務時的命令引數錯誤
- 由 @wux1an 修正在 Linux 上執行 'wails3 build GOOS=darwin GOARCH=arm64' 時的 Docker 錯誤「undefined symbol: **<em>ubsan</em>handle_xxxxxxx」
- 由 @leaanthony 整合自訂通訊協定文件，並新增 Universal Links 章節
- 由 @leaanthony 新增防護以避免同時呼叫 TrackPopupMenuEx，修正在 Windows 上重複點擊系統匣圖示時選單當機的問題（#4151）
- 由 @leaanthony 防止在 app.Run() 之前呼叫 systray.Run() 時應用程式當機
- 由 @leaanthony 修正啟用 ApplicationShouldTerminateAfterLastWindowClosed 時，透過 Hide()/Show() 切換視窗可見性會導致 macOS 當機的問題（#4389）
- 由 @leaanthony 修正在 macOS 與 Windows 上重複開啟內容選單時的記憶體洩漏問題（#4012）
- 由 @leaanthony 修正 macOS 上未重複使用內容選單原生資源，導致每次顯示時都重新建立選單的問題（#4012）
- 由 @leaanthony 修正在應用程式以 `Hidden: true` 啟動時，點擊 macOS Dock 圖示不會顯示隱藏視窗的問題（#4583）
- 由 @leaanthony 修正 CGO 呼叫使用錯誤的視窗指標型別，導致 macOS 列印對話方塊無法開啟的問題（#4290）
- 由 @leaanthony 修正 appmenu-gtk-module 存取尚未實現的視窗，導致 Wayland 上的視窗選單當機的問題（#4769）
- 由 @leaanthony 修正應用程式名稱含有無效字元（空格、括號等）時 GTK 應用程式當機的問題
- 由 @overlordtm 修正在 Windows 上初始化拖放功能時出現「not enough memory」錯誤的問題（#4701）
- 由 @leaanthony 修正 URI 跳脫不正確，導致 Linux 上的檔案總管開啟錯誤目錄的問題（#4397）
- 由 @leaanthony 自動偵測 `.relr.dyn` ELF 區段並停用剝除，修正現代 Linux 發行版（Arch、Fedora 39+、Ubuntu 24.04+）上的 AppImage 建置失敗問題（#4642）
- 由 @leaanthony 修正 `wails doctor` 在 Fedora／DNF 系統上誤報已安裝 webkit 套件的問題（#4457）
- 由 @mbaklor 修正預設的 `config.yml` 會搭配生產環境建置執行 `wails3 dev` 的問題
- 由 @leaanthony 修正 iOS 服務存根匯入不存在的套件而導致建置失敗的問題
- 由 @leaanthony 修正 debug/info 方法中的結構化記錄導致「no formatting directives」錯誤的問題
- 由 @leaanthony 移除從行動平台合併中意外納入的臨時偵錯列印陳述式
- 由 @leaanthony 自動停用 DMA-BUF 轉譯器，修正在使用 NVIDIA GPU 的 Wayland 上 WebKitGTK 當機的問題（錯誤 71 Protocol error）
- 解決 Linux 上 `application.WebviewWindowOptions.BackgroundColour` 中的 Alpha 值遭忽略的問題（[#4722](https://github.com/wailsapp/wails/pull/4722)，@BradHacker）
- 修正在未提供自訂圖示時，Windows 系統匣圖示不會預設為應用程式圖示的問題（#4704）
- 追蹤 `HICON` 的擁有權，確保只銷毀使用者建立的控制代碼，以防止 Explorer 回收時發生當機（#4653）。
- 在銷毀期間釋放 Windows 系統佈景主題監聽器及保留的系統匣圖示，以免 goroutine 與裝置內容持續洩漏（#4653）。
- 將系統匣工具提示截斷至 127 個 UTF-16 單位，以免損毀代理配對與多位元組字符（#4653）。
- 修正 Windows 套件任務失敗的問題（#4667）
- 修正 Linux taskfile 中的 Linux AppImage appicon 變數 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 修正 go-webview2 v1.0.22 簽章變更所導致的 Windows 建置錯誤（#4513、#4645）
- 修正 Linux taskfile 中的 Linux AppImage appicon 變數 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 由 @Tolfx 在 #4510 中移除 `<.Info.Protocol>` 至 `<.Protocol>`，修正 Linux desktop.tmpl 的通訊協定範圍
- 由 @Etesam913 修正 Liquid Glass 示範中的重複定義錯誤，詳見 [#4542](https://github.com/wailsapp/wails/pull/4542)
- 由 [@JackDoan](https://github.com/JackDoan) 修正 Linux 上的系統匣選單更新，詳見 [#4604](https://github.com/wailsapp/wails/issues/4604)
- 由 @leaanthony 在 [#4612](https://github.com/wailsapp/wails/pull/4612) 中修正在 Windows 上建立隱藏視窗時出現白色視窗的問題
- 由 @rxliuli 在 [#4617](https://github.com/wailsapp/wails/pull/4617) 中修正文件內 notifications 套件的匯入路徑
- 由 @leaanthony 在 #4616 中修正使用 npm 套件 @wailsio/runtime 時拖放功能無法運作的問題（#4489）
- Windows：由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/4600) 中修正啟動時視窗閃爍，以及隱藏視窗遭錯誤顯示的問題。
- 由 [@samstanier](https://github.com/samstanier) 修正 Wayland 視窗最大化時的尺寸問題（https://github.com/wailsapp/wails/issues/4429）
- 由 [@samstanier](https://github.com/samstanier) 修正 Wayland 視窗最大化時的尺寸問題（https://github.com/wailsapp/wails/issues/4429）
- 由 @Etesam913 修正 Liquid Glass 示範中的重複定義錯誤，詳見 [#4542](https://github.com/wailsapp/wails/pull/4542)
- 由 @jghiloni 修正 AssetServer 在 macOS 上可能當機的問題，詳見 [#4576](https://github.com/wailsapp/wails/pull/4576)
- 修正使用 NextJs 建置時的編譯問題。由 @rev42 在 [#4585](https://github.com/wailsapp/wails/pull/4585) 中修正
- 由 @riadafridishibly 修正每夜版本發行管線，詳見 [#4597](https://github.com/wailsapp/wails/pull/4597)
- 由 @Etesam913 修正 Liquid Glass 示範中的重複定義錯誤，詳見 [#4542](https://github.com/wailsapp/wails/pull/4542)
- 由 @jghiloni 修正 AssetServer 在 macOS 上可能當機的問題，詳見 [#4576](https://github.com/wailsapp/wails/pull/4576)
- 修正使用 NextJs 建置時的編譯問題。由 @rev42 在 [#4585](https://github.com/wailsapp/wails/pull/4585) 中修正
- 由 @riadafridishibly 修正每夜版本發行管線，詳見 [#4597](https://github.com/wailsapp/wails/pull/4597)
- 由 @Etesam913 修正 Liquid Glass 示範中的重複定義錯誤，詳見 [#4542](https://github.com/wailsapp/wails/pull/4542)
- 由 @PPTGamer 在 [PR](https://github.com/wailsapp/wails/pull/4492) 中修正 Windows 上的 SetBackgroundColour
- 由 @yulesxoxo 在 [PR #4476](https://github.com/wailsapp/wails/pull/4476) 中更新文件，以反映 Manager API 重構所帶來的變更
- 修正 Linux taskfile 中 Linux .desktop 檔案的 appicon 變數 [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- 由 @yulesxoxo 在 [PR #4476](https://github.com/wailsapp/wails/pull/4476) 中更新文件，以反映 Manager API 重構所帶來的變更
- 由 @leaanthony 在 [#4460](https://github.com/wailsapp/wails/pull/4460) 中修正於 [#4456](https://github.com/wailsapp/wails/issues/4456) 回報的 Windows nil 指標解除參照錯誤
- 在 macOS WKWebView 中新增對 `allowsBackForwardNavigationGestures` 的支援，以啟用雙指滑動導覽手勢（#1857）
- 修正由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中處理的問題：起初設為停用的選單項目無法使用 onClick。感謝 @IanVS 進行初步調查。
- 修正建置失敗時未清理 Vite 伺服器的問題（#4403）
- 修正在 Windows 上關閉或取消`SaveFileDialog`時發生 panic 的問題。由 @hkhere 在[PR](https://github.com/wailsapp/wails/pull/4284)中修正
- 由[@mbaklor](https://github.com/mbaklor)在[#4259](https://github.com/wailsapp/wails/pull/4259)中修正 Windows 上 HTML 層級的拖放功能
- 新增 macOS WKWebView 對`allowsBackForwardNavigationGestures`的支援，以啟用雙指滑動導覽手勢（#1857）
- 修正由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中處理的問題：起初設為停用的選單項目無法使用 onClick。感謝 @IanVS 進行初步調查。
- 修正建置失敗時未清理 Vite 伺服器的問題（#4403）
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4450)中修正 Windows 上的通知剖析
- 由[@kodumulo](https://github.com/kodumulo)在[#4390](https://github.com/wailsapp/wails/issues/4390)中修正 doctor 命令，使其檢查 Windows SDK 相依套件
- 由[@etesam913](https://github.com/etesam913)在[#4366](https://github.com/wailsapp/wails/pull/4366)中修正 Mac 上 processURLRequest 的 nil 指標解參照問題
- 由[@bh90210](https://github.com/bh90210)在[#4287](https://github.com/wailsapp/wails/pull/4287)中修正 Linux 上導致已篩選對話方塊無法運作的錯誤
- 由[@leaanthony](https://github.com/leaanthony)在[#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)中修正 Windows 與 Linux 上的「編輯」選單問題
- 由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中將 macOS .plist 檔案的最低系統版本從10.13.0更新為10.15.0
- 由[@leaanthony](https://github.com/leaanthony)修正視窗 ID 跳號問題
- 由[@leaanthony](https://github.com/leaanthony)修正呼叫 RegisterContextMenu 時選單為 nil 的問題
- 由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中修正繫結產生器輸出中的相依性循環
- 由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中修正繫結產生器輸出中的先使用後定義錯誤
- 由[@fbbdev](https://github.com/fbbdev)在[#4023](https://github.com/wailsapp/wails/pull/4023)中將建置旗標傳遞給繫結產生器
- 由[@leaanthony](https://github.com/leaanthony)將 Windows Taskfile 中的路徑改用正斜線，以確保其可在非 Windows 平台上運作
- 由[@leaanthony](https://github.com/leaanthony)修正 Mac + Mac JS 事件
- 由[@leaanthony](https://github.com/leaanthony)修正 macOS 上的事件死結
- 由[@leaanthony](https://github.com/leaanthony)修正在 Windows 上初始化 Window 時，已提供 HTML 但未提供 JS 所導致的`Parameter incorrect`錯誤
- 由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中修正資產伺服器用於偵測內容類型的回應前綴大小
- 由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中修正資產伺服器在根索引路徑上對非404回應的處理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器測試泛型型別屬性時的未定義行為
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正當基礎型別與具名包裝型別的屬性不相同時，繫結產生器對模型的輸出
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器對 map 鍵型別的輸出與預處理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器對實作 marshaler 介面的 struct 所產生的輸出
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器對涉及泛型型別之型別循環的偵測
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器輸出中對未匯出模型的無效參照
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中將注入的程式碼移至服務檔案末尾
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器對檔案關閉操作所傳回錯誤的處理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中隱藏針對只定義生命週期或 http 方法、但未定義其他繫結方法之服務的警告
- 由[@marcus-crane](https://github.com/marcus-crane)在[#4056](https://github.com/wailsapp/wails/pull/4056)中修正非 React 範本在使用淺色系統配色時無法顯示 Hello World 頁尾的問題
- 由[@leaanthony](https://github.com/leaanthony)修正 macOS 上隱藏的選單項目
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中修正訊息處理器對錯誤的處理與格式化
-  由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中修正結束應用程式時略過服務關閉程序的問題
-  由[@leaanthony](https://github.com/leaanthony)確保選單更新在主執行緒上進行
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中強化拖曳與調整大小機制，使其更貼近預期的平台行為
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中修正[#4097](https://github.com/wailsapp/wails/issues/4097)：Webpack/angular 會捨棄執行階段初始化程式碼的問題
- 由[@IanVS](https://github.com/IanVS)在[#4116](https://github.com/wailsapp/wails/pull/4116)中修正起初隱藏的選單項目
- 修正當無副檔名的請求中`[request]`不存在但`[request].html`存在時，assetFileServer 不會提供`.html`檔案的問題
- 由[@robin-samuel](https://github.com/robin-samuel)在[#4125](https://github.com/wailsapp/wails/pull/4125)中修正圖示產生路徑
- 由[@oSethoum](https://github.com/osethoum)在[#4130](https://github.com/wailsapp/wails/pull/4130)中修正未發出`fullscreen`、`unfullscreen`、`unminimise`和`unmaximise`事件的問題
- 由[@robin-samuel](https://github.com/robin-samuel)在[#4126](https://github.com/wailsapp/wails/pull/4126)中修正設定內預設版本的前綴不正確所造成的 NSIS 錯誤
- 由[TheGB0077](https://github.com/TheGB0077)在[#4188](https://github.com/wailsapp/wails/pull/4188)中修正 Dialogs 執行階段函式在 Windows 上傳回逸出路徑的問題
- 由[@leaanthony](https://github.com/leaanthony)修正 HKCU 中的 Webview2 偵測路徑。
- 由[@leaanthony](https://github.com/leaanthony)修正 macOS 上的輸入問題。
- 由 [@yulesxoxo](https://github.com/yulesxoxo) 在 [#4219](https://github.com/wailsapp/wails/pull/4219) 中修正 Windows 圖示產生任務的檔案名稱。
- 由 [@leaanthony](https://github.com/leaanthony) 根據 @kron 的工作，修正無框線視窗的透明度問題。
- 由 [@leaanthony](https://github.com/leaanthony) 根據 @kron 的工作，修正視窗停用或最小化時的焦點呼叫。
- 由 [@leaanthony](https://github.com/leaanthony) 根據 @kron 的工作，修正工作列重新啟動後系統匣未顯示的問題。
- 在 [#4245](https://github.com/wailsapp/wails/pull/4245) 中修正 fallbackResponseWriter 未實作 Flush() 的問題
- 由 [@superDingda] 在 [#4236](https://github.com/wailsapp/wails/issues/4236) 中修正 fallbackResponseWriter 未實作 Flush() 的問題
- 由 [@joshhardy](https://github.com/joshhardy) 在 [#4354](https://github.com/wailsapp/wails/pull/4354) 中修正 macOS 視窗關閉時仍有待處理的非同步 Go 綁定函式呼叫所造成的當機
- 由 [@leaanthony](https://github.com/leaanthony) 修正 Windows 效率模式啟動時的競爭條件
- 由 [@leaanthony](https://github.com/leaanthony) 修正 Windows 圖示控制代碼的清理。
- 由 [@PPTGamer](https://github.com/PPTGamer) 在 [#4375](https://github.com/wailsapp/wails/pull/4375) 中修正 Windows 上的 `OpenFileManager`。
- 由 @atterpac 在 [#3979](https://github.com/wailsapp/wails/pull/3979) 中修正 Linux 的最小／最大寬度選項
- 由 @atterpac 在 [#3966](https://github.com/wailsapp/wails/pull/3966) 中透過提升 npm 版本更新 TypeScript 範本的型別定義
- 由 @atterpac 在 [#3945](https://github.com/wailsapp/wails/pull/3945) 中修正 SvelteKit 範本的 CSS 參照
- 由 [@leaanthony](https://github.com/leaanthony) 確保 window run() 中的關鍵回呼會在主執行緒上呼叫
- 由 [@leaanthony](https://github.com/leaanthony) 修正對話方塊目錄選擇器範例
- 由 [@leaanthony](https://github.com/leaanthony) 建立 index.html 遺失時顯示的新中文錯誤頁面
-  由 [@leaanthony](https://github.com/leaanthony) 確保 `windowDidBecomeKey` 回呼在主執行緒上執行
-  由 [@leaanthony](https://github.com/leaanthony) 為無框線視窗新增全螢幕支援
-  由 [@leaanthony](https://github.com/leaanthony) 改善視窗銷毀邏輯
-  由 [@leaanthony](https://github.com/leaanthony) 修正視窗附加至系統匣時的位置邏輯
-  由 [@leaanthony](https://github.com/leaanthony) 為無框線視窗新增全螢幕支援
- 由 [@leaanthony](https://github.com/leaanthony) 修正事件處理
- 由 [@leaanthony](https://github.com/leaanthony) 修正視窗關閉邏輯
- 由 [@leaanthony](https://github.com/leaanthony) 將通用 Taskfile 改為預設替 TypeScript 範本產生 TypeScript 綁定
- 由 [@mmalcek](https://github.com/mmalcek) 在 [#3990](https://github.com/wailsapp/wails/pull/3990) 中修正沒有開啟任何視窗／僅使用系統匣時，收到 WM_CLOSE 訊息便關閉應用程式的問題
- 由 @5aaee9 在 [#3192](https://github.com/wailsapp/wails/pull/3192) 中修正 garble 建置
- 由 [@leaanthony](https://github.com/leaanthony) 修正 Windows NSIS 建置
- 修正 Linux 多選對話方塊中因未關閉而造成的死結
- 修正 Windows 建置期間對 .syso 檔案的跨平台清理，由
- 修正 amd64 AppImage 編譯，由 @atterpac 在
- 修正建置資產更新，由 @ansxuman 在
- 由 @atterpac 修正 Linux 系統匣的 `OnClick` 與 `OnRightClick` 實作
- 修正 `AlwaysOnTop` 無法在 Mac 上運作的問題，由
-  修正 `application.NewEditMenu` 包含重複項目的問題
- 🐧 修正 aarch64 編譯
- ⊞ 修正選項按鈕群組選單項目，由
- 修正 MacOS 上當 'name' 與 'outputfilename' 時，建置可執行 .app 所發生的錯誤
- 修正拖放範例中使用 customEventProcessor 的錯誤，由
- 🐧 修正新增 IgnoreMouseEvents 所引發的 Linux 編譯錯誤，由
- ⊞ 修正 syso 圖示檔案產生錯誤，由
- 🐧 納入可在 Wayland 上原生執行的修正，來源為
- 不要在以下項目中綁定內部服務方法
- ⊞ 修正系統匣啟動時的 panic，在
- 不要在以下項目中綁定內部服務方法
- ⊞ 修正系統匣啟動時的 panic，在
- 大幅重構選單項目與事件處理。目前主要改善 macOS。由
- 修正外掛程式與事件重構後的測試，在
- ⊞ 修正 `Failed to unregister class Chrome_WidgetWin_0` 警告。由
- 模組問題
- 由 [atterpac](https://github.com/atterpac) 修正調整大小事件的訊息傳遞，在
- 🐧修正 NixOS 上的佈景主題處理錯誤，由
- 修正 Windows 上跨磁碟區安裝專案的問題，由
- 修正 React 範本 CSS，使頁尾得以顯示，由
- 透過更新至最新的 refresh，修正在開發模式下工作時產生的殭屍處理程序
- 由 [Atterpac](https://github.com/atterpac) 修正 AppImage 的 WebKit 檔案來源設定
- 由 [Atterpac](https://github.com/Atterpac) 修正 Doctor 的 apt 套件驗證，在
- 修正結束時應用程式凍結的問題（Darwin），由 @5aaee9 在
- 修正 Windows 上範例的背景色彩，由
- 由 [mmghv](https://github.com/mmghv) 修正預設內容功能表，在
- 修正 Darwin 上方向鍵的十六進位值，由
- 讓 Windows 上的拖放功能得以運作。新增者為
- 修正使用者沒有適當驅動程式時，Linux 上 Doctor 的錯誤
- 修正啟動時的 DPI 縮放（Windows）。由 [@almas-x](https://github.com/almas-x) 變更，在
- 修正 `go.mod` 中的取代行，使其使用相對路徑。修正包含以下內容的 Windows 路徑
- 修正 MacOS 系統匣未附加視窗時的點擊處理，由
- 修正未知選項導致 Windows 建置失敗的問題，由
- 修正 Windows 上沒有附加視窗時，以滑鼠左鍵按一下系統匣圖示所造成的當機
- 修正開啟視窗兩次時錯誤的 baseURL，由 @5aaee9 在 PR
- 修正 `WebviewWindow.Restore` 方法中 if 分支的順序，由
- 在以下情況下，於多次 `GetStartURL` 呼叫間正確計算 `startURL`
- 修正 `Screen` 結構的 JS 型別，使其與對應的 Go 型別一致，由
- 修正 `WML.Reload` 方法，以確保正確清理已註冊的事件
- 修正 Linux 上自訂內容功能表立即關閉的問題，由
- 修正繫結產生的模型檔案輸出路徑與副檔名
- 修正繫結產生的 JS 程式碼中模型檔案的匯入路徑
- 修正某些 Linux 發行版上的拖放功能，作者為
- 修正使用`wails3 task dev`時缺少 macOS 工作的問題，作者為
- 修正註冊事件時造成對 nil map 指派的問題，作者為
- 修正繫結方法參數的反序列化問題，作者為
- 修正處理繫結方法多個回傳值的問題，作者為
- 修正 doctor 無法偵測並非透過系統套件管理員安裝之 npm 的問題
- 修正缺少 MicrosoftEdgeWebview2Setup.exe 的問題。感謝
- 修正 Linux 上因視窗 ID 處理而隨機當機的問題，作者為 @leaanthony。依據
- 修正 systemTray.setIcon 在 Linux 上造成當機的問題，作者為
- 修正`setFrameless`函式首次呼叫時未套用視窗框架的問題，平台為

### 變更

- **重大變更**：產生的 JS/TS 繫結現在會將 map 鍵標示為選用，以準確反映 Go map 的語意。TypeScript 中存取 map 值現在會回傳`T | undefined`而非`T`，因此必須進行 null 檢查或使用斷言（#4943），作者為`@fbbdev`
- 因應`@wailsio/runtime`的變更，將`Event`改為`Events`，並在`Features/Events/Event System`的文件中使用適當的函式呼叫，作者為 @AbdelhadiSeddar
- 將`EnabledFeatures`、`DisabledFeatures`和`AdditionalBrowserArgs`從各視窗選項移至應用程式層級的`Options.Windows`（#4559），作者為 @leaanthony
- 更新`Drag N Drop`範例的 README，並特別說明該範例展示了`Internal Drag and Drop`，作者為 @ndianabasi
- 將多項偵錯日誌的層級從 Info 改為 Debug（作者為 @mbaklor）
- <strong>重大變更：</strong>將視窗選項中的`EnableDragAndDrop`重新命名為`EnableFileDrop`
- <strong>重大變更：</strong>將事件內容中的`DropZoneDetails`重新命名為`DropTargetDetails`
- <strong>重大變更：</strong>將`WindowEventContext`上的`DropZoneDetails()`方法重新命名為`DropTargetDetails()`
- <strong>重大變更：</strong>移除`WindowDropZoneFilesDropped`事件，請改用`WindowFilesDropped`
- <strong>重大變更：</strong>將 HTML 屬性從`data-wails-dropzone`改為`data-file-drop-target`
- <strong>重大變更：</strong>將 CSS hover 類別從`wails-dropzone-hover`改為`file-drop-target-active`
- <strong>重大變更：</strong>從 Windows 移除`DragEffect`、`OnEnterEffect`及`OnOverEffect`選項（這些選項原屬於已移除的 IDropTarget）
- 所有執行階段 JSON 處理（方法繫結、事件、webview 請求、通知、kvstore）均改用 goccy/go-json，使效能提升21-63%，記憶體配置減少40-60%
- 最佳化 BoundMethod 結構配置並快取 isVariadic 旗標，以降低每次呼叫的額外負擔
- 對具有`<=8`個引數的方法使用堆疊配置的引數緩衝區，以避免堆積配置
- 最佳化方法呼叫的結果收集方式，避免為單一回傳值配置切片
- MIME 類型快取改用 sync.Map，以提升並行效能
- 讀取 HTTP 傳輸請求本文時使用緩衝區集區
- 在內容類型嗅探器中延遲配置 CloseNotify 通道，以減少每個請求的配置
- 移除資產伺服器的 CSS 偵錯日誌
- 擴充 MIME 類型副檔名對應表，涵蓋50種以上的常見 Web 格式（字型、音訊、視訊等）
- 更新視窗`X/Y`選項的文件，作者為 @ruhuang2001
- 更新`Frontend Runtime`文件，加入更多產生前端繫結的選項，作者為 @ndianabasi
- 更新 Wails v3資產伺服器的文件頁面，作者為 @ndianabasi
- **重大變更**：移除套件層級的對話方塊函式（`application.InfoDialog()`、`application.QuestionDialog()`等）。請改用`app.Dialog`管理器：`app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()`
- 更新對話方塊文件以符合實際 API：使用帶有回呼的`app.Dialog.*`、`AddButton()`（而非`SetButtons()`）、`SetDefaultButton(*Button)`（而非字串）、`AddFilter()`（而非`SetFilters()`）、`SetFilename()`（而非`SetDefaultFilename()`），並使用`app.Dialog.OpenFile().CanChooseDirectories(true)`選取資料夾
- **重大變更**：現在預設建立生產版本。若要建立開發版本，請在 Taskfile 中設定`DEV=true`。產生新專案以供範例使用，作者為 @leaanthony
- 發出包含零個或一個資料引數的自訂事件時，資料值將直接指派給 Data 欄位，不再包裝於切片中，作者為[@fbbdev](https://github.com/fbbdev)，見[#4633](https://github.com/wailsapp/wails/pull/4633)
- Windows 系統匣現在會透過切換`NIS_HIDDEN`來遵循`SystemTray.Show()`/`Hide()`，讓應用程式可以真正消失並再次顯示（#4653）。
- 系統匣註冊現在會重複使用已解析的圖示、僅設定一次`NOTIFYICON_VERSION_4`，並啟用`NIF_SHOWTIP`，讓工具提示能在 Explorer 重新啟動後恢復（#4653）。
- macOS：置中視窗時改用`visibleFrame`而非`frame`，以排除選單列和 Dock 區域
- macOS：置中視窗時改用`visibleFrame`而非`frame`，以排除選單列和 Dock 區域
- 使用`-config`參數執行`wails3 update build-assets`時，透過`-product*`參數設定的值會
- `window.NativeWindowHandle()` -> `window.NativeWindow()`，作者為 @leaanthony，見[#4471](https://github.com/wailsapp/wails/pull/4471)
- 重構內部視窗處理，作者為 @leaanthony，見[#4471](https://github.com/wailsapp/wails/pull/4471)
- 移除`application.WindowIDKey`和`application.WindowNameKey`（由`application.WindowKey`取代），作者為[@leaanthony](https://github.com/leaanthony)
- ContextMenuData 現在回傳 string 而非 any，作者為[@leaanthony](https://github.com/leaanthony)
- 在 JS/TS 繫結中，固定長度陣列類型的類別欄位現在會以預期長度初始化，而不再是空陣列，作者為[@fbbdev](https://github.com/fbbdev)，見[#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData 現在回傳 string 而非 any，作者為[@leaanthony](https://github.com/leaanthony)
- `application.NewService`不再接受 options 作為選用參數（請改用`application.NewServiceWithOptions`），作者為[@leaanthony](https://github.com/leaanthony)，見[#4024](https://github.com/wailsapp/wails/pull/4024)
- 移除`nanoid`相依性，作者為[@leaanthony](https://github.com/leaanthony)
- 更新視窗範例，示範 mica/acrylic/tabbed 視窗樣式，作者為[@leaanthony](https://github.com/leaanthony)
- 在 JS/TS 繫結中，已移除`internal.js/ts`模型檔案；現在所有模型皆可在`models.js/ts`中找到，作者為[@fbbdev](https://github.com/fbbdev)，見[#4045](https://github.com/wailsapp/wails/pull/4045)
- 在 JS/TS 繫結中，具名類型絕不再呈現為其他具名類型的別名；舊行為現在僅適用於別名，作者為[@fbbdev](https://github.com/fbbdev)，見[#4045](https://github.com/wailsapp/wails/pull/4045)
- 在 JS/TS 繫結的類別模式中，型別為型別參數的結構欄位會標示為選用，且絕不會自動初始化，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 貢獻
- 從範本中移除 ESLint，由 [@IanVS](https://github.com/IanVS) 於 [#4059](https://github.com/wailsapp/wails/pull/4059) 貢獻
- 將著作權日期更新為 2025，由 [@IanVS](https://github.com/IanVS) 於 [#4037](https://github.com/wailsapp/wails/pull/4037) 貢獻
- 新增 event.Sender 的文件，由 [@IanVS](https://github.com/IanVS) 於 [#4075](https://github.com/wailsapp/wails/pull/4075) 貢獻
- 支援 Go 1.24，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- `ServiceStartup`掛鉤現在會在呼叫`App.Run`時執行，而非在`application.New`中執行，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 貢獻
- `ServiceStartup`錯誤現在會從`App.Run`傳回，而不再終止處理程序，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 貢獻
- 從 JS 發出的繫結與對話方塊呼叫現在會以錯誤物件而非字串拒絕，由 [@fbbdev](https://github.com/fbbdev) 於 [#4066](https://github.com/wailsapp/wails/pull/4066) 貢獻
- 改善 Windows 上系統匣選單的定位，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- JS 執行階段已移植至 TypeScript，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 貢獻
- 執行階段現在一經匯入便會初始化，無須等待視窗載入，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 貢獻
- 執行階段不再匯出 init 方法；可使用具副作用的匯入來初始化，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 貢獻
- 繫結方法現在會傳回`CancellablePromise`；若取消，該物件會以`CancelError`拒絕。呼叫的實際結果會被捨棄，由 [@fbbdev](https://github.com/fbbdev) 於 [#4100](https://github.com/wailsapp/wails/pull/4100) 貢獻
- 內建服務型別現在統一命名為`Service`，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 貢獻
- 含選項的內建服務建立函式現在統一命名為`NewWithConfig`，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 貢獻
- 為與 Go API 保持一致，`sqlite`服務的`Select`方法現在命名為`Query`，由 [@fbbdev](https://github.com/fbbdev) 於 [#4067](https://github.com/wailsapp/wails/pull/4067) 貢獻
- 範本：將執行階段移至「dependencies」，並整理 package.json 檔案，由 [@IanVS](https://github.com/IanVS) 於 [#4133](https://github.com/wailsapp/wails/pull/4133) 貢獻
- 在開發模式中建立並臨時簽署應用程式套件，以啟用特定 macOS API，由 [@popaprozac](https://github.com/popaprozac) 於 [#4171](https://github.com/wailsapp/wails/pull/4171) 貢獻
- 將建置資產移至平台專用目錄，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 將 Taskfile 移至平台專用目錄並重新命名，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 大幅改善缺少`index.html`時的使用體驗，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- [Windows] 改善最小化與還原的效能，由 [@leaanthony](https://github.com/leaanthony) 貢獻。以 [562589540](https://github.com/562589540) 原先提出的[PR](https://github.com/wailsapp/wails/pull/3955)為基礎
- 移除`ShouldClose`選項（請改為註冊 events.Common.WindowClosing 的掛鉤），由 [@leaanthony](https://github.com/leaanthony) 貢獻
- [Windows] 減少開啟視窗時的閃爍，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 移除`Window.Destroy`，因為它原本是內部函式，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 將`WindowClose`事件重新命名為`WindowClosing`，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 前端建置現在會依建置類型使用 vite 環境「development」或「production」，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 更新至 go-webview2 v1.19，由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 確保使用 Taskfile 的分支版本，由 @leaanthony 貢獻
- 更新 Taskfile 的分支版本，以修正使用以下方式安裝時的版本問題
- 使用 Taskfile 的分支版本，以修正使用以下方式安裝時的版本問題
- `service.OnStartup`現在會在發生錯誤時關閉應用程式，並執行
- 重構系統匣點擊訊息傳遞，使其更符合使用者互動，由
- 資產嵌入包含`all:frontend/dist`，以支援會產生以下內容的框架
- Taskfile 重構，由 [leaanthony](https://github.com/leaanthony) 於
- 升級至`go-webview2` v1.0.16，由
- 修正`Screen`型別，使其包含`ID`而非`Id`，由
- 更新`go.mod.tmpl`的 Wails 版本以支援`application.ServiceOptions`，由
- 修正服務名稱的判定方式，由 [windom](https://github.com/windom/) 於
- mkdocs serve 現在使用 Docker，由 [leaanthony](https://github.com/leaanthony) 貢獻
- 將開發設定整合至`config.yml`，由
- 系統匣對話方塊現在預設使用應用程式圖示（若有；Windows），由
- 改善 macOS 的 GPU 與記憶體資訊回報，由
- 移除`WebviewGpuIsDisabled`和`EnableFraudulentWebsiteWarnings`
- 事件 API 變更：`On`/`Emit` -> 使用者事件，`OnApplicationEvent` ->
- 修正 Linux 上的事件 API，由 [TheGB0077](https://github.com/TheGB0077) 於
- [CI] 改善 actions，並使 actions 也能在分支版本及以下環境中執行
- 將`AbsolutePosition()`重新命名為`Position()`，由
- 將 Linux WebKit 相依套件從 webkitgtk2-4.0 更新為 webkit2gtk-4.1，以
- 隨附的 JS 執行階段指令碼現在是 ESM 模組：匯入它的 script 標籤
- `@wailsio/runtime`套件不會在`window.wails`上公開其 API
- 視窗 API 模組`@wailsio/runtime/src/window`現在會公開其所包含的
- JS 視窗 API 已更新，以符合目前的 Go `WebviewWindow`
- 繫結產生器現在預設使用依 ID 呼叫。`-id` CLI 選項
- 新的繫結程式碼配置：輸出檔案先前依資料夾組織
- 結構欄位`application.Options.Bind`已重新命名為
- 繫結服務的新語法：服務執行個體現在必須包裝在
- 在非終端機或 CI 環境中停用旋轉指示器，由

### 已移除

- **重大變更**：從各視窗的`WindowsWindow`選項中移除`EnabledFeatures`、`DisabledFeatures`和`AdditionalLaunchArgs`。請改用應用程式層級的`Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures`和`Options.Windows.AdditionalBrowserArgs`。這些旗標會全域套用至共用的 WebView2 環境（#4559），由 @leaanthony 完成
- 移除 Windows 上原生的`IDropTarget`實作，改用以 JavaScript 為基礎的方法（與v2的行為一致）
- 移除 github.com/wailsapp/mimetype 相依套件，改用擴充的副檔名對照表與標準函式庫的 http.DetectContentType，使二進位檔案大小減少約1.2MB
- 為 Linux 檔案總管實作精簡的 .desktop 檔案剖析器，以移除 gopkg.in/ini.v1相依套件，節省約45KB
- 改用 Go 1.21+ 標準函式庫的 slices 套件與精簡的內部輔助函式，從執行階段程式碼中移除 samber/lo，節省約310KB
- 從 Darwin URL 配置處理常式中移除偵錯 printf 陳述式（#4834）
- **重大變更**：移除`linux:WindowLoadChanged`事件；若要偵測 WebView 何時完成載入，請改用`linux:WindowLoadFinished`（#3896），由 @leaanthony 完成

### 重大變更

- **Manager API 重構**：將應用程式 API 從扁平結構重組為各個管理器，以改善程式碼組織並提高可探索性；由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中完成
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
- 重新命名 Service 方法：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`，由[@leaanthony](https://github.com/leaanthony)完成
- 由[@leaanthony](https://github.com/leaanthony)將`Path`和`Paths`方法移至`application`套件
- 應用程式選單現在僅適用於 macOS，由[@leaanthony](https://github.com/leaanthony)完成

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

## 已修正

## v3.0.0-alpha.56 - 2026-01-04

## 新增

## 已變更

## 已修正

## 已移除

## v3.0.0-alpha.55 - 2026-01-02

## 已變更

## 已修正

## 已移除

## v3.0.0-alpha.54 - 2025-12-29

## 新增

## 已修正

## 已移除

## v3.0.0-alpha.53 - 2025-12-27

## 新增

## 已修正

## v3.0.0-alpha.52 - 2025-12-26

## 已修正

## v3.0.0-alpha.51 - 2025-12-23

## 已修正

## v3.0.0-alpha.50 - 2025-12-21

## 已變更

## v3.0.0-alpha.49 - 2025-12-18

## 已變更

## v3.0.0-alpha.48 - 2025-12-16

## 新增

## 已變更

## 已修正

## v3.0.0-alpha.47 - 2025-12-15

## 新增

## 已修正

## v3.0.0-alpha.46 - 2025-12-14

## 新增

## 已移除

## v3.0.0-alpha.45 - 2025-12-13

## 新增

## 已修正

## v3.0.0-alpha.44 - 2025-12-12

## 新增

## 已變更

## 已修正

## v3.0.0-alpha.43 - 2025-12-11

## 新增

## v3.0.0-alpha.42 - 2025-12-10

## 新增

## v3.0.0-alpha.41 - 2025-11-23

## 已修正

## v3.0.0-alpha.40 - 2025-11-13

## 已修正

## v3.0.0-alpha.39 - 2025-11-12

## 新增

## 已變更

## v3.0.0-alpha.38 - 2025-11-04

## 新增

## 已變更

## 已修正

## v3.0.0-alpha.37 - 2025-11-02

## 已修正

## v3.0.0-alpha.36 - 2025-10-15

## 已修正

## v3.0.0-alpha.35 - 2025-10-14

## 已修正

## v3.0.0-alpha.34 - 2025-10-06

## 新增

## 已修正

## v3.0.0-alpha.33 - 2025-10-04

## 已修正

## v3.0.0-alpha.32 - 2025-10-02

## 已修正

## v3.0.0-alpha.31 - 2025-09-27

## 已修正

## v3.0.0-alpha.30 - 2025-09-26

## 已修正

## v3.0.0-alpha.29 - 2025-09-25

## 新增

## 已變更

## 已修正

## v3.0.0-alpha.29 - 2025-09-25

## 新增

## 已變更

## 已修正

## v3.0.0-alpha.27 - 2025-09-07

## 已修正

## v3.0.0-alpha.26 - 2025-08-24

## 新增

## v3.0.0-alpha.25 - 2025-08-16

## 已變更

不再被忽略，並會覆寫設定值。

## v3.0.0-alpha.24 - 2025-08-13

## 新增

## v3.0.0-alpha.23 - 2025-08-11

## 已修正

## v3.0.0-alpha.22 - 2025-08-10

## 新增

## 已變更

+ 修正範圍過廣的 Linux 套件相依性，並修正過時的 RPM 相依性。

## v3.0.0-alpha.21 - 2025-08-07

## 已修正

## v3.0.0-alpha.20 - 2025-08-06

## 已修正

## v3.0.0-alpha.19 - 2025-08-05

## 新增

## 已修正

## v3.0.0-alpha.18 - 2025-08-03

## 新增

## 已修正

## v3.0.0-alpha.17 - 2025-07-31

## 已修正

## v3.0.0-alpha.16 - 2025-07-25

## 新增

## v3.0.0-alpha.15 - 2025-07-25

## 新增

## v3.0.0-alpha.14 - 2025-07-25

## 新增

## v3.0.0-alpha.12 - 2025-07-15

### 新增

### 已修正

## v3.0.0-alpha.11 - 2025-07-12

## 新增

## v3.0.0-alpha.10 - 2025-07-06

### 破壞性變更

### 新增

### 修正

### 變更

## v3.0.0-alpha.9 - 2025-01-13

### 新增

### 修正

### 變更

## v3.0.0-alpha.8.3 - 2024-12-07

### 變更

## v3.0.0-alpha.8.2 - 2024-12-07

### 變更

`go install`，由 @leaanthony 提供

## v3.0.0-alpha.8.1 - 2024-12-07

### 變更

`go install`，由 @leaanthony 提供

## v3.0.0-alpha.8 - 2024-12-06

### 新增

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### 變更

`service.OnShutdown`適用於先前已啟動的所有服務，由 @atterpac 於 [#3920](https://github.com/wailsapp/wails/pull/3920) 提供；@atterpac 於 [#3907](https://github.com/wailsapp/wails/pull/3907) 提供；子資料夾，由 @atterpac 於 [#3887](https://github.com/wailsapp/wails/pull/3887) 提供；[#3748](https://github.com/wailsapp/wails/pull/3748)；[leaanthony](https://github.com/leaanthony)；[etesam913](https://github.com/etesam913) 於 [#3778](https://github.com/wailsapp/wails/pull/3778) 提供；[northes](https://github.com/northes) 於 [#3836](https://github.com/wailsapp/wails/pull/3836) 提供；[#3827](https://github.com/wailsapp/wails/pull/3827)；[leaanthony](https://github.com/leaanthony)；[@leaanthony](https://github.com/leaanthony)；[@leaanthony](https://github.com/leaanthony)；（已由 `EnabledFeatures` 與 `DisabledFeatures` 選項取代），由 [leaanthony](https://github.com/leaanthony) 提供

### 修正

通道變數，由 @michael-freling 於 [#3925](https://github.com/wailsapp/wails/pull/3925) 提供；[ansxuman](https://github.com/ansxuman) 於 [#3924](https://github.com/wailsapp/wails/pull/3924) 提供；[#3898](https://github.com/wailsapp/wails/pull/3898)；[#3901](https://github.com/wailsapp/wails/pull/3901)；於 [#3886](https://github.com/wailsapp/wails/pull/3886)；[leaanthony](https://github.com/leaanthony) 於 [#3841](https://github.com/wailsapp/wails/pull/3841) 提供；Darwin 編輯選單中的 `PasteAndMatchStyle` 角色，由 [johnmccabe](https://github.com/johnmccabe) 於 [#3839](https://github.com/wailsapp/wails/pull/3839) 提供；[#3840](https://github.com/wailsapp/wails/issues/3840) 於 [#3854](https://github.com/wailsapp/wails/pull/3854)，由 [kodflow](https://github.com/kodflow) 提供；[@leaanthony](https://github.com/leaanthony)；兩者不同，由 @nickisworking 於 [#3789](https://github.com/wailsapp/wails/pull/3789) 提供

## v3.0.0-alpha.7 - 2024-09-18

### 新增

[mmghv](https://github.com/mmghv) 於 [#3665](https://github.com/wailsapp/wails/pull/3665) 提供；[#3682](https://github.com/wailsapp/wails/pull/3682)；[atterpac](https://github.com/atterpac) 與 [leaanthony](https://github.com/leaanthony) 於 [#3570](https://github.com/wailsapp/wails/pull/3570) 提供

### 變更

應用程式事件 `OnWindowEvent` -> 視窗事件，由 [leaanthony](https://github.com/leaanthony) 提供；[#3734](https://github.com/wailsapp/wails/pull/3734)；以 `v3/` 或 `v3-` 為前綴的分支，由 [stendler](https://github.com/stendler) 於 [#3747](https://github.com/wailsapp/wails/pull/3747) 提供

### 修正

[etesam913](https://github.com/etesam913) 於 [#3742](https://github.com/wailsapp/wails/pull/3742) 提供；[atterpac](https://github.com/atterpac) 於 [#3721](https://github.com/wailsapp/wails/pull/3721) 提供；[atterpac](https://github.com/atterpac) 於 [#3675](https://github.com/wailsapp/wails/pull/3675) 提供；[#1811](https://github.com/wailsapp/wails/pull/1811) 於 [#3614](https://github.com/wailsapp/wails/pull/3614)，由 [@stendler](https://github.com/stendler) 提供；[#3720](https://github.com/wailsapp/wails/pull/3720)，由 [leaanthony](https://github.com/leaanthony) 提供；[#3693](https://github.com/wailsapp/wails/issues/3693)，由 [@DeltaLaboratory](https://github.com/DeltaLaboratory) 提供；[#3720](https://github.com/wailsapp/wails/pull/3720)，由 [leaanthony](https://github.com/leaanthony) 提供；[#3693](https://github.com/wailsapp/wails/issues/3693)，由 [@DeltaLaboratory](https://github.com/DeltaLaboratory) 提供；[leaanthony](https://github.com/leaanthony)；[#3746](https://github.com/wailsapp/wails/pull/3746)，由 [@stendler](https://github.com/stendler) 提供；[leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 修正

## v3.0.0-alpha.5 - 2024-07-30

### 新增

[#3580](https://github.com/wailsapp/wails/pull/3580)；[#3580](https://github.com/wailsapp/wails/pull/3580)；圖示點擊，由 @5aaee9 於 [#2991](https://github.com/wailsapp/wails/pull/2991) 提供；[#2618](https://github.com/wailsapp/wails/pull/2618)；[@fbbdev](https://github.com/fbbdev) 於 [#3282](https://github.com/wailsapp/wails/pull/3282) 提供；@[Atterpac](https://github.com/Atterpac) 於 [#3022](https://github.com/wailsapp/wails/pull/3022]) 提供；[@marcus-crane](https://github.com/marcus-crane) 於 [#3146](https://github.com/wailsapp/wails/pull/3146) 提供；[PR](https://github.com/wailsapp/wails/pull/3147)；[PR](https://github.com/wailsapp/wails/pull/3189)；[@fbbdev](https://github.com/fbbdev) 於 [#3281](https://github.com/wailsapp/wails/pull/3281) 提供；[aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)；基於 @Mai-Lapyst 的 [PR](https://github.com/wailsapp/wails/pull/2044)；[@fbbdev](https://github.com/fbbdev) 於 [#3295](https://github.com/wailsapp/wails/pull/3295) 提供；[@fbbdev](https://github.com/fbbdev) 於 [#3295](https://github.com/wailsapp/wails/pull/3295) 提供；[@fbbdev](https://github.com/fbbdev) 於 [#3295](https://github.com/wailsapp/wails/pull/3295) 提供；npm 套件，由 [@fbbdev](https://github.com/fbbdev) 於 [#3334](https://github.com/wailsapp/wails/pull/3334) 提供；於 [#3354](https://github.com/wailsapp/wails/pull/3354)；`WAILS_VITE_PORT`，由 [@abichinger](https://github.com/abichinger) 於 [#3429](https://github.com/wailsapp/wails/pull/3429) 提供；[@abichinger](https://github.com/abichinger) 於 [#3431](https://github.com/wailsapp/wails/pull/3431) 提供；[@bruxaodev](https://github.com/bruxaodev) 於 [#3667](https://github.com/wailsapp/wails/pull/3667) 提供；[@OlegGulevskyy](https://github.com/OlegGulevskyy) 於 [#3674](https://github.com/wailsapp/wails/pull/3674) 提供

### 修正

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) 於   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) 於   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) 於   [#3477](https://github.com/wailsapp/wails/pull/3477)   由 [Atterpac](https://github.com/atterpac) 於   [#3320](https://github.com/wailsapp/wails/pull/3320)。   於 [#3306](https://github.com/wailsapp/wails/pull/3306)。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) 於   [#2750](https://github.com/wailsapp/wails/pull/2750)。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   [jaybeecave](https://github.com/jaybeecave) 於   [#3052](https://github.com/wailsapp/wails/pull/3052)。   [@pylotlight](https://github.com/pylotlight) 於   [PR](https://github.com/wailsapp/wails/pull/3039)   安裝。由 [@pylotlight](https://github.com/pylotlight) 於   [PR](https://github.com/wailsapp/wails/pull/3032) 新增   [PR](https://github.com/wailsapp/wails/pull/3145)   空格 — @leaanthony。   [thomas-senechal](https://github.com/thomas-senechal) 於 PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) 於 PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   附加視窗 [tw1nk](https://github.com/tw1nk) 於 PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) 於   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` 存在。   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   監聽器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) 於   [#3330](https://github.com/wailsapp/wails/pull/3330)   產生器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3334](https://github.com/wailsapp/wails/pull/3334)   產生器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) 於   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) 於   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) 於   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) 於   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) 於   [#3431](https://github.com/wailsapp/wails/pull/3431)   由 [@pekim](https://github.com/pekim) 於   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel)。   PR [#3466](https://github.com/wailsapp/wails/pull/3622)，由   [@5aaee9](https://github.com/5aaee9) 提交。   [@windom](https://github.com/windom/) 於   [#3636](https://github.com/wailsapp/wails/pull/3636)。   Windows，由 [@bruxaodev](https://github.com/bruxaodev/) 於   [#3691](https://github.com/wailsapp/wails/pull/3691)。

### 已變更

[mmghv](https://github.com/mmghv) 於   [#3611](https://github.com/wailsapp/wails/pull/3611)   支援 Ubuntu 24.04 LTS，由 [atterpac](https://github.com/atterpac) 於   [#3461](https://github.com/wailsapp/wails/pull/3461)   必須具有 `type="module"` 屬性。由   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   物件，且不會啟動 WML 系統。此變更旨在改善封裝性。如有需要，可呼叫新的 `WML.Enable` 方法手動啟動 WML 系統。隨附的 JS 執行階段指令碼仍會自動執行這兩項作業。由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   視窗物件作為預設匯出。現在無法再透過 ESM 具名匯入或命名空間匯入語法匯入個別方法。   API。部分方法已變更名稱或原型，具體如下：`Screen`   變為 `GetScreen`；`GetZoomLevel`/`SetZoomLevel` 變為 `GetZoom`/`SetZoom`；   `GetZoom`、`Width` 和 `Height` 現在會直接傳回值，而不再將值包裝在物件中。由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   已移除。使用 `-names` CLI 選項可切換回依名稱呼叫。   由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   先前依其所屬套件命名；現在則使用完整的 Go 匯入路徑，包括模組路徑。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   呼叫 `application.NewService`。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) 於   [#3574](https://github.com/wailsapp/wails/pull/3574)
