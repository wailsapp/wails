---
title: "變更記錄"
description: "Wails 的版本歷程與發行說明 v3"
slug: "changelog"
sourcePath: "changelog.md"
---

圖例：

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- 此專案的所有重要變更都會記錄在此檔案中。

格式以[維護變更記錄](https://keepachangelog.com/en/1.0.0/)為基礎，且本專案遵循[語意化版本](https://semver.org/spec/v2.0.0.html)規範。

- `Added`表示新增功能。
- `Changed`表示現有功能的變更。
- `Deprecated`表示即將移除的功能。
- `Removed`表示已移除的功能。
- `Fixed`表示錯誤修正。
- `Security`表示安全漏洞。

_/

/_   * 自動更新 *   發行說明會依據已合併的 PR 產生。請勿在一般貢獻 PR 中編輯此檔案，也不要重新建立已移除的`v3/UNRELEASED_CHANGELOG.md`。如需維護變更記錄，請參閱`AGENTS.md`。 _/

## [尚未發行]

## v3.2.0-nightly.20261102 - 2026-11-02

## 已變更
- 改用 DisableOverscroll，而非 DisableBounce，見 @lethusang 的[PR](https://github.com/wailsapp/wails/pull/6077)

## v3.2.0-nightly.20261101 - 2026-11-01

## 已新增
- 新增實驗性的無 CGO macOS 與 Linux 後端，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6521)
- 原生 YAML 建置、可攜式封裝與簽署，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6531)

## 已修正
- 修正各種v3夜間建置問題並清理相依套件，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6519)
- 修復v3發行驗證與測試資料分類，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6533)
- 解決跨平台發行驗證失敗問題，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6538)

## v3.1.0 - 2026-11-01

## 已新增
- 新增實驗性的無 CGO macOS 與 Linux 後端，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6521)
- 原生 YAML 建置、可攜式封裝與簽署，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6531)

## 已變更
- 避免複製已捨棄的桌面中繼資料，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6391)
- 不使用 fmt 格式化內容長度，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6393)
- 選取索引根目錄時不收集候選項目，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6404)
- 重複使用不可搜尋檔案的複製緩衝區，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6408)
- 避免配置用於偵測前綴的讀取器，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6410)
- 對已知副檔名略過 MIME 偵測讀取，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6414)
- 防止佇列增長與縮減來回震盪，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6415)
- 避免配置隨即捨棄的請求內容，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6419)
- 略過已停用的偵錯記錄準備作業，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6420)
- 延遲配置 MIME 偵測緩衝區，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6428)
- 減少超連結格式化配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6443)
- 移除項目符號格式化的暫時配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6449)
- 將 UUID 字串轉換減少至一次配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6453)
- 減少陳述式控制代碼的格式化配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6454)
- 消除名稱型 UUID 摘要的暫時配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6455)
- 避免複製尚未完成的命令輸出行，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6460)
- 在套件暫存期間以串流方式複製執行檔，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6472)
- 避免重新掃描尚未完成的輸出行，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6473)
- 重複使用已編譯的命令路由模式，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6475)
- 在探索執行階段時避免查詢中繼資料，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6482)
- 減少本地化備用頁面的配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6485)
- 避免為每個處理器重新編譯模式，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6490)

## 已修正
- 為系統匣實作 SetTooltip，見 @fan711 的[PR](https://github.com/wailsapp/wails/pull/6017)
- 從設定入口網站讀取深色模式，見 @TechHutTV 的[PR](https://github.com/wailsapp/wails/pull/6072)
- 允許 Window.SetZoom 低於1.0，見 @APshenkin 的[PR](https://github.com/wailsapp/wails/pull/6166)
- 依發出順序傳遞給每個 Go 接聽程式，見 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/6275)
- 在 Emit 中於 hookLock 鎖定期間讀取事件掛鉤，見 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/6276)
- 一律公開系統匣的 Menu 屬性，避免在 Run 之後建立的系統匣遭捨棄，見 @AlbinoGeek 的[PR](https://github.com/wailsapp/wails/pull/6284)
- 將 Run 之前傳送至主執行緒的工作保留到 Run 啟動，見 @AlbinoGeek 的[PR](https://github.com/wailsapp/wails/pull/6287)
- 在 GitHub 更新提供者中使用正確的資產 URL，見 @gavidroselj 的[PR](https://github.com/wailsapp/wails/pull/6304)
- WebView2 省略 Content-Length 時保留請求本文，見 @strraaa 的[PR](https://github.com/wailsapp/wails/pull/6326)
- 以 DIP 空間回報 macOS 與 GTK 螢幕（混合 DPI 邊界），見 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/6346)
- 不要將已關閉的通知回報為點擊，見 @bercedev 的[PR](https://github.com/wailsapp/wails/pull/6367)
- 終端機關閉時停止開發程序，見 @xmasdev 的[PR](https://github.com/wailsapp/wails/pull/6371)
- 在佐證檢查中，將變更記錄執行對應至 squash merge 提交，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6373)
- 將工作列重新啟動的防彈跳延遲維持一秒，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6386)
- 進入全螢幕時若視窗起始狀態即為全螢幕，仍保留尺寸限制，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6387)
- 不要調整全螢幕或最大化視窗的大小來符合新的限制，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6388)
- 不要在 liquid-glass 範例中於 Run 之前顯示對話方塊，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6389)
- 在 Run 之前要求顯示的對話方塊不再永遠等待，而是直接不執行任何動作，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6390)
- 桌面項目沒有命令時改用備援方案，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6392)
- 在不保留檔名的情況下偵測目前內容，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6394)
- 辨識大小寫混合的資產副檔名，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6395)
- 釋放已使用的內容偵測器緩衝區，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6396)
- 釋放已移除的 ringqueue 參照，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6397)
- 在選擇備援方案之前遵循標頭刪除操作，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6398)
- 在清送前提交備援回應標頭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6399)
- 只尋找名稱完全相符的索引檔案，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6400)
- 對不可搜尋的 HEAD 回應避免產生主體及完整讀取，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6401)
- 跨多次短讀取填滿 MIME 偵測前綴，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6402)
- 提交成功但為空的不可搜尋回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6403)
- 避免重複展開後的 URI 目標，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6405)
- 在欄位展開期間保留路徑，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6406)
- 遵循 XDG 桌面檔案搜尋目錄，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6407)
- 在空清送時提交偵測器標頭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6409)
- 略過遮蔽桌面項目的目錄，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6411)
- 保留在可執行檔路徑內找到的目標，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6412)
- 計時器觸發後釋放回呼擷取的內容，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6413)
- 避免 ldconfig 程式庫誤比對，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6416)
- 限制檔案系統程式庫的備援比對範圍，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6417)
- 消除快取失效時的資料競爭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6418)
- 遵循選用的記錄器設定，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6421)
- 避免修改共用的啟動 URL 及使用過時的連接埠，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6422)
- 請求取消時停止重試開發代理，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6423)
- 缺少內嵌資產時傳回404，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6424)
- 處理內嵌執行階段的 HEAD 請求，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6425)
- 路由服務時保留逸出的路徑，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6426)
- 正確排序超大數值的預發行版本識別碼，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6427)
- 儲存工具下載項目之前先拒絕 HTTP 錯誤回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6429)
- 在尋找桌面項目期間遵循啟動期限，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6430)
- 在開發代理重試之間保留請求主體，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6432)
- 在每條傳回路徑上關閉 HTTP 回應主體，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6434)
- 在平行搜尋程式庫時保留取消錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6435)
- 解碼前先拒絕失敗的 GitHub 標籤回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6436)
- 保留含方括號的字面模板路徑，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6437)
- 只安裝一次 wails3，並保留安裝程式錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6438)
- 忽略隱藏項目及註解區塊中的標頭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6439)
- 顯示查詢穩定版本時的錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6440)
- 輸出純文字訊息而不解讀格式，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6441)
- 避免載入 JSON null 後寫入時發生 panic，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6442)
- 讀取整個儲存區時傳回快照，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6444)
- 部分迭代資料列後傳播錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6446)
- 將解壓縮目的地路徑保留為字面值，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6447)
- 關閉及重新開啟時釋放預備陳述式控制代碼，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6448)
- 處理預備作業與關閉或重新開啟之間的競爭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6450)
- 正確處理備援索引頁面的 HEAD 請求，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6451)
- 要求完整的十六進位陳述式控制代碼，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6452)
- 並行處理時保留預設的記憶體內資料庫，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6456)
- 將 stderr 與已剖析的 HEAD 輸出分開，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6457)
- 從嵌入式套件組合產生執行階段資產，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6458)
- 偵測 TODO、FIXME 及 TBD 預留位置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6459)
- 驗證版本資訊 API 回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6461)
- 在狀態檢查中遵循任務目錄，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6462)
- 防止複製至同一檔案造成資料遺失，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6463)
- 將任務目錄納入快取識別資訊，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6464)
- 每次檢查僅探測連接埠一次，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6465)
- 避免快取項目解碼為 null 時發生 panic，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6466)
- 遞增數值預發行版本時不受整數範圍限制，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6467)
- 防止舊版檢查掩蓋快取未命中，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6468)
- 在執行器的基底目錄中執行預設任務，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6469)
- 在檔案總管命令中保留字面百分比符號，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6470)
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6471)中對齊前置條件與任務的命令目錄
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6474)中以產生的一般檔案取代符號連結
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6476)中在建置失敗時保留現有的 DMG
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6477)中保留加上引號的正規表示式及路徑中的反斜線
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6478)中在配置影像記憶體前驗證產生的圖示尺寸
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6479)中透過 shell 執行命令鏈與重新導向
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6480)中在建置失敗時保留先前的 Linux 套件
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6481)中使用平台專屬的 shell 備援執行方式
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6483)中支援以重複的路徑前導斜線提供資產
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6484)中辨識含有建置中繼資料的穩定版本
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6486)中準確回報資源產生失敗
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6487)中避免取代作用中的報告器時發生競爭條件
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6488)中將並行的快取檔案儲存作業序列化
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6489)中支援多個範本相等性候選項目
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6491)中正規化已儲存但為 null 的類別對應表
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6492)中顯示執行失敗的變數 shell 命令所產生的診斷資訊
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6493)中排除 shell 命令鏈，不將其納入隱含的 Go 快取
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6494)中將 Linux 類別持久化作業序列化
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6495)中快取含有以引號括住或逸出的空格之 Go 輸出
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6496)中在編碼失敗時保留 Mac 圖示
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6497)中不再捨棄無效的 Go 測試次數旗標
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6498)中讓新複製的檔案保留來源權限
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6519)中修正各種v3夜間建置問題並清理相依性
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6533)中修復v3版本驗證與測試資料分類
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6538)中解決跨平台版本驗證失敗

## v3.1.0-nightly.20261028 - 2026-10-28

## 變更
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6490)中重複使用已編譯的正規表示式模式，以加快變更記錄的產生速度
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6485)中減少提供本地化備援頁面時的記憶體配置
- 由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/6482)中避免不必要的中繼資料查詢，以加快 AppImage 執行階段的探索速度
- 重複使用已編譯的命令路由模式，以提升 wake 效能，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6475)
- 避免重複掃描尚未完成的輸出行，以提升效能，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6473)
- 在 MSIX 套件暫存期間以串流方式複製可執行檔，以減少記憶體用量並加快建置速度，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6472)
- 避免複製尚未完成的命令輸出行，以提升效能，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6460)
- 減少產生名稱型 UUID 時的記憶體配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6455)
- 減少格式化 SQLite 陳述式控制代碼時的記憶體配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6454)
- 將 UUID 字串轉換減少為僅配置一次記憶體，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6453)
- 減少格式化變更記錄項目符號條目時的記憶體配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6449)
- 減少格式化終端機超連結時的記憶體配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6443)
- 延遲探測 MIME 緩衝區，以減少提供資產時的記憶體配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6428)
- 停用偵錯記錄時，略過資產伺服器的偵錯記錄準備工作，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6420)
- 避免建立會被捨棄的要求內容，以減少處理 WebView 要求時的記憶體配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6419)
- 略過已知副檔名的 MIME 探測，以加快資產提供速度，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6414)
- 避免使用讀取器探測內容，以減少提供資產時的記憶體配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6410)
- 減少提供不可搜尋的資產檔案時的記憶體配置，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6408)
- 避免收集候選項目，以加快資產伺服器選取索引根目錄的速度，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6404)
- 不使用 fmt 格式化內容長度，以提升資產伺服器效能，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6393)
- 減少讀取桌面檔案總管中繼資料時的記憶體複製，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6391)

## 已修正
- 將來源檔案複製至新的目的地檔案時保留其權限，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6498)
- 變更最小或最大尺寸限制時，避免全螢幕及最大化視窗調整大小，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6388)
- 發布 wailsapp/wails 儲存庫排程的星標歷史記錄更新，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6500)
- 停止在 wake 中捨棄無效的 Go 測試計數旗標，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6497)
- 編碼失敗時保留 Mac 圖示，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6496)
- 修正含有加上引號或逸出空格的 Go 輸出發生建置快取未命中的問題，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6495)
- 在 Linux 上依序持久化通知類別，以防止並行寫入競爭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6494)
- 將 shell 命令鏈排除在隱含的 Go 建置快取之外，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6493)
- 修正已儲存資料為 null 時無法載入通知類別對應表的問題，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6491)
- 顯示執行失敗的變數 shell 命令所產生的診斷資訊，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6492)
- 修正存在多個相等候選項目時的 wake 範本比對問題，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6489)
- 將並行的快取檔案儲存作業序列化，以防止快取寫入損毀，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6488)
- 防止替換目前使用中的報告器時發生競爭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6487)
- 在 CLI 中準確回報資源產生失敗，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6486)
- 讓 CLI 能識別包含建置中繼資料的穩定版本，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6484)
- 提供請求路徑含有重複前導斜線的資產，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6483)
- 修正 wake 在各平台上的 shell 備援執行，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6481)
- 在建置失敗時保留先前的 Linux 套件，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6480)
- 透過 shell 處理命令鏈與重新導向，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6479)
- 在解析 wake 命令時保留帶引號的規則運算式及路徑中的反斜線，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6477)
- 在封裝建置失敗時保留現有的 DMG 檔案，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6476)
- 在繫結輸出中，以產生的一般檔案取代符號連結，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6474)
- 在 Linux 檔案總管命令中保留字面百分比符號，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6470)
- 對齊前置條件與工作命令的目錄，使其解析結果一致，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6471)
- 在執行器的基底目錄中執行預設工作，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6469)
- 防止舊版檢查隱藏快取未命中，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6468)
- 修正 CLI 的預發行版本編號，使其遞增不受整數溢位限制，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6467)
- 避免快取項目解碼為 null 時發生 panic，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6466)
- 讓 CLI 每次檢查僅探測連接埠一次，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6465)
- 將工作目錄納入快取識別資訊，以避免過時的 wake 結果，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6464)
- 防止封裝作業將檔案複製到自身時遺失資料，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6463)
- 讓狀態檢查遵循工作目錄設定，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6462)
- 偵測變更日誌項目中的 TODO、FIXME 與 TBD 預留文字，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6459)
- 在 CLI 中從內嵌套件組合產生執行階段資產，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6458)
- 在 git 命令中，將 stderr 與經過解析的 HEAD 輸出分開，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6457)
- 在並行存取時維持預設的記憶體內 SQLite 資料庫，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6456)
- 針對備援索引頁面的 HEAD 請求傳回正確回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6451)
- 處理 SQLite 陳述式準備作業與資料庫關閉或重新開啟之間的競爭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6450)
- 在 SQLite 資料庫關閉或重新開啟時釋放預備陳述式控制代碼，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6448)
- 使用 gosod 時保持解壓縮目的地路徑的字面形式，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6447)
- 在部分資料列反覆運算後，傳播 SQLite 查詢的錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6446)
- 在 kvstore 中，對整個儲存區的讀取傳回快照，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6444)
- 防止 kvstore 載入 JSON null 值後在寫入時發生 panic，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6442)
- 輸出純文字終端訊息，不解析格式動詞，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6441)
- 在更新程式查詢最新穩定版本時顯示錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6440)
- 忽略變更日誌註解區塊中的隱藏項目與標頭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6439)
- 僅安裝一次 wails3，並在更新程式中顯示安裝程式錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6438)
- 在 gosod 中保留含方括號之範本路徑的字面形式，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6437)
- 在更新程式中，於解碼前拒絕失敗的 GitHub 標籤回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6436)
- 在並行搜尋程式庫時保留取消錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6435)
- 在每條傳回路徑上關閉 HTTP 回應主體，以防止資源洩漏，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6434)
- 在開發代理重試請求時保留請求主體，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6432)
- 探索桌面檔案總管時遵循啟動期限，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6430)
- 儲存 AppImage 工具下載項目之前拒絕 HTTP 錯誤，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6429)
- 在 semver 比較中正確排序超大數值的預發行版本識別碼，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6427)
- 在資產伺服器中路由服務時保留逸出的路徑，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6426)
- 在資產伺服器中處理內含執行階段的 HEAD 請求，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6425)
- 針對缺少的內含資產傳回404，而非不正確的回應，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6424)
- 在請求取消時停止開發代理重試，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6423)
- 避免資產伺服器共用的起始 URL 發生異動，並避免使用過時的連接埠，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6422)
- 讓資產伺服器遵循選用的記錄器設定，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6421)
- 消除程式庫路徑快取失效處理中的資料競爭，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6418)
- 搜尋 ldconfig 輸出時避免誤判程式庫相符，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6416)
- 防止資產伺服器佇列在負載下成長及反覆縮減，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6415)
- 防彈跳計時器觸發後釋放回呼擷取內容，以避免保留過時的參照，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6413)
- 保留在可執行檔路徑內找到的檔案總管目標，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6412)
- 在檔案總管中略過遮蔽桌面項目的目錄，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6411)
- 修正回應主體為空時未傳送資產伺服器標頭的問題，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6409)
- 尋找應用程式圖示時遵循 XDG 桌面檔案搜尋目錄，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6407)
- 在檔案總管展開欄位時保留檔案路徑，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6406)
- 避免檔案總管重複顯示展開後的 URI 目標，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6405)
- 提交成功但內容為空且不可搜尋的資產回應，而非讓其維持待處理狀態，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6403)
- 在資產伺服器中跨越多次短讀取正確填滿 MIME 探測前綴，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6402)
- 避免資產伺服器針對不可搜尋的 HEAD 要求讀取完整回應本文，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6401)
- 資產伺服器解析目錄要求時，僅提供完全相符的索引檔案，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6400)
- 修正在不支援動作的伺服器上，已關閉的 Linux 通知被回報為點擊的問題，詳見 @bercedev 的[PR](https://github.com/wailsapp/wails/pull/6367)
- 終端機關閉時停止開發程序，詳見 @xmasdev 的[PR](https://github.com/wailsapp/wails/pull/6371)
- Windows 工作列重新啟動的防彈跳旗標在清除前保留一秒，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6386)
- 修正資產伺服器遺失後備回應標頭的問題，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6399)
- 資產伺服器選擇後備回應前會遵循標頭刪除操作，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6398)
- 釋放資產伺服器中已移除的環形佇列參照，以避免記憶體洩漏，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6397)
- 釋放資產伺服器中已使用的內容探測緩衝區，以減少記憶體用量，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6396)
- 提供資產時辨識大小寫混合的資產副檔名，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6395)
- 資產伺服器無須保留檔名即可正確偵測內容類型，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6394)
- 桌面項目沒有命令時改用預設值，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6392)
- 避免 liquid-glass 範例在非 macOS 平台上於應用程式啟動前顯示對話方塊，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6389)
- 在 Windows 上從受限制的大小進入全螢幕時，保留視窗大小限制，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6387)
- 修正在 Run 之前呼叫的對話方塊不回傳而永久阻塞的問題，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6390)

## 安全性
- 在配置影像前驗證產生的圖示大小，以防止過度使用記憶體，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6478)
- 使用版本資訊 API 回應前先加以驗證，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6461)
- SQLite 外掛程式要求完整的十六進位陳述式控制代碼，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6452)
- 限制檔案系統程式庫的後備比對，以防止載入非預期的程式庫，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6417)

## v3.1.0-nightly.20261027 - 2026-10-27

## 變更
- 將404頁面連結至v2.wails.io 上相符的 Wails v2文件頁面，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6366)
- 將v2文件網站移至v2.wails.io，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6365)
- 從 wails.io 提供v3文件，並從v2.wails.io 提供v2文件，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6364)
- 更新v3.0.0版本資訊並修正安全公告中的說法，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6363)

## 修正
- 依觸發順序循序傳遞 Go 事件監聽器回呼，詳見 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/6275)
- 一律公布系統匣的 Menu 屬性，避免在 Linux 上捨棄 Run 之後建立的系統匣，詳見 @AlbinoGeek 的[PR](https://github.com/wailsapp/wails/pull/6284)
- 將 Run 之前提交的主執行緒工作保留至應用程式啟動，詳見 @AlbinoGeek 的[PR](https://github.com/wailsapp/wails/pull/6287)
- 下載 GitHub 更新程式的發行資產時使用正確的資產 API URL，詳見 @gavidroselj 的[PR](https://github.com/wailsapp/wails/pull/6304)
- 修正 macOS 與 Linux 混合 DPI 設定中的螢幕邊界和位置，詳見 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/6346)
- 當 WebView2 省略 Content-Length 標頭時，在 Windows 上保留要求本文，詳見 @strraaa 的[PR](https://github.com/wailsapp/wails/pull/6326)
- 在鎖定期間建立掛鉤快照，以修正事件掛鉤觸發時的資料競爭，詳見 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/6276)
- 修正`SystemTray.SetTooltip()`在 Linux 上沒有作用的問題：現在會設定 StatusNotifierItem 工具提示說明（#6017）
- 在每週發行工作流程中執行變更記錄證據檢查，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6374)
- 修正發行變更記錄證據檢查，以將執行記錄對應至 squash merge 提交，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6373)
- 修正 Linux 深色模式偵測在深色桌面上回傳淺色，直到配色變更才恢復正常，以及`linux:SystemThemeChanged`在 GTK4 後端永遠不會觸發的問題。兩個後端現在共用同一套實作，並以`org.freedesktop.appearance`入口網站命名空間為基礎（#6072）
- 在有框線視窗（Windows、macOS、Linux）上，`Window.SetZoom`不再將低於1.0的倍率提高至1.0：WebView2 `ZoomFactor`、WKWebView `magnification`及 WebKitGTK 的縮放層級都接受任何正值，而`Zoom`視窗選項原本就會直接傳遞這類值。無框線視窗仍維持1.0下限，因為其拖曳和調整大小區域不會配合縮放進行調整（[#4590](https://github.com/wailsapp/wails/issues/4590)）；現在會忽略非正值、NaN 和 Inf 倍率，而不再將其提高至1.0 — @APshenkin

## v3.0.0 - 2026-10-26

## 變更
- 每夜預發行版本會以`vX.Y.0-nightly.YYYYMMDD`發布，預覽下一個每週發行版本；每週發行版本則會遞增次要版本（`vX.Y.0`），因此 Go 的`@latest`一律解析至每週發行版本或修補程式版本；`@wailsio/runtime` npm 套件版本一律與 Wails 發行版本一致，每夜版本使用`nightly` dist-tag
- 專案範本會將`@wailsio/runtime`鎖定至產生範本的 CLI 版本
- Wails v3需要 Go 1.26或更新版本（Go 1.26.10，或1.27.3及更新版本，包含`crypto/x509`安全性修正）

## 修正
- 修正`wails3 update cli`找不到v3發行版本的問題：它只考慮v2標籤，且從錯誤的路徑安裝
- 修正在某些 Windows 組建上繪製應用程式功能表列時，Windows 應用程式因`panic: SetTextColor failed`而當機的問題：功能表列現在會在有效的裝置內容上繪製，或改用預設繪製方式；`UAHMENUITEM.Position`現在與 C 的`int`一致，而且功能表繪製永遠不會發生 panic（#6311）
- 在等待重新啟動期間確保暫存更新安全：關閉期間啟動的下載或第二次重新啟動，不會刪除交換輔助程式正等待安裝的酬載，也不會與其爭用該酬載 — @julianstorer
- 修正當`pkg/w32`的初始化在`uxtheme.dll`的序號函式指標仍在使用時卸載該模組，導致 Windows 應用程式在`applyDarkMode`中當機的問題：現在會讓該模組在程序的整個存續期間保持載入（#6356）
- 在交換前驗證暫存的承載內容，並在捨棄備份前驗證已安裝的應用程式，讓損壞的更新能夠回復，而不是取代正常運作的應用程式，由 @julianstorer 貢獻

## 已移除
- 由 @github-actions[bot] 在最終翻譯同步完成後，於[PR](https://github.com/wailsapp/wails/pull/5736)中移除 Crowdin 整合

## 安全性
- 在所有框架中注入防護指令碼，阻擋 Linux（GTK3 與 GTK4）上來自跨來源子框架的保留`wails:`訊息，以緩解 GHSA-2p5f-2gx8-mh8m
- 將`github.com/klauspost/compress`更新至1.18.7（GHSA-259r-337f-4rfw），並將`golang.org/x/image`更新至0.45.0；同時使用 Go 1.26.10建置，以納入`crypto/x509`修正（GO-2026-6809）
- 在執行階段套件、設定精靈、範例與專案範本中，將`vite`更新至8.2.4或更新版本（GHSA-9jrq-w75r-8gcw、GHSA-rq7h-c2jc-7f22、GHSA-vfpm-58rq-9qcg、GHSA-r825-w8mm-5pxw、GHSA-587h-7wmj-h5pm、GHSA-698v-4r3x-3vcv）

### 已知限制
- **Linux**：由於 WebKitGTK 不會向處理常式提供任何框架資訊，子框架防護機制無法阻止*同源*子框架（例如`srcdoc`或`about:blank` iframe）透過其父框架存取原生訊息處理常式。跨來源框架，以及未設定`allow-same-origin`的沙箱框架都會遭到封鎖。請勿在 Wails 視窗中嵌入不受信任的同源內容。
- **Android（實驗性）：**`window.wails` JavaScript 橋接會公開給 WebView 中的每個框架，且不會檢查主框架或來源。嵌入式 iframe 中的內容（包括跨來源內容）可以呼叫繫結方法及傳送保留的`wails:`訊息。請勿在應用程式的 WebView 中嵌入不受信任的第三方內容。這與 macOS 和 iOS 不同，後兩者會忽略來自子框架的`wails:`訊息。預計在3.0.1中修正。
- <strong>所有平台：</strong>子框架對繫結方法的呼叫會透過 HTTP 傳送至`/wails/runtime`，而非原生訊息處理常式，因此此防護不會限制這些呼叫。請將任何嵌入的 iframe 視為可存取您繫結的 Go 方法。

## v3.0.0-beta.34 - 2026-10-20

## 已修正
- 修正 Windows webview 在視窗啟用時未重新取得焦點，以及 Tab 鍵被困在 webview 內、無法移至原生視窗再移回的問題（#3783）
- 修正`internal/webview2/pkg/webview2` COM vtable 結構（`ICoreWebView2_2`..`_27`及其他衍生介面）遺漏繼承的方法鏈，導致呼叫落在錯誤 vtable 槽位的問題（#6129）
- 修正 Windows 詢問對話方塊忽略自訂按鈕標籤且從不觸發按鈕回呼的問題；自訂按鈕現在會使用 TaskDialog 顯示（若 comctl32 v6無法使用，則回退至 MessageBox，並依位置對應回呼）（#5794）
- 修正`wails3 dev`搭配 Vite 時會停滯60秒的問題：自動前端就緒閘門現在使用 TCP 探測，而非傳送不含`Accept`標頭並收到404的 HTTP GET（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6301)中確保系統匣選單在開啟期間更新後仍可正常運作
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6302)中，讓執行階段操作能妥善處理暫時性的 WebView2 COM 錯誤

## v3.0.0-beta.33 - 2026-10-19

## 已新增
- 繫結產生器現在會輸出`bindings/events`模組，針對以編譯期常數名稱透過`application.RegisterEvent`註冊的每個可探索事件，匯出具型別的名稱常數（例如為`user:login`輸出`UserLogin`）；無需 Vite 外掛即可搭配`Events.On`/`Emit`使用（#5645）

## 已修正
- 修正 Windows 詢問對話方塊忽略自訂按鈕標籤且從不觸發按鈕回呼的問題；自訂按鈕現在會使用 TaskDialog 顯示（若 comctl32 v6無法使用，則回退至 MessageBox，並依位置對應回呼）（#5794）
- 修正`wails3 dev`搭配 Vite 時會停滯60秒的問題：自動前端就緒閘門現在使用 TCP 探測，而非傳送不含`Accept`標頭並收到404的 HTTP GET（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6301)中確保系統匣選單在開啟期間更新後仍可正常運作
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6302)中，讓執行階段操作能妥善處理暫時性的 WebView2 COM 錯誤

## v3.0.0-beta.32 - 2026-10-18

## 已新增
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6297)中，為 macOS 14+ 新增原生 macOS 選單區段標題；請參閱[文件網站](https://v3.wails.io)（已更新20個頁面）

## 已修正
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6299)中，讓 macOS 上的 SetRelativePosition 以螢幕工作區為相對基準
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6298)中，讓更新程式輔助工具在父程序逾時或取代失敗後，重新啟動未變更的目標
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6293)中，將版本化文件使用的 pnpm 鎖定至10版

## v3.0.0-beta.31 - 2026-10-18

## 已新增
- 公開`@wailsio/runtime/cancellable`子路徑匯出，讓`CancellablePromise`無需對根執行階段進入點求值即可匯入（#6153）

## 已變更
- 要求 Go 1.26或更新版本：v3模組現在宣告`go 1.26.0`，而由`wails3 init`和`wails3 service init`產生的專案也宣告相同的最低版本

## 已修正
- 修正使用 Go 1.27建置時繫結產生器輸出會改變的問題；當時`encoding/json`型別（例如`Marshaler`、`Unmarshaler`和`RawMessage`）會成為`encoding/json/v2`及`encoding/json/jsontext`型別的別名。現在會再次將它們呈現為`any`，且不會為`encoding/json/v2`輸出模型（#6292）
- 複製前先驗證傳給`wails3 init -t`的遠端範本 URL：僅接受`https://`、`http://`、`ssh://`、`git://`和`user@host:path` URL；值會在`--`之後傳給 git，且複製期間 git 僅限使用這些傳輸方式
- 拒絕 HTTP 執行階段傳輸中超出範圍的`x-wails-window-id`標頭值，而非無聲地將其截斷為另一個視窗的 ID
- 強化更新程式的封存檔解壓縮，防範 zip-slip：現在會使用`filepath.IsLocal`檢查項目路徑，並檢查合併後的路徑是否仍位於目標範圍內；同時新增測試，涵蓋 zip 與 tar.gz 封存檔中巢狀`..`、反斜線、絕對路徑及透過符號連結寫入的項目
- 修正`wails3 generate appimage`產生的 AppImage 在 Debian 和 Ubuntu 以外的系統上啟動時，因「Unable to spawn a new child process: … WebKitNetworkProcess」而當機的問題。隨附的 WebKit 現在會執行 AppImage 內附的輔助程序。WebKit 的 bubblewrap 沙箱仍會保持啟用
- 修正`-tags server`建置在 Android 和 iOS 上無法編譯的問題，並修正[#6271](https://github.com/wailsapp/wails/issues/6271)
- 修正 Linux 上暫存目錄位於另一個檔案系統（例如 tmpfs `/tmp`）時，更新程式無法交換的問題：現在會在已安裝二進位檔旁進行暫存，並跨檔案系統複製
- 修正更新程式嘗試取代 AppImage 的唯讀掛載點，而非`.AppImage`檔案的問題
- 修正在 macOS 和 iOS 上，保留的`wails:`視窗控制訊息會被接受自 iframe 的問題。現在只接受來自視窗主框架的這類訊息；來自 iframe 的其他訊息仍會連同其來源資訊傳遞至`RawMessageHandler`
- 停止在 macOS 上發送`WindowDidUpdate`和`WindowWillUpdate`：`NSWindowDidUpdateNotification`會在每次執行迴圈運行時觸發，因此頁面每秒會收到約120次空的`evaluateJavaScript:`呼叫；在 macOS 上，每次呼叫27.0都會消耗使用者點擊的暫時啟用狀態，導致`navigator.clipboard.writeText()`及其他需要啟用狀態的 API 對真實點擊回傳`NotAllowedError`而拒絕執行（[#6170](https://github.com/wailsapp/wails/issues/6170)），由 @APshenkin 貢獻
- 將 Windows 的預設 info.json 語言修正為 en-US（0409），讓檔案總管能顯示版本詳細資料（#6210）
- 修正在同時配備 NVIDIA 獨立 GPU 與 AMD 或 Intel 整合式 GPU 的筆記型電腦上，Linux DMA-BUF 轉譯器會遭停用的問題：在 Wayland 上，現在只有當 WebKit 使用 NVIDIA 進行轉譯時（應用程式已卸載至獨立 GPU，或合成器在其上執行）才會停用，而不是只要載入 nvidia 模組便停用；後者會使整合式 GPU 上的轉譯速度慢如蝸牛。X11 工作階段維持舊有行為，因為 Xorg 使用的轉譯 GPU 是由其組態設定
- 在[PR](https://github.com/wailsapp/wails/pull/6135)中加入分階段複製的備援機制，處理 Linux 更新程式跨裝置重新命名失敗的情況，由 @john-okeefe 貢獻
- 修正 Windows 上的`WebviewWindow.SetMenu`只有在設定選單列主題時才會重繪選單列的問題，修復[#6103](https://github.com/wailsapp/wails/issues/6103)
- 修正`WebviewWindow.SetMenu(nil)`在 Windows 上會引發 panic，而不像 macOS 上不執行任何操作的問題，修復[#6104](https://github.com/wailsapp/wails/issues/6104)
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中透過重試機制，在彼此競爭的推送之間保留變更記錄項目，由 @leaanthony 貢獻
- 移除參照從未隨模組發布之二進位檔的嵌入項目，修正`go mod vendor`在所有平台上都會因`pattern arm64/WebView2Loader.dll: no matching files found`而失敗的問題；修復[#5782](https://github.com/wailsapp/wails/issues/5782)及[#5376](https://github.com/wailsapp/wails/issues/5376)，由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6031)中貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6159)中防止 macOS 重複拆除視窗，由 @charlesbalderson 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6162)中處理系統匣選單的 WM_CONTEXTMENU，使其可在按一下滑鼠右鍵時開啟，由 @wzk2025 貢獻
- 修正全部五種`InvokeSync*`變體在分派的回呼發生 panic 時，會永遠阻塞呼叫端的問題；修復[#6107](https://github.com/wailsapp/wails/issues/6107)，由 @RALIST 在[PR](https://github.com/wailsapp/wails/pull/6108)中貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6236)中安全地重建已銷毀的 GTK3 子選單，由 @leaanthony 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6261)中實作供巢狀 worker 使用的 Windows worker 轉送相容層，由 @taliesin-ai 貢獻
- 修正 Windows 上的`OpenFile()`在`CanChooseDirectories(true)`時會忽略`PromptForMultipleSelection()`的問題：使用者現在可以像在 macOS 和 Linux 上一樣，一次選取多個資料夾（#6263）
- 修正 windows/arm64 上傳值給 WebView2 COM 方法的`double`引數問題——`ZoomFactor`、`RasterizationScale`、Cookie `Expires`，以及列印設定的縮放比例、頁面大小和邊界都受到影響（列印設定的 setter 先前會回傳`ErrDoubleArgUnsupported`）：受影響的呼叫現在會經由載入`d0`的`runtime.cgocall`跳板執行，而 Go 的 syscall 路徑無法載入該值（golang.org/issue/62583），由 @APshenkin 貢獻

## 安全性
- 伺服器模式（`-tags server`）現在只接受來自伺服器自身來源或`ServerOptions.WebSocketOriginPatterns`所列來源、以`POST`要求傳送的`/wails/runtime`呼叫；其他方法會收到 HTTP 405，其他瀏覽器來源則會收到 HTTP 403。若部署環境位於會重寫`Host`標頭的反向 Proxy 後方，應將其公開來源加入`WebSocketOriginPatterns`
- 繫結至迴路位址（預設的`localhost`、`127.0.0.0/8`或`::1`）的伺服器模式（`-tags server`），現在會拒絕`Host`標頭不是`localhost`、IP 常值、繫結主機或符合`ServerOptions.WebSocketOriginPatterns`之主機的要求，並回傳 HTTP 403，藉此阻擋針對`/wails/runtime`和 WebSocket 端點的 DNS 重新繫結攻擊
- Windows 通知在 WinRT 快顯通知 API 失敗時不再改用 PowerShell；現在會改為回傳錯誤

## v3.0.0-beta.30 - 2026-10-15

## 已修正
- 修正更新程式嘗試取代 AppImage 的唯讀掛載點，而不是`.AppImage`檔案的問題
- 停止在 macOS 上發送`WindowDidUpdate`和`WindowWillUpdate`：`NSWindowDidUpdateNotification`會在每次執行迴圈運行時觸發，因此頁面每秒會收到約120次空的`evaluateJavaScript:`呼叫；在 macOS 上，每次呼叫27.0都會消耗使用者點擊的暫時啟用狀態，導致`navigator.clipboard.writeText()`及其他需要啟用狀態的 API 對真實點擊回傳`NotAllowedError`而拒絕執行（[#6170](https://github.com/wailsapp/wails/issues/6170)），由 @APshenkin 貢獻
- 將 Windows 的預設 info.json 語言修正為 en-US（0409），讓檔案總管能顯示版本詳細資料（#6210）
- 在[PR](https://github.com/wailsapp/wails/pull/6135)中加入分階段複製的備援機制，處理 Linux 更新程式跨裝置重新命名失敗的情況，由 @john-okeefe 貢獻
- 修正 Windows 上的`WebviewWindow.SetMenu`只有在設定選單列主題時才會重繪選單列的問題，修復[#6103](https://github.com/wailsapp/wails/issues/6103)
- 修正`WebviewWindow.SetMenu(nil)`在 Windows 上會引發 panic，而不像 macOS 上不執行任何操作的問題，修復[#6104](https://github.com/wailsapp/wails/issues/6104)
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中透過重試機制，在彼此競爭的推送之間保留變更記錄項目，由 @leaanthony 貢獻
- 移除參照從未隨模組發布之二進位檔的嵌入項目，修正`go mod vendor`在所有平台上都會因`pattern arm64/WebView2Loader.dll: no matching files found`而失敗的問題；修復[#5782](https://github.com/wailsapp/wails/issues/5782)及[#5376](https://github.com/wailsapp/wails/issues/5376)，由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6031)中貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6159)中防止 macOS 重複拆除視窗，由 @charlesbalderson 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6162)中處理系統匣選單的 WM_CONTEXTMENU，使其可在按一下滑鼠右鍵時開啟，由 @wzk2025 貢獻
- 修正全部五種`InvokeSync*`變體在分派的回呼發生 panic 時，會永遠阻塞呼叫端的問題；修復[#6107](https://github.com/wailsapp/wails/issues/6107)，由 @RALIST 在[PR](https://github.com/wailsapp/wails/pull/6108)中貢獻

## v3.0.0-beta.29 - 2026-10-13

## 已新增
- macOS 上的全域快速鍵現在能感知鍵盤配置：在非 QWERTY 配置（AZERTY、QWERTZ、Dvorak……）上，字母快速鍵會透過`UCKeyTranslate`繫結至標示該字母的按鍵，而不是固定的 QWERTY 實體位置。

## 已修正
- 在[PR](https://github.com/wailsapp/wails/pull/6233)中，先初始化原生選單，再於 Linux 上設定該選單，由 @abhicodes-007 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6237)中解決 macOS 正式版就緒過程中，輸出處理與設定精靈的迴歸問題，由 @taliesin-ai 貢獻
- 修正 Windows 上的 WebView2 COM `AddRef`/`Release`包裝函式會將過時的執行緒最後錯誤回報為失敗的問題。它們現在會回傳參考計數，因此釋放參考不再可能中止應用程式或導致資產要求失敗（#3931）
- 修正`ICoreWebView2Deferral.Release`呼叫`AddRef`槽位，導致 Windows 上每個延後處理的資產要求都洩漏其延後物件的問題（#3931）

## v3.0.0-beta.28 - 2026-10-05

## 已修正
- 在[PR](https://github.com/wailsapp/wails/pull/6225)中，讓 macOS 應用程式輸出避開檔案提供者，由 @taliesin-ai 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6224)中完成設定精靈並顯示系統詳細資料，由 @taliesin-ai 貢獻
- 在[PR](https://github.com/wailsapp/wails/pull/6222)中改善 Windows 上 WebView2 復原作業的錯誤處理與資源清理，由 @leaanthony 貢獻
- 改善 Windows 上的 WebView2 復原機制，避免應用程式當機及控制器擱置，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6220)
- 修正 macOS 媒體授權指引，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6221)
- 修正以滑鼠右鍵按一下 Linux 系統匣圖示時，除了開啟選單還會執行點擊處理常式，導致每次按一下右鍵都切換所附加視窗的問題（#6018）
- 在 Windows 上維持 WebView2 回呼處理常式存活，直到原生端釋放為止，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6184)

## v3.0.0-beta.27 - 2026-10-01

## 已修正
- 修正 macOS 12及更新版本忽略`Permissions`選項的問題：現已實作`WKUIDelegate`媒體擷取方法，因此`PermissionAllow`和`PermissionDeny`會如同在 Linux 與 Windows 上一樣套用至攝影機和麥克風要求。仍須提供`NSCameraUsageDescription`／`NSMicrophoneUsageDescription`；若在沙盒中執行，也仍須提供相符的裝置權利項目

## v3.0.0-beta.26 - 2026-09-25

## 新增
- Watcher 對前端開發伺服器執行 HTTP 就緒檢查，詳見 @atterpac 的[PR](https://github.com/wailsapp/wails/pull/6150)

## 已修正
- 從 WebView2 處理程序失敗中復原，而非留下空白視窗，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6002)

## v3.0.0-beta.25 - 2026-09-22

## 已修正
- 主要應用程式結束時停止`wails3 dev`及其背景處理程序（#6048）

## v3.0.0-beta.24 - 2026-09-20

## 安全性
- 限制建置與發布工作流程中的 GITHUB_TOKEN 權限，以遵循最小權限原則
- 修正畫面範例資產中介軟體中的路徑周遊和符號連結逸出問題
- 在將設定精靈限制為僅能執行核准的套件管理員命令時，仍維持相依套件安裝功能正常運作

## v3.0.0-beta.23 - 2026-09-16

## 新增

- 允許在 Linux 上覆寫 GTK 應用程式 ID，詳見 @overlordtm 的[PR](https://github.com/wailsapp/wails/pull/5972)

## 已修正

- Linux 單一執行個體以應用程式本身的 ID 宣告匯流排名稱，詳見 @overlordtm 的[PR](https://github.com/wailsapp/wails/pull/5973)
- 在 Linux 上取消已中止的自訂配置要求，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6126)

## v3.0.0-beta.22 - 2026-09-14

## 新增

- 完成翻譯並移植社群媒體內容——請參閱[文件網站](https://v3.wails.io)（已更新1812個頁面），詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6122)
- 記錄各種樣式與大小的可嵌入 Wails 按鈕，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6120)

## v3.0.0-beta.21 - 2026-09-13

## 新增

- 使用 M-Press 提供 Wails v3文件，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6116)

## 已修正

- 解析 MPD 前置中繼資料中的 JSON slug 值以產生變更記錄，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6118)
- 更新程式會清除輔助程式環境變數，並在備份失敗後重新啟動原始目標，詳見 @cnmax 的[PR](https://github.com/wailsapp/wails/pull/6080)
- 在 App.Run 期間啟動預設訊號處理常式，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6098)
- Windows 選單現可處理 nil 選單、釋放被取代的資源，並重新繪製選單列，詳見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6112)
- 為使用共用 YAML 設定的新專案恢復 MSIX 封裝，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6115)
- 修正當泛型模型建立器參照稍後宣告的輔助項目時，產生的 JavaScript 和 TypeScript 繫結無法載入的問題，並防止建立相互依賴的泛型模型時發生堆疊溢位（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 已變更

- 將 Clave 展示連結更新至目前的網站與儲存庫，詳見 @01xR4in 的[PR](https://github.com/wailsapp/wails/pull/6082)

## 已修正

- 取消已中止的 Windows 資產要求（包括 worker 要求），同時在導覽期間保留 keepalive 處理常式。在 Apple 平台上，透過應用程式包裝函式轉送原生要求內容。（#5963、#5969）
- 透過重試，在相互競爭的推送之間保留變更記錄項目，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6094)
- 移除參照模組中從未隨附之二進位檔的嵌入項目，修正`go mod vendor`在所有平台上都因`pattern arm64/WebView2Loader.dll: no matching files found`而失敗的問題，並修正[#5782](https://github.com/wailsapp/wails/issues/5782)和[#5376](https://github.com/wailsapp/wails/issues/5376)，詳見 @Grantmartin2002 的[PR](https://github.com/wailsapp/wails/pull/6031)

## 已移除

- 移除已由純 Go 載入器取代的原生 WebView2 載入器支援。這會移除內嵌的`WebView2Loader.dll`二進位檔和`github.com/jchv/go-winloader`相依套件。`native_webview2loader`建置標籤仍可接受且不再報錯，但對v3建置沒有任何作用，詳見 @Grantmartin2002 的[PR](https://github.com/wailsapp/wails/pull/6031)
- 從 macOS API 指南移除未使用的建置標籤和 FPS 選項，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6097)

## v3.0.0-beta.19 - 2026-09-09

## 新增

- 以建置標籤限制私有 macOS API，讓使用者選擇啟用——請參閱[文件](https://v3.wails.io/features/browser/integration)、[文件](https://v3.wails.io/features/environment/info)、[文件](https://v3.wails.io/features/windows/basics)、[文件](https://v3.wails.io/features/windows/frameless)、[文件](https://v3.wails.io/features/windows/notch-windows)、[文件](https://v3.wails.io/features/windows/options)、[文件](https://v3.wails.io/guides/build/macos)、[文件](https://v3.wails.io/guides/build/private-macos-apis)及[文件](https://v3.wails.io/reference/overview)，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6087)

## 已修正

- 拒絕超過64 MiB 的執行階段要求，並回傳 HTTP 413，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6091)

## 安全性

- 透過權杖驗證強化 MCP 來源與遠端存取安全性，詳見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6092)

## v3.0.0-beta.18 - 2026-09-08

## 已修正

- 透過使用指標接收器修正 Linux 和 Darwin 上的 Calloc 記憶體洩漏，詳見 @4RH1T3CT0R7 的[PR](https://github.com/wailsapp/wails/pull/6083)

## v3.0.0-beta.17 - 2026-09-06

## 已修正

- Windows：WebResourceRequested 處理常式中失敗或為 nil 的`GetRequest`不再終止處理程序（`log.Fatal`／nil 解參照 panic）；現在會改為捨棄該要求並記錄日誌，詳見 @midagedev 的[PR](https://github.com/wailsapp/wails/pull/6006)

## v3.0.0-beta.16 - 2026-08-29

## 已變更

- 在新的終端機視窗中提示輸入公證密碼，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6029)

## 已修正

- 在 macOS 上正確處理系統匣點擊類型，見 @ChewbaccaCookie 的[PR](https://github.com/wailsapp/wails/pull/5919)
- CI 在更新前移除未使用的 Microsoft apt 軟體庫，見 @Grantmartin2002 的[PR](https://github.com/wailsapp/wails/pull/6041)

## v3.0.0-beta.15 - 2026-08-27

## 已修正

- 將 WebView2 內嵌逾時時間增加至60秒，見 @Grantmartin2002 的[PR](https://github.com/wailsapp/wails/pull/6043)

## v3.0.0-beta.14 - 2026-08-26

## 已修正

- 在 macOS 上正確命名 Control 字母按鍵操作，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6032)
- 修正 ICO 系統匣圖示，並在 Windows 上跟隨工作列主題，見 @nik9play 的[PR](https://github.com/wailsapp/wails/pull/6016)

## v3.0.0-beta.13 - 2026-08-25

## 已修正

- 在 macOS 執行模態迴圈期間，持續處理主執行緒工作，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6026)
- 讓行動裝置安全儲存作業可回報失敗，並在失敗時預設拒絕，見 @mortenolsrud 的[PR](https://github.com/wailsapp/wails/pull/5923)
- 即使未註冊監聽器，仍執行應用程式事件掛鉤，見 @archy-rock3t-cloud 的[PR](https://github.com/wailsapp/wails/pull/5999)
- 修正註解與本地化文件中的錯字，見 @haoku123 的[PR](https://github.com/wailsapp/wails/pull/6023)
- 移除提交於`v3/examples`下的預先編譯 macOS 二進位檔，見 @4RH1T3CT0R7 的[PR](https://github.com/wailsapp/wails/pull/6025)

## v3.0.0-beta.12 - 2026-08-21

## 新增

- 新增具備生命週期與遙測範例的 macOS 瀏海通知視窗——請參閱[文件](https://v3.wails.io/features/windows/notch-windows)，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6010)
- 新增具備新選項與原生整合的 macOS NSPanel 視窗支援——請參閱[文件](https://v3.wails.io/features/windows/options)，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6008)

## 已修正

- 避免並行呼叫導致 SQLite Prepare 停滯，見 @archy-rock3t-cloud 的[PR](https://github.com/wailsapp/wails/pull/5998)

## v3.0.0-beta.11 - 2026-08-20

## 已移除

- 從文件中移除過時的實作追蹤器，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6005)

## v3.0.0-beta.10 - 2026-08-19

## 已修正

- 修正 GTK4 Linux 主機會捨棄自訂通訊協定與檔案關聯啟動引數的問題，見 @midagedev 的[PR](https://github.com/wailsapp/wails/pull/6000)
- 在變更日誌驗證中正確處理已刪除的行與相同來源的修正，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5993)

## v3.0.0-beta.9 - 2026-08-16

## 新增

- 新增安全的 wails3 mcp 伺服器，以供代理程式輔助管理專案，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5896)
- 新增繫結中模型的文件——請參閱[文件](https://v3.wails.io/features/bindings/models)，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5988)
- 支援在不可變 Linux 系統上透過 rpm-ostree 安裝，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5987)
- 新增原生的每週星標歷史圖表產生與發布功能——請參閱[文件](https://v3.wails.io/credits)、[文件](https://v3.wails.io/de/credits)、[文件](https://v3.wails.io/fr/credits)、[文件](https://v3.wails.io/id/credits)、[文件](https://v3.wails.io/ja/credits)、[文件](https://v3.wails.io/ko/credits)、[文件](https://v3.wails.io/pt/credits)、[文件](https://v3.wails.io/ru/credits)、[文件](https://v3.wails.io/zh-cn/credits)、[文件](https://v3.wails.io/zh-tw/credits)，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5986)
- 新增僅適用於 Darwin、用來解析應用程式套件資源的 mac 套件——請參閱[文件](https://v3.wails.io/guides/build/macos)，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5965)
- 新增 Condui 展示頁面與索引項目——請參閱[文件](https://v3.wails.io/community/showcase/condui)及[文件](https://v3.wails.io/community/showcase)，見 @mgueregath 的[PR](https://github.com/wailsapp/wails/pull/5962)
- 新增包含螢幕擷取畫面與專案連結的 Redis Viewer 展示頁面——請參閱[文件](https://v3.wails.io/community/showcase)及[文件](https://v3.wails.io/community/showcase/redisviewer)，見 @redisviewer 的[PR](https://github.com/wailsapp/wails/pull/5984)

## 已變更

- 在 Linux 上將 GTK 應用程式旗標更新為 G<em>APPLICATION</em>NON_UNIQUE，見 @overlordtm 的[PR](https://github.com/wailsapp/wails/pull/5971)
- 將缺少視窗事件的記錄層級從警告改為偵錯，見 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5914)

## 已修正

- 讓已註冊的 macOS 鍵盤快速鍵優先於 webview 處理，見 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5902)
- 修復文件側邊欄中失效的連結，見 @northes 的[PR](https://github.com/wailsapp/wails/pull/5937)
- 當 WebKit 中止相符的自訂配置任務時，取消 macOS 與 iOS 資產請求的 context（#5963）
- 處理 WindowSetFullscreenButtonEnabled 訊息，見 @archy-rock3t-cloud 的[PR](https://github.com/wailsapp/wails/pull/5976)
- 在 preact-ts 範本中匯入 Fragment，以解決建置失敗問題，見 @haoku123 的[PR](https://github.com/wailsapp/wails/pull/5979)
- 避免舊版僅提供服務的 GTK3 應用程式，在尚無可用的使用中視窗或顯示器便執行螢幕探索時當機（#5966）
- 即使未發布的變更日誌為空，也允許明確指定版本的發布作業繼續進行（#5977）

## 安全性

- 將網站的 nanoid 鎖定檔更新至已修補的3.3.18，以解決安全性公告指出的問題，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5985)

## v3.0.0-beta.8 - 2026-08-12

## 新增

- 在自動產生的變更日誌項目中加入文件 URL 產生功能，見 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5957)
- 新增 Streams：使用 WebSocket 程式設計模型，在 Go 與 JavaScript 之間提供雙向位元組串流，且不需要監聽通訊端。在 Go 中以`app.HandleStream(name, handler)`宣告串流，並從前端以`Stream(name)`連線；後者會傳回形如`WebSocket`的物件。Go→JS 透過資產伺服器為每個視窗保留一個輪詢來傳輸，JS→Go 則透過一般 POST 傳輸；不會繫結任何 TCP 連接埠，也不會有任何資料經過`evaluateJavaScript`。在伺服器建置（`-tags server`）中，相同的處理常式會改由真正的 WebSocket 提供服務，因此各種建置使用的應用程式程式碼完全相同。作者為 @leaanthony
- 將信箱變更日誌項目移至「未發布」，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5935)

## 已變更

- 更新文件側邊欄自動產生功能與部落格作者型別推導，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5938)

## 已修正

- WebView2 初始化使用截止期限與訊息泵，見 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5952)
- WebView2 Cookie 測試在 CI 中除非明確啟用，否則會略過，並將執行鎖定於目前的作業系統執行緒；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5951)中完成
- Windows 選單建構器恢復子選單父項目的命令 ID；由 @gilad-ch 於[PR](https://github.com/wailsapp/wails/pull/5944)中完成
- 使官方交叉編譯映像與 GTK 4.14+ Linux 支援基準一致（#5928）
- 設定 iOS Xcode 專案以保留繼承的連結器旗標並加入 -ObjC；由 @mortenolsrud 於[PR](https://github.com/wailsapp/wails/pull/5915)中完成
- 修正大型前端在`wails3 dev`資產代理中過度頻繁地建立及關閉 TCP 連線的問題；此問題可能耗盡主機的暫時連接埠，導致不相關的程序因`EADDRNOTAVAIL`而失敗
- 將各視窗的事件 JavaScript 排入佇列，以依序分派並提供反壓機制；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5934)中完成

## 已移除

- 移除桌面二進位檔發布管線：v3版本僅發布標籤，而`wails3` CLI 透過`go install`安裝。刪除`release-v3.yml`以及分派該工作流程的每夜建置步驟；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5946)中完成

## v3.0.0-beta.7 - 2026-08-11

## 新增

- 新增 macOS 自動播放偏好設定，以停用媒體播放需要使用者操作的要求；由 @Eyalm321 於[PR](https://github.com/wailsapp/wails/pull/5512)中完成
- 將信箱變更記錄項目移至「尚未發布」；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5935)中完成

## 已變更

- macOS 縮放動畫改用 CADisplayLink 或 NSTimer，使效能更流暢；由 @savely-krasovsky 於[PR](https://github.com/wailsapp/wails/pull/5945)中完成

## 已修正

- 設定 iOS Xcode 專案以保留繼承的連結器旗標並加入 -ObjC；由 @mortenolsrud 於[PR](https://github.com/wailsapp/wails/pull/5915)中完成
- 修正大型前端在`wails3 dev`資產代理中過度頻繁地建立及關閉 TCP 連線的問題；此問題可能耗盡主機的暫時連接埠，導致不相關的程序因`EADDRNOTAVAIL`而失敗
- 將各視窗的事件 JavaScript 排入佇列，以依序分派並提供反壓機制；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5934)中完成

### 新增

- 實作泛型非同步 FIFO 信箱，以依序傳遞事件；由 @savely-krasovsky 與 @DevLumuz 於[PR](https://github.com/wailsapp/wails/pull/5851)中完成

## v3.0.0-beta.6 - 2026-08-09

## 新增

- 針對過大的事件實作有界的主機端儲存空間，並依序傳遞 JavaScript；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5930)中完成
- 實作 macOS Dock 彈跳效果以提示視窗；由 @julianstorer 於[PR](https://github.com/wailsapp/wails/pull/5921)中完成

## 已修正

- 資產伺服器在清除緩衝區時會保留內容類型偵測錯誤及尚未寫入的前置資料；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5931)中完成
- 避免從 Wails 回呼取代應用程式選單時造成 macOS 應用程式當機
- 修正 Windows 10 1809／Windows Server 2019（組建 17763）上的原生選單無法閱讀問題。深色模式的 uxtheme 匯出函式原本僅在組建 18334上啟用，因此應用程式層級的深色模式明確啟用程序從未在這些主機上執行：選單背景被繪製成深色，但 Windows 仍以淺色佈景主題繪製選單文字，造成深色背景上顯示深色文字。這些序號自17763起即已存在，因此現在的條件限制已與其相符。
- 修正`w32.GetStockObject`呼叫`GetDeviceCaps`而非`GetStockObject`的問題；該問題會使其針對每個預設物件都傳回0。
- 改善 WebView2 啟動載入器的下載錯誤處理與回報；由 @jannskiee 於[PR](https://github.com/wailsapp/wails/pull/5924)中完成

## v3.0.0-beta.5 - 2026-08-07

## 已修正

- macOS 應用程式啟用現在只會對一般應用程式遵循啟用原則；由 @julianstorer 於[PR](https://github.com/wailsapp/wails/pull/5897)中完成
- 在 Linux 組建中防護尚未初始化的 GTK 視窗；由 @julianstorer 於[PR](https://github.com/wailsapp/wails/pull/5898)中完成
- 在載入 URL 前，為 Linux WebKit 視窗明確設定不透明背景色彩；由 @julianstorer 於[PR](https://github.com/wailsapp/wails/pull/5899)中完成

## v3.0.0-beta.4 - 2026-08-05

## 已變更

- Android 建置工作預設使用 arm64，而 deploy-emulator 會選取主機架構；由 @mortenolsrud 於[PR](https://github.com/wailsapp/wails/pull/5890)中完成

## 已修正

- 拖曳期間保留 macOS 視窗的縮放狀態並減少動態效果；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5900)中完成
- 將`!server`加入`webview_window_windows_nonclient.go`建置條件約束，以修正 Windows 伺服器模式組建

## v3.0.0-beta.3 - 2026-08-03

## 新增

- 在實作詳細資料中記錄第10階段 Beta 驗證已完成；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5881)中完成

## 已修正

- 將視窗控制代碼傳遞給 Windows 深色模式 API 並驗證引數；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5877)中完成
- 集中處理 macOS 無框線視窗的標題列按鈕狀態解析；由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5870)中完成
- 避免 Windows 應用程式要求使用深色模式、但 Windows 應用程式佈景主題為淺色時，原生選單文字無法閱讀。現在選單會使用相符的淺色原生背景，直到 Windows 能夠呈現深色選單文字為止。
- 修正 Windows 10 1809／Windows Server 2019（組建 17763）上的原生選單無法閱讀問題。深色模式的 uxtheme 匯出函式原本僅在組建 18334上啟用，因此應用程式層級的深色模式明確啟用程序從未在這些主機上執行：選單背景被繪製成深色，但 Windows 仍以淺色佈景主題繪製選單文字，造成深色背景上顯示深色文字。這些序號自17763起即已存在，因此現在的條件限制已與其相符。

## v3.0.0-beta.2 - 2026-08-02

## 已變更

- 將v3從 Alpha 升級為 Beta
- 記錄系統匣的智慧預設值與快顯視窗自動隱藏行為，並加入點擊處理常式選取的迴歸測試涵蓋範圍（#5840）。
- GitHub 更新程式預設排除 Windows 安裝程式資產；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5861)中完成
- 新增 macOS 無框線視窗的圓角、直角及自訂半徑邊角支援；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5866)中完成

## 已修正

- 回報即時 GTK4 視窗大小，並從已設定的表面發出調整大小、最大化、最小化及全螢幕狀態事件（#5830）。
- 修正在 fetch 要求中傳送 Blob 或 FormData 時造成 Linux WebKit 當機的問題；由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5854)中完成
- 缺少 Blob/FormData 標頭時，Fetch 墊片會傳遞 undefined；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5865)中完成

## v3.0.0-alpha2.122 - 2026-08-01

## 新增

## 已變更

- 新增 macOS 無框線視窗的圓角、直角及自訂半徑邊角支援；由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5866)中完成

## 已修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5865)中讓 fetch shim 對缺少的 Blob/FormData 標頭傳遞 undefined

## v3.0.0-alpha2.121 - 2026-07-31

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5857)中新增 macOS DMG 封裝支援，包括新的選項與建置工作

## 變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5861)中讓 GitHub 更新程式預設排除 Windows 安裝程式資產

## 修正

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5854)中修正 Linux WebKit 在 fetch 請求中傳送 Blob 或 FormData 時當機的問題

## v3.0.0-alpha2.120 - 2026-07-31

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5853)中實作 macOS 標題列按兩下時的最大化或最小化動作

## 變更

- 由 @jeongkyu 在[PR](https://github.com/wailsapp/wails/pull/5849)中更新 QR 服務教學以使用 NewServiceWithOptions，並加入間距

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5856)中讓 WKWebView 在 macOS 縮放期間保持回應
- 修正 GTK4 視窗大小查詢，並從已設定的`GdkSurface`發出調整大小、最大化、最小化及全螢幕狀態事件。

## v3.0.0-alpha2.119 - 2026-07-27

## 修正

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5833)中更新多種語言的文件，加入架構圖

## v3.0.0-alpha2.118 - 2026-07-26

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5825)中為圖示產生作業的輸入與輸出提供預設路徑
- 由 @savely-krasovsky 在[PR](https://github.com/wailsapp/wails/pull/5797)中將來源進入點模組加入執行階段 package.json 的 sideEffects
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5816)中為貢獻指南新增授權與來源追溯章節

## 修正

- 由 @savely-krasovsky 在[PR](https://github.com/wailsapp/wails/pull/5800)中套用限定範圍的 GTK4 無邊框 CSS，以移除圓角半徑
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5789)中讓 Windows 的快顯選單與螢幕列舉妥善處理游標位置取得失敗
- 由 @phergul 在[PR](https://github.com/wailsapp/wails/pull/5678)中讓 macOS 開啟檔案對話方塊正確篩選副檔名，並依檔名後綴驗證允許的檔案
- 由 @roachadam 在[PR](https://github.com/wailsapp/wails/pull/5793)中防止 Windows 深色模式初始化時呼叫 nil API
- 修正更新程式在32位元建置時失敗的問題：在`GOARCH=386`上將`maxArchiveTotalSize`常數（2 GiB）傳給`fmt.Errorf`時，超出平台`int`的範圍。現在已將其明確指定為`int64`型別。
- 修正 Windows 建置版本啟動時發生 nil 指標 panic 的問題。當視窗使用深色（或系統深色）標題列，而該 Windows 建置版本未載入深色模式 uxtheme API（例如 Windows 10 1809 / Windows Server 2019〔組建17763〕）時，便會發生此問題。視窗佈景主題設定中的`AllowDarkModeForWindow`呼叫現在已加入 nil 防護，與`w32.SetMenuTheme`中既有的防護一致。

## v3.0.0-alpha2.117 - 2026-07-08

## 新增

- 由 @savely-krasovsky 在[PR](https://github.com/wailsapp/wails/pull/5462)中為 Windows 的非用戶端區域實作自訂點擊測試邏輯

## 變更

- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5761)中根據 UseVisualHosting 設定 WebView2 顯示器縮放比例偵測

## v3.0.0-alpha2.116 - 2026-07-07

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5763)中更新常見問題文件，使其著重於 Wails v3的功能與指引

## v3.0.0-alpha2.115 - 2026-07-06

## 修正

- 修正`Menu.Update()`未在 GTK4 Linux 上重建原生選單的問題（#5659；@puneetdixit200 也在 #5539中獨立診斷並修正此問題）
- 透過複製螢幕 ID／名稱字串並擷取數量快照，修正顯示器變更時列舉 macOS 螢幕會當機的問題（#5565；@x-haose 也在 #5584中獨立診斷並修正此問題）
- 修正 Windows 上的當機問題：在最小化／還原轉換期間，`GetClientRect`傳回 nil，而`WM_ERASEBKGND`繪製純色背景時會觸發此問題（防護措施由 @sinspired 在 #5636中回報）
- 由 @mbaklor 在 #5690中修正前端繫結錯誤一律解析為文字的問題
- 修正使用`server`建置標籤時 Windows 建置失敗的問題。原因是 Windows GUI 檔案缺少 macOS 與 Linux 對應檔案已有的`!server`建置條件（#5680）

## v3.0.0-alpha2.114 - 2026-07-05

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5720)中實作更新資訊清單通訊協定與端點提供者

## 變更

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5711)中將`webview2`繫結以`v3/internal/webview2`形式整合至v3模組，移除獨立模組、其每夜版本發布／同步工作流程，以及 go.mod 版本調整機制（v3是其唯一使用者）

## 修正

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5750)中將 WebView2 顯示器縮放比例偵測及 DPI 變更時的主機重新同步修正移至「未發布」章節
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5741)中更新 WebView2 對 float64 與 BOOL 參數的 COM 封送處理
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5703)中避免 Windows 系統匣圖示更新與銷毀時發生 panic 及 nil 解參照
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5743)中修正隱藏視窗在 Windows 上未正確再次隱藏的問題
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5742)中讓 WebView2 控制器的可見性與視窗最小化、最大化及還原狀態同步

### 修正

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5734)中重新啟用 WebView2 顯示器縮放比例偵測，並僅在 DPI 變更時重新同步主機；此變更以 @randalmurphal 驗證的修正為基礎，由 @eleclin 驗證根本原因，並由 @qq540491950 進行硬體測試

## v3.0.0-alpha2.113 - 2026-07-04

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中，新增在未設定`ANDROID_KEYSTORE_FILE`時建置發布版 AAB 的警告（Google Play 會拒絕以偵錯金鑰簽署的套件），並記錄 App Bundle 封裝與簽署方式
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5739)中新增多種語言的「為何選擇 Wails」文件
- 由 @fbbdev 在[PR](https://github.com/wailsapp/wails/pull/5398)中支援在繫結中將 Go time.Time 對應至 JS Date 或 string
- 新增用於提交至 Play Store 的 Android App Bundle（AAB）封裝工作（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）；APK 工作仍保留供本機／模擬器測試使用，由 @mortenolsrud 於[PR](https://github.com/wailsapp/wails/pull/5728)中實作（修正[#5726](https://github.com/wailsapp/wails/issues/5726)）
- 新增 Android 實體裝置工作目標，並恢復相機／位置權限，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5735)中實作

## 變更

- 將`webview2`升級至v1.0.28（[版本資訊](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- 將 Android 範本的`compileSdk`／`targetSdk`從34升級至35；這是 Google Play 對新應用程式提交的要求，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5730)中實作

## 修正

- 修正 sponsorkit 中預先烘焙的頭像遮罩，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5745)中實作
- 修正 Android AVD 自動建立功能因按字典順序排序版本，而選取錯誤系統映像或 cmdline-tools 版本的問題，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5730)中實作
- 修正設定精靈建議過時 Android NDK 版本的問題（現為26.3.11579264，與文件記載的需求一致），由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5730)中實作
- 更新 SvelteKit 與選項的法文文件，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5744)中實作
- 修正 macOS 在顯示器變更期間列舉螢幕時發生的 SIGSEGV，由 @flofreud 於[PR](https://github.com/wailsapp/wails/pull/5516)中實作

## v3.0.0-alpha2.112 - 2026-07-03

## 新增

- 新增以 Go 為基礎的貢獻者 SVG 產生器，並更新 docs／website 的致謝頁面，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5724)中實作
- 支援在繫結中將 Go time.Time 對應至 JS Date 或字串，由 @fbbdev 於[PR](https://github.com/wailsapp/wails/pull/5398)中實作

## 變更

- 以 Go 產生器取代以 Node 為基礎的贊助者圖片處理管線，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5719)中實作

## 修正

- 修正 Android 建置資產相依套件安裝指令碼，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5729)中實作
- 在`ValidateAndSanitizeURL`中拒絕 U+0085（NEXT LINE）控制字元，使 URL 驗證器完整涵蓋空白字元
- 在無框線視窗的 DPI 變更時重新計算 DWM 框架，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/4785)中實作
- 修正 Windows 使用非100%縮放比例時，DnD 放置區域偵測失敗的問題，由 @yulesxoxo 於[PR](https://github.com/wailsapp/wails/pull/4632)中實作
- 為 Darwin 對話方塊、選單、系統匣及通知中的 Cocoa 物件新增明確的 Objective-C 記憶體管理，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5714)中實作
- 修正 Linux CGO 後端錯誤與系統匣問題，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5718)中實作

## v3.0.0-alpha2.111 - 2026-07-01

## 新增

- 將 HappyTools 加入社群展示，由 @Aliuyanfeng 於[PR](https://github.com/wailsapp/wails/pull/5061)中實作
- 新增印尼文語系支援與完整文件，由 @triadmoko 於[PR](https://github.com/wailsapp/wails/pull/5643)中實作
- 為 WindowsWindow 新增 DisableMenu 選項，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/4813)中實作

## 變更

- 更新 Taskfile 範本與 CLI，以使用 GOOS 和 ARCH 分派建置／封裝工作，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5617)中實作

## 修正

- 修正 Mac 視窗分頁問題，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5708)中實作

## 移除

- 從貢獻、功能與指南中移除德文翻譯的 MDX 檔案，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5702)中實作

## v3.0.0-alpha2.110 - 2026-06-30

## 新增

- 實作 macOS WebView 重新載入與強制重新載入，並新增 WebContent 處理程序終止後的復原機制，由 @wayneforrest 於[PR](https://github.com/wailsapp/wails/pull/5129)中實作
- 新增涵蓋貢獻、功能與指南的完整德文文件，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5396)中實作
- 強化通知功能，加入音效、附件、排程與更新 API，由 @popaprozac 於[PR](https://github.com/wailsapp/wails/pull/5333)中實作

## 修正

- 在無框線視窗的 DPI 變更時重新計算 DWM 框架，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/4785)中實作
- 修正 Windows 使用非100%縮放比例時，DnD 放置區域偵測失敗的問題，由 @yulesxoxo 於[PR](https://github.com/wailsapp/wails/pull/4632)中實作

## v3.0.0-alpha2.109 - 2026-06-29

## 新增

- 在 EventsEmit 文件中新增程式碼範例，由 @iamhabbeboy 於[PR](https://github.com/wailsapp/wails/pull/5026)中實作
- 新增 Windows WebView2 視覺化託管選項，由 @MerIijn 於[PR](https://github.com/wailsapp/wails/pull/5380)中實作
- 將 Klustr 加入社群展示文件，由 @SametKUM 於[PR](https://github.com/wailsapp/wails/pull/5536)中實作
- 將 Kira 加入社群展示，並新增頁面與變更記錄項目，由 @thiennguyen93 於[PR](https://github.com/wailsapp/wails/pull/5685)中實作
- 在 MCP 服務指南中新增意見回饋章節，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5694)中實作

## 變更

- 伺服器模式現在具備正式的一級生產環境建置流程，與桌面建置工作一致（#5693）。`task build:server`預設會建置生產環境二進位檔（`-tags server,production`、`-trimpath`、已移除符號），並接受`DEV=true`（開發伺服器）、`OBFUSCATED=true`（garble）及`EXTRA_TAGS`。`task run:server`會執行開發伺服器。`Dockerfile.server`／`task build:docker`會先建置生產環境伺服器（`-tags server,production`）及生產環境前端；映像預設是在 distroless/static 上建置的純 Go 靜態版本，並將`CGO_ENABLED`、`GO_IMAGE`及`RUNTIME_IMAGE`公開為可覆寫的建置引數，以供 CGO 應用程式使用。

## 修正

- 避免關閉仍有待處理非同步呼叫的視窗時發生當機，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/4435)中實作
- 避免在 Windows 上開啟隱藏的應用程式時啟用視窗，由 @leaanthony 於[PR](https://github.com/wailsapp/wails/pull/5249)中實作
- 確保 WebKit 的要求中繼資料、回應完成處理及本文串流處理皆在 GTK 主執行緒上執行，由 @taliesin-ai 於[PR](https://github.com/wailsapp/wails/pull/5668)中實作
- 修正`Menu.Update()`未在 GTK4 Linux 上重新建置原生選單的問題（#5659；@puneetdixit200 亦於 #5539中獨立診斷並修正此問題）
- 透過複製螢幕 ID／名稱字串並建立數量快照，修正在顯示器變更時列舉 macOS 螢幕所發生的當機（#5565；@x-haose 亦於 #5584中獨立診斷並修正此問題）
- 透過在`WM_DPICHANGED`處理常式中重新設定控制器邊界，比照取消最小化時的 DPI 重新同步機制，修正在 Windows 上將視窗拖曳跨越不同 DPI 的顯示器後，WebView2 內容先縮小再消失的問題（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 新增

- 透過`app.GlobalShortcut`新增全域（系統範圍）鍵盤快速鍵（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）。即使應用程式未取得焦點，快速鍵仍會觸發。各平台皆採原生實作，且不含任何第三方相依套件：macOS 使用 Carbon hot keys、Windows 使用`RegisterHotKey`、X11 使用`XGrabKey`，Wayland 則使用 XDG Desktop Portal 全域快速鍵介面。
- 新增內建 MCP 伺服器：使用`mcp`標籤建置應用程式時，Model Context Protocol 伺服器會自動啟動，讓 LLM 代理程式測試及控制執行中的 Wails 應用程式，包括視窗控制、DOM 檢查、JavaScript 求值、繫結方法呼叫、事件，以及以螢幕上的動畫游標呈現的模擬滑鼠／鍵盤輸入。無須撰寫使用者程式碼：設定`WAILS_MCP=1`後，`wails3 build`／`wails3 dev`會自動加入`mcp`標籤。完全透過環境變數設定（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）。

## 已修正

- 修正`Menu.Update()`未在 GTK4 Linux 上重建原生選單的問題（#5659；@puneetdixit200 亦在 #5539中獨立診斷並修正）
- 透過複製螢幕 ID／名稱字串並建立數量快照，修正顯示器變更時列舉 macOS 螢幕所造成的當機（#5565；@x-haose 亦在 #5584中獨立診斷並修正）

## v3.0.0-alpha2.107 - 2026-06-27

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5613)中新增實驗性 Wake 文件及其側邊欄導覽

## v3.0.0-alpha2.106 - 2026-06-24

## 變更

- 將`webview2`升級至v1.0.27。
  - ci(webview2)：修正發布版本建置（交叉編譯 Windows 並補全 go.sum）（#5671）\

  **完整差異：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5672)中，從 webview2 發布工作流程的交叉編譯移除 go vet
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5670)中，將自動變更記錄使用的 OpenRouter 模型更新為 google/gemini-2.5-flash-lite
- 將`webview2`升級至v1.0.26。

### 修正

- **從暫時性的執行階段 COM 錯誤復原，而非結束程式**（#5658、#5580）。`Chromium.errorCallback`先前會對*任何* COM 錯誤呼叫`os.Exit(1)`，因此啟動後一次可復原的短暫異常就會終止整個應用程式。執行階段路徑（`Resize`／`GetClientRect`、`Navigate`／`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）現在會記錄錯誤並復原。尤其是`MessageReceived`中格式錯誤／不受信任的 Web 訊息，現在會被捨棄，而不會使處理程序終止。這可處理跨越混合 DPI 顯示器時的當機類型（#5544、#5650）。環境／控制器建立路徑的錯誤仍是致命錯誤。\

**完整差異：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 已修正

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5671)中修正 release-webview2 工作流程，使其能正確處理 go.sum 檔案
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5659)中，透過清除並重建原生選單修正 Linux GTK4 選單更新

## v3.0.0-alpha2.105 - 2026-06-21

## 新增

- 新增`application.System`，供共用程式碼在執行階段偵測平台：`System.IsMobile()`（iOS／Android）、`System.IsDesktop()`（macOS／Windows／Linux）、`System.IsServer()`（`server`建置標籤），以及用來直接測試單一目標的`System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`。它可在所有目標上編譯，因此無須建置標籤即可建立條件分支。`@wailsio/runtime`亦提供對應的前端輔助函式（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）
- 新增「使用其他前端框架」指南，說明如何將自己的 Vite 專案放入`frontend/`（涵蓋 Solid、Preact、Lit、SvelteKit、Qwik、Angular 等）
- `wails3 setup`精靈現在會檢查行動平台（iOS／Android）工具鏈，包括 Xcode 與 iOS Simulator 執行階段、JDK、Android SDK／NDK 和模擬器；在適用處並提供一鍵安裝及可複製的 shell 設定修正方式
- 產生的專案隨附`frontend/.npmrc`，其中設定7天的`minimum-release-age`，以降低接觸剛發布（且可能已遭入侵）套件的風險（pnpm 和 bun 會遵循此設定；npm 會忽略，但不會造成影響）

## 變更

- 以全新的霓虹山景主視覺重新設計所有內建入門範本（Web、iOS 和 Android）
- **TypeScript 現在是入門範本的預設語言，並使用不含語言後綴的範本名稱。** `wails3 init`（不含`-t`）會建立 TypeScript 專案的骨架；`-t vanilla`、`-t react`、`-t vue`和`-t svelte`均為 TypeScript，JavaScript 變體則為`-t vanilla-js`、`-t react-js`、`-t vue-js`和`-t svelte-js`。內建範本會在`template.yaml`中以`typescript:`宣告其語言；使用`-ts`後綴的社群範本仍可作為備援方式運作
- 以霓虹「digital Wails」主題重新設計`wails3 setup`精靈（山景背景上具有毛玻璃般的鮮明效果）

## 已修正

- 修正在 Windows 上還原應用程式時的當機；此情況發生於應用程式最小化時間過長，導致 WebView2 暫停，或其算繪／GPU 處理程序遭回收。最小化／還原時的 DPI 重新同步（#5544）現在只會在視窗 DPI 確實變更時操作 WebView2 控制器，避免在常見的相同 DPI 還原情況下，對已暫停的控制器進行致命 COM 呼叫（#5605）
- 修正在頻繁載入資產／媒體的長時間執行 Linux 應用程式中，原生`SIGABRT`／`SIGSEGV`反覆當機的問題（通常發生在 GTK 主迴圈期間的`g_object_unref`內）。資產伺服器先前會從工作 goroutine 完成`WebKitURISchemeRequest`，因而在 GTK 主執行緒外呼叫非執行緒安全的 WebKit2GTK 函式；現在完成作業（`webkit_uri_scheme_request_finish_with_response`／`finish_error`）會在主執行緒上執行。這完成了 #5566中的部分修正。GTK3 和 GTK4／WebKitGTK 的6.0建置皆受影響（#5631、#5557）
- 修正 Linux／GTK3 上`setupSignalHandlers`中間歇出現的`fatal error: invalid pointer found on stack`。以訊號`user_data`傳遞的視窗 ID 原本儲存在 Go 的`unsafe.Pointer`區域變數中，因此垃圾回收器在堆疊複製期間掃描到該（非指標）值時會中止。現在 Go 端會將 ID 保持為整數型別（`uintptr_t`），把 #4958對 GTK4 路徑套用的相同修正反向移植至舊版 GTK3 路徑（該修正將 C 訊號函式改為`uintptr_t`，以消除`-race`／checkptr 錯誤）（#5631）

## 移除

- 移除`react-swc`、`preact`、`lit`、`solid`、`qwik`和`sveltekit`入門範本（以及它們的`-ts`變體）。目前支援的內建集合為`vanilla`、`react`、`vue`和`svelte`；每個範本預設皆使用 TypeScript，並提供`-js` JavaScript 變體。仍可透過[使用自備前端](https://v3.wails.io/guides/dev/frontend-frameworks)或自訂範本來使用任何其他框架

## v3.0.0-alpha2.104 - 2026-06-18

## 已修正

- 修正繫結的 Go 服務方法傳回空字串時發生的 iOS 當機（SIGABRT）。iOS 資產回應寫入器先前以`buf != nil`而非長度檢查本文指標，因此長度為零的本文會使`&buf[0]`發生 panic；現在改為檢查長度，與桌面平台寫入器一致

## v3.0.0-alpha2.103 - 2026-06-15

## 變更

- 將 iOS 和 Android 原生功能移至平台管理器：改為透過`application.IOS.*`和`application.Android.*`呼叫（例如`application.IOS.Haptic("medium")`、`application.Android.Share(payload)`），而不再使用舊的`application.IOS*`／`application.Android*`獨立函式（#5602）
- 重新命名行動裝置橋接事件：跨平台事件現在使用`common:*`前綴（例如`common:haptic`、`common:location`），平台專屬事件則使用`ios:*`／`android:*`（例如`ios:backgroundTask`、`android:foregroundService`）；不再使用`native:*`前綴（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 新增

- 新增實驗性的`wails3 setup`精靈，用於互動式專案設定與相依性檢查
- 為`wails3 doctor`新增`--json`旗標，以提供機器可讀的輸出
- 在`wails3 doctor`命令中新增簽署狀態區段

## 修正

- 修正 Linux 上的 npm 偵測，除了套件管理器外也會檢查 PATH

## v3.0.0-alpha.101 - 2026-06-13

## 新增

- iOS：原生訊息對話方塊（UIAlertController），以及開啟檔案／多個檔案／目錄的對話方塊（UIDocumentPickerViewController）；儲存對話方塊會明確傳回錯誤
- iOS：透過 UIPasteboard 支援剪貼簿
- iOS：透過 UIScreen 取得真實螢幕度量（點、像素、縮放比例、安全區域工作區）
- iOS：裝置組建（`IOS_PLATFORM=device`）、程式碼簽署身分／佈建描述檔／權利設定支援、`.ipa`封裝，以及透過 devicectl 執行`deploy-device`
- iOS：可設定最低 iOS 版本（build/config.yml 中的`ios.minIOSVersion`）
- iOS：`wails3 doctor`會回報 macOS 上 Xcode 與 iOS SDK 的可用狀態
- iOS：系統事件——電池、網路、主題、螢幕鎖定及低記憶體狀態會以`events.IOS.*`和平台中立的`events.Common.*`應用程式事件呈現
- iOS：原生行動功能橋接（匯出的`application.IOS*`）——分享面板、開啟 URL、保持喚醒、手電筒、安全區域內距、亮度、應用程式資訊、方向鎖定、狀態列、生物辨識（Face ID／Touch ID）、本機通知及 Keychain 安全儲存空間
- iOS：感測器與硬體——觸覺回饋、單次地理定位、加速度計、接近感測器、文字轉語音、儲存空間資訊、電源／電池狀態、網路狀態、鍵盤內距及螢幕擷取偵測
- iOS：文件（IOS.md 和文件網站指南）
- Android：原生訊息對話方塊（AlertDialog），以及開啟檔案／多個檔案的對話方塊（Storage Access Framework，匯入為快取副本）；開啟目錄和儲存對話方塊會明確傳回錯誤
- Android：透過 ClipboardManager 支援剪貼簿
- Android：透過 WindowMetrics／DisplayMetrics 取得真實螢幕度量（dp、像素、縮放比例、系統列工作區）
- Android：觸覺回饋（`Android.Haptics.Vibrate`）、裝置資訊（`Android.Device.Info`）和快顯通知（`Android.Toast.Show`）執行階段方法
- Android：具型別的生命週期事件（`events.Android.*`，由 events.txt 產生），其中`ActivityCreated`會對應至`Common.ApplicationStarted`
- Android：組建管線會產生可安裝的偵錯與發行版 APK（`android:run`、`android:package`、`android:package:fat`）；發行版預設使用偵錯金鑰庫簽署，也可透過`ANDROID_KEYSTORE_*`環境變數使用正式金鑰庫
- Android：`wails3 doctor`會回報 Android SDK、NDK 和 JDK
- Android：系統事件——電池、網路、主題、螢幕鎖定及低記憶體狀態會以`events.Android.*`和平台中立的`events.Common.*`應用程式事件呈現
- Android：原生行動功能橋接（匯出的`application.Android*`）——分享、開啟 URL、保持喚醒、手電筒、安全區域內距、亮度、應用程式資訊、方向鎖定、狀態列、生物辨識（BiometricPrompt）、本機通知及 EncryptedSharedPreferences 安全儲存空間
- Android：感測器與硬體——觸覺回饋、單次地理定位、加速度計、接近感測器、文字轉語音、儲存空間資訊、電源／電池狀態、網路狀態、鍵盤內距，以及使用 FLAG_SECURE 阻擋螢幕擷取
- Android：文件（ANDROID.md 和文件網站指南）
- 範例：`mobile`綜合範例新增 Mobile 和 Hardware 分頁，示範跨 iOS 與 Android 的原生功能橋接（膠囊式分頁可換行為多列）
- 行動裝置：電池——應用程式進入背景時會暫停加速度計、接近感測器、手電筒和範例中的週期性時鐘，返回前景時再恢復（Android 會讓程序繼續在背景執行，而手電筒在 iOS 上是會持續保留的硬體狀態）；此外，只有當應用程式位於前景時，才會註冊 Android 系統事件接收器
- iOS：相機擷取——`application.IOSCapturePhoto`／`IOSCaptureVideo`（UIImagePickerController → 帶有 base64 縮圖的`native:capture`事件）
- iOS：背景執行——`application.IOSBeginBackgroundTask`／`IOSEndBackgroundTask`（UIApplication 背景工作時間窗），以及可設定的`ios.backgroundModes`（build/config.yml），它會將`UIBackgroundModes`套用至產生的 Info.plist
- Android：相機擷取——`application.AndroidCapturePhoto`／`AndroidCaptureVideo`（透過 FileProvider 使用系統相機 → `native:capture`事件）
- Android：前景服務——`application.AndroidStartForegroundService`／`AndroidStopForegroundService`（帶有持續通知的`WailsForegroundService`可讓程序持續執行，以處理長時間執行的背景工作）
- 範例：新增 Camera 分頁，示範相片／影片擷取及背景執行（Android 使用前景服務，iOS 使用背景工作時間窗）

## 修正

- 修正`getUserMedia`在 Linux 上一律因`NotAllowedError`而失敗的問題：WebKitGTK 會拒絕未由任何處理常式處理的權限要求，而`permission-request`訊號先前並未連接。現在會依照新的跨平台`WebviewWindowOptions.Permissions`對應表（`map[PermissionType]Permission`）處理相機／麥克風權限，Linux（WebKitGTK）和 Windows（WebView2）都會遵循此設定。在沒有原生提示的 Linux 上，相機／麥克風預設為允許（恢復`getUserMedia`），並可使用`PermissionDeny`停用（#5552）
- iOS：`GOOS=ios`已可再次編譯（匯出`events.IOS`、加入行動方法名稱虛設常式），帶有 production 標籤的組建也可編譯（修正 pkg/application 和數個服務中的組建標籤）
- iOS：Go→JS 事件和 ExecJS 現已正常運作——頁面不再於啟動時載入兩次，`wails:runtime:ready`交握也不再可能遺失
- iOS：`ApplicationDidFinishLaunching`／`ApplicationStarted`不再與應用程式啟動產生競爭；已移除固定2秒的啟動等待
- iOS：修正每次執行 Go→JS JavaScript 時發生的 C 字串洩漏
- iOS：`hasListeners`現在會反映實際的監聽器註冊狀態
- iOS：正式環境組建在編譯時會排除框架偵錯記錄
- Android：`GOOS=android`已可再次編譯——定義`events.Android`、移除越界的`events_android.go`監聽器陣列、新增行動方法名稱虛設常式，並防止桌面版 Linux 檔案（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）混入 Android 組建
- Android：JS→Go 繫結現已正常運作——WebView 無法將`fetch()` POST 本文傳遞至`shouldInterceptRequest`，因此執行階段呼叫現在改由 JavascriptInterface 傳輸（`nativeHandleRuntimeCall`）路由，不再因要求本文為 nil 而當機
- Android：`Screens.*`執行階段呼叫現在會傳回真實資料——ScreenManager 現在會在啟動時填入資料（先前從未接線，因此`GetAll`會傳回 nil）
- Android：正式環境組建在編譯時會排除框架偵錯記錄；偵錯組建則會使用`Wails`標籤將記錄傳送至 logcat
- Android：實作真正的`hasListeners`登錄機制、JNI 參照／例外處理，以及只載入一次的頁面生命週期（不再重複導覽）
- 修正當 Vite 開發伺服器正在執行時，`wails3 generate bindings`在 Windows 上因「存取遭拒」而失敗的問題：現在會將產生的檔案同步至輸出目錄，而非以重新命名方式覆蓋該目錄（#5515）
- 修正在 macOS 上顯示器變更後讀取螢幕資訊時偶發的嚴重當機：儲存的螢幕 ID 和名稱指標指向已自動釋放的`UTF8String`緩衝區，可能在 Go 複製前就已釋放（釋放後使用）。現在會對字串執行`strdup`，並在轉換後將其釋放；螢幕列舉也會在明確的自動釋放集區中執行，因此從 Go goroutine 呼叫時不再造成記憶體洩漏（#5556）
- 修正 Linux 上 assetserver 關閉`WebKitURISchemeRequest`時偶發的 SIGSEGV：最後一個`g_object_unref`在 assetserver goroutine 上執行，導致 WebKit GObject 在 GTK 主執行緒之外完成終結處理。現在會透過`g_main_context_invoke`將 unref 操作封送至 GTK 主內容（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 新增

- 擴充`MacWebviewPreferences`，加入其他 WKWebView 設定選項：`EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize`及`ApplicationNameForUserAgent`（#5549）

## 修正

- 修正在 Windows 上執行 Vite 開發伺服器時，`wails3 generate bindings`因「存取遭拒」而失敗的問題：現在會將產生的檔案同步至輸出目錄，而非透過重新命名覆蓋該目錄（#5561）
- 修正 Linux 上無邊框視窗不會觸發 JS 調整大小事件的問題；修正無邊框視窗的捲軸邊緣偵測（#5368）
- 修正 Windows 更新程式在暫存目錄與安裝目錄位於不同磁碟區時，因「無效的跨裝置連結」而失敗的問題（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 修正

- 修正在 Windows 上執行 Vite 開發伺服器時，`wails3 generate bindings`因「存取遭拒」而失敗的問題：現在會將產生的檔案同步至輸出目錄，而非透過重新命名覆蓋該目錄（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 修正

- 修正 Linux 上 WebKit 閒置時（例如開啟檢查器時）UI 凍結的問題：不再對`SIGUSR1`強制套用`SA_ONSTACK`，因為這會破壞 JavaScriptCore 的 GC 執行緒同步（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 新增

- 新增偵錯頁面及使用`runtime/trace`的說明

## 變更

- 移除一些不必要的`_ "embed"`匯入，稍微整理程式碼

## 修正

- 修正 Windows 上視窗取消最大化後未強制套用最小寬度／高度限制的問題（#4593）
- 修正使用 Frameless + Transparent 視窗選項時，全螢幕模式下滑鼠點擊會穿透的問題（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 新增

- 新增 Garble 混淆支援（[#4563](https://github.com/wailsapp/wails/issues/4563)）：穩定的繫結方法 ID、建置／Taskfile 串接（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`），以及所有面向執行階段之承載資料上的 JSON 結構標籤（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`），使線上格式不受 Garble 重新命名匯出欄位的影響。

## v3.0.0-alpha.95 - 2026-05-20

## 新增

- 補上缺少的專案結構頁面

## 變更

- 文件：將架構頁面上的幾張圖表改為循序圖，以獲得更清晰的顯示效果
- 文件：加入執行前必須先安裝 D2 的附註

## 修正

- 修正 GTK4 預設組態上的`wails3 generate appimage`：打包工具現在會先從二進位檔偵測 GTK 技術堆疊，再搜尋執行階段檔案，因此 GTK4 建置會選用`libwebkitgtkinjectedbundle.so`（位於`webkitgtk-6.0/`之下），`-tags gtk3`建置則會選用`libwebkit2gtkinjectedbundle.so`（位於`webkit2gtk-4.1/`之下）。`.relr.dyn`探測現在也會檢查`libgtk-4.so.1`，因此無論使用哪個技術堆疊，現代工具鏈都能正確停用符號剝除。（#5475）
- 修正使用相對`-builddir`叫用`wails3 generate appimage`時失敗的問題：打包工具現在會預先將`-binary`、`-icon`、`-desktopfile`、`-builddir`及`-outputdir`解析為絕對路徑，因此流程中途的`s.CD`不會再中斷 AppRun 下載 goroutine 或複製後的`ldd`探測。
- 修正桌面檔案的`Name=`欄位與二進位檔基本名稱不符時，`wails3 generate appimage`無法將最終 AppImage 移至`-outputdir`的問題：打包工具現在會透過`OUTPUT`環境變數，強制 linuxdeploy 的 appimage 外掛程式將 AppImage 寫入`<binary>-<arch>.AppImage`，而非使用從桌面檔案推導出的名稱。
- 修正在 alpha.93 中將 GTK4 + WebKitGTK 6.0技術堆疊提升為預設值後，Linux 上不會觸發`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`及`Common.SystemDidWake`的問題。新的預設`application_linux.go` `run()`未呼叫`setupCommonEvents()`（此函式會將`Linux.*`事件轉送至其對應的`Common.*`事件）或`monitorPowerEvents()`。現在透過`application_linux_dbus.go`，讓 GTK3 與 GTK4 建置路徑共用 DBus 電源監視輔助程式。（#5474）

## v3.0.0-alpha.94 - 2026-05-19

## 修正

- 修正在 alpha.93 中將 GTK4 + WebKitGTK 6.0技術堆疊提升為預設值後，Linux 上不會觸發`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`及`Common.SystemDidWake`的問題。新的預設`application_linux.go` `run()`未呼叫`setupCommonEvents()`（此函式會將`Linux.*`事件轉送至其對應的`Common.*`事件）或`monitorPowerEvents()`。現在透過`application_linux_dbus.go`，讓 GTK3 與 GTK4 建置路徑共用 DBus 電源監視輔助程式。（#5474）

## v3.0.0-alpha.93 - 2026-05-17

## 新增

- 由 @leaanthony 在 Linux 的`wails3 doctor`輸出中加入`XDG_SESSION_TYPE`

## 修正

- 由 @leaanthony 修正 Wayland 上 appmenu-gtk-module 存取尚未實現的視窗所導致的視窗選單當機（#4769）
- 由 @leaanthony 修正應用程式名稱包含無效字元（空格、括號等）時 GTK 應用程式當機的問題
- 由 @overlordtm 修正 Windows 上初始化拖放功能時出現「記憶體不足」錯誤的問題（#4701）
- 由 @leaanthony 修正 mainthread 回呼儲存區使用錯誤的 RLock 刪除 map 項目所造成的競爭條件（Linux、macOS、iOS）（#4424）
- 修正將命令列引數傳遞給任務時的變數處理。現在會正確初始化以 KEY=VALUE 配對指定的 CLI 變數，並在整個任務執行期間傳播這些變數。
- 修正 macOS 上的 NSWindowZoomButton 衝突：`MaximiseButtonState`與`FullscreenButtonState`現在會在啟動及執行階段套用限制較嚴格的狀態；任一 setter 都無法再默默覆寫另一個（#5319）
- 修正 CodeRabbit 在 #5463 中發現、舊版 GTK3 建置路徑（`-tags gtk3`）原已存在的一系列錯誤：透過檔案關聯啟動時不再略過啟動處理常式；`getTheme`現在具備邊界與型別安全性；`appName`不再釋放由 GLib 擁有的記憶體；`clipboardGet`不再洩漏 GTK 傳回的`gchar*`；`Calloc`現在使用指標接收器（且`NewCalloc`會傳回`*Calloc`），讓集區能確實追蹤配置；`zoomOut`現在使用`zoomInFactor`的倒數，而非會被箝制為1.0的負乘數；`execJS`會重複使用預先配置的空世界名稱，不再每次呼叫都洩漏一個`C.CString("")`；並已從`menuItem.setAccelerator`移除開發用的`fmt.Println`。解決 #5465。
- 修正預設 GTK4 建置路徑（`linux_cgo.go`）中相同的 `Calloc` 值接收器洩漏：改用指標接收器與 `NewCalloc() *Calloc`，以確實追蹤並釋放每個視窗的 `c.String(...)` 配置。

## v3.0.0-alpha.92 - 2026-05-15

## 新增

- 修改 Taskfile，以允許透過 `PACKAGE_MANAGER` 選項控制使用的前端套件管理器
- 以 `{{.Opn}}` 和 `{{.Cls}}` 擴充範本資料，讓 Taskfile 範本的撰寫更可預期

## 變更

- 修改數個現有 Taskfile，以使用 `{{.Opn}} and {{.Cls}}`

## 修正

- 修正面板讀取系統匣選單時更新該選單，會在 `linuxSystemTray` 中觸發 `concurrent map read and map write` 執行階段嚴重錯誤的問題。
- WebView2 錯誤與堆疊追蹤輸出改用 `log`，而非 `fmt`，避免應用程式在 Windows 上未連接主控台執行時遺失訊息。

## v3.0.0-alpha.91 - 2026-05-12

## 變更

- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5414) 中更新贊助者 SVG
- <strong>重大變更（macOS）：</strong>統一 macOS 座標系統，讓 `GetScreens`、`Position` 和 `SetPosition` 全部使用相同的座標空間：以邏輯點為單位、Y 軸向下，且 `(0,0)` 位於主要螢幕左上角。這與 Windows、GTK，以及 Electron 和 Web 的公開 API 一致。實體位置在主要螢幕上方的螢幕現在會回報負的 `Bounds.Y`（先前為正值），而 `Position()`／`SetPosition()` 值現在使用邏輯點，而非 `points × primaryScale`。仍可維持 `Position()` → `SetPosition()` 的往返轉換；從先前 alpha 建置記錄的絕對值，或手動計算的因應措施（例如乘以 `primaryScale`，或依螢幕高度翻轉 Y 軸）都需要更新。解決 [#5117](https://github.com/wailsapp/wails/issues/5117)。

## 修正

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5416) 中加入防禦性驗證 DBus 訊號名稱與訊息本文長度，以防止發生 panic
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5363) 中修正 Linux 上 GTK 選單處理的記憶體安全問題
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5295) 中加入 NVIDIA GPU 偵測，並在 Linux 上停用 DMA-BUF 轉譯器
- 在 [#5117](https://github.com/wailsapp/wails/issues/5117) 中修正 macOS 上 `SetPosition` 的跨螢幕 Y 座標轉換：以主要螢幕高度作為全域參考，讓視窗在相對於主要顯示器有垂直偏移的顯示器上落於正確位置
- 由 @wayneforrest 在 [PR](https://github.com/wailsapp/wails/pull/5109) 中修正 Git PR 範本，使其指向正確的意見回饋 URL
- 修正一系列 Windows 系統匣 `SetMenu` 當機問題。這些問題源於損壞的 `DestroyMenu` 系統呼叫傳入四個引數而非一個，導致每次呼叫都回傳 FALSE，且未釋放任何資源。此外，在重建選單時釋放 HMENU 和 HBITMAP 控制代碼（包括執行階段透過 `MenuItem.SetBitmap` 配置的控制代碼）、重設 `Win32Menu.Update` 中過時的核取方塊／選項按鈕對應，並移除 `systemtray.updateMenu` 中重複配置記憶體的多餘 `Update()` 呼叫。長時間執行的系統匣應用程式不再於每次重建選單時洩漏 GDI／USER 物件。

## v3.0.0-alpha.90 - 2026-05-11

## 新增

- 由 @vinhvoit225 在 [PR](https://github.com/wailsapp/wails/pull/5261) 中新增可設定的應用程式名稱，用於 macOS 上的 WKWebView User-Agent
- 由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5400) 中，為 gin-service 範例新增間接相依套件 github.com/coder/websocket
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5402) 中為建置資產測試新增深度相等比較支援

## 變更

- 由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5401) 中將建置輸出整合至 assets 目錄
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5399) 中更新贊助者 SVG

## 修正

- 由 @overlordtm 在 [PR](https://github.com/wailsapp/wails/pull/5289) 中，將通知物件用於 macOS 單一執行個體訊息
- 由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5383) 中批次處理 Windows 回呼，防止高負載下遺失 Promise

## v3.0.0-alpha.89 - 2026-05-10

## 新增

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5316) 中新增 go<em>test</em>results 工作，以彙整 Go 測試結果

## 變更

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5369) 中依條件將大型 RPC 承載資料拆分成分塊 POST 請求
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5386) 中，將所有前端範本的 Vite 從 5.x.x 升級至 8.0.0
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5365) 中將 Vite 開發伺服器的連接埠設定移轉至環境變數
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5361) 中，將所有範本的 Vite 開發伺服器設定為繫結至 127.0.0.1
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5384) 中更新贊助者 SVG

## 修正

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5312) 中，於 build-assets 更新期間清理 Info.plist 範本存根
- 在主執行緒同步套用選單項目修改方法（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`），修正 macOS 選單中的過時狀態，消除快速重新開啟選單時仍會呈現先前狀態的 `dispatch_async` 競爭情況（#5002）
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5203) 中，於開發模式忽略 `*_test.go` 檔案，以避免不必要的重新建置
- 由 @wucm667 在 [PR](https://github.com/wailsapp/wails/pull/5291) 中，防止應用程式未執行時 Menu.Update() 發生記憶體區段錯誤
- 由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5382) 中，使用 lastSizeWParam 控制 Windows 上的選單列重繪

## v3.0.0-alpha.88 - 2026-05-09

## 變更

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5371) 中將 HiddenOnTaskbar 改為使用 WS<em>EX</em>TOOLWINDOW
- 由 @atterpac 在 [PR](https://github.com/wailsapp/wails/pull/5370) 中重新排列相依套件，並移除 go.mod 中的 webview2 replace 指令
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5358) 中更新贊助者 SVG

## 修正

- 由 @fbbdev 在 [PR](https://github.com/wailsapp/wails/pull/5331) 中移除泛型間接別名，並整合 map 鍵型別

## 移除

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5377) 中刪除 PR-master 工作流程，並移除文件、Go 測試與略過測試

## v3.0.0-alpha.87 - 2026-05-07

## 新增

- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5352) 中新增 Wails v3 的韓文文件
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5354) 中新增安裝與快速入門的法文文件
- 由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5355) 中新增快速入門、概念與社群的葡萄牙文文件

## v3.0.0-alpha.86 - 2026-05-06

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5328)中新增法文文件本地化
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5343)中為文件網站新增德文語系

## 變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5347)中，於文件設定註冊所有8翻譯語系
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5317)中更新多個與 Windows WebView2 相關的檔案

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5340)中，將 Linux 的對話方塊分派拆分為 GTK3 與 GTK4 版本
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5339)中確保對話方塊回呼於 GTK 執行緒上執行，修正記憶體區段錯誤

## v3.0.0-alpha.85 - 2026-05-05

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5179)中將 PR 範本 URL 新增至儲存庫
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5330)中新增 Wails v3德文文件

## v3.0.0-alpha.84 - 2026-05-03

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5307)中新增可停用 macOS 按 Escape 鍵退出全螢幕的選項
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5310)中新增可停用 macOS 按 Escape 鍵退出全螢幕的選項
- 由 @yuseferi 在[PR](https://github.com/wailsapp/wails/pull/5288)中新增 Pausa 社群展示文件

## 變更

- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5308)中更新贊助者 SVG
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5309)中更新圖示產生命令，以處理不支援的平台
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5224)中以三態 ButtonState 取代布林值全螢幕 API，並實作平台繫結

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5315)中防止 WebView2 在控制器狀態為 nil 時執行焦點操作
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5313)中更新 GitHub Actions 工作流程，使其正確參照 PR 的基底分支
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5203)中，於開發模式忽略`*_test.go`檔案，以避免不必要的重新建置
- 由 @wucm667 在[PR](https://github.com/wailsapp/wails/pull/5291)中防止應用程式未執行時 Menu.Update() 發生記憶體區段錯誤

## v3.0.0-alpha.83 - 2026-05-02

## 新增

- 由 @symball 在[PR](https://github.com/wailsapp/wails/pull/5094)中新增 InstallScope 旗標，以及供電腦／使用者安裝使用的建置選項
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5294)中為 BrowserWindow 新增不執行任何操作的 SetScreen 方法，以滿足 Window 介面

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5295)中偵測 NVIDIA GPU，並在 Linux 上停用 DMA-BUF 轉譯器
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5109)中修正 Git PR 範本，使其指向正確的意見回饋 URL
- 修正一系列 Windows 系統匣`SetMenu`當機問題。其原因是損壞的`DestroyMenu`系統呼叫傳入了四個引數而非一個，導致每次呼叫都傳回 FALSE 且未釋放任何資源。此外，在重建選單時釋放 HMENU 與 HBITMAP 控制代碼（包括執行階段透過`MenuItem.SetBitmap`配置的控制代碼）、重設`Win32Menu.Update`中過期的核取方塊／選項按鈕對應，並移除`systemtray.updateMenu`中會使配置量加倍的多餘`Update()`呼叫。長時間執行的系統匣應用程式不再於每次重建選單時洩漏 GDI／USER 物件。

## v3.0.0-alpha.82 - 2026-05-01

## 修正

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5232)中修正桌面檔案產生功能，使其正確處理桌面名稱

## v3.0.0-alpha.81 - 2026-04-30

## 變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5286)中將每夜版本發布排程調整為 UTC 15:00

## 修正

- 修正 Retina Mac 上 Screen Bounds、WorkArea 與 Size 數值減半的問題 - （#5168）

## v3.0.0-alpha.80 - 2026-04-29

## 變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5285)中更新文件相依套件與內容集合載入器

## v3.0.0-alpha.79 - 2026-04-29

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5270)中授予 trigger-release 工作 actions: write 權限

## 變更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5283)中將發布工作預設為 master 分支，並更新變更日誌措辭
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5282)中更新自動變更日誌工作流程，以使用最新版本
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5280)中新增路徑篩選條件並移除無效的工作流程，以提升工作流程效率
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5274)中更新文件，使範例連結參照 master 分支
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5272)中更新v3的文件與範例

## 修正

- 由 @AkagiYui 在[PR](https://github.com/wailsapp/wails/pull/5265)中增強反向 Proxy，加入重試邏輯，並在開發期間強制使用 IPv4
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5281)中重寫未發布變更日誌的觸發工作流程

## 移除

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5267)中移除用於各種測試目的的 Shell 測試指令碼
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5266)中刪除v3-alpha 文件部署工作流程與 CNAME 記錄

### 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5196)中將「前端路由」項目新增至側邊欄導覽
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5185)中新增前端路由指南，並提供各框架專屬的建議
- 新增對模態工作表的支援（macOS）
- 由 @leaanthony 升級 ghw 版本，以更妥善支援 Apple 裝置（#4977）
- 將`GetBadge`方法新增至 Dock 服務
- 為`wails3 build`命令新增`-tags`旗標，以傳遞自訂 Go 建置標籤（例如`wails3 build -tags gtk4`）（#4957）
- 新增繫結產生器自動產生列舉的相關文件，包括專用的「列舉」頁面與側邊欄導覽（#4972）
- 為`wails3 build`命令新增`-tags`旗標，以傳遞自訂 Go 建置標籤（例如`wails3 build -tags gtk4`）（#4957）
- 在`v3/examples/web-apis/`中新增 Web API 範例，示範41項瀏覽器 API，包括儲存（localStorage、sessionStorage、IndexedDB、Cache API）、網路（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、媒體（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、裝置（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、效能（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）等
- 新增 WebView API 相容性檢查器範例（`v3/examples/webview-api-check/`），可跨平台測試超過200項瀏覽器 API
- 新增`internal/libpath`套件，用於在 Linux 上尋找原生程式庫路徑，支援平行搜尋、快取及 Flatpak/Snap/Nix
- <strong>開發中：</strong>新增實驗性的 Linux WebKitGTK 6.0 / GTK4 支援，可透過`-tags gtk4`啟用（GTK3/WebKit2GTK 4.1仍為預設值）
- 注意：在平鋪式視窗管理員（例如 Hyprland、Sway）上，由於視窗管理員會控制視窗幾何配置，最小化／最大化操作可能無法如預期運作
- 由 @AbdelhadiSeddar 在文件的**在 JavaScript 中監聽事件**章節新增如何使用**單次處理常式**的說明
- 由 @leaanthony 為`WebviewWindowOptions`新增`UseApplicationMenu`選項，讓 Windows/Linux 上的視窗可繼承透過`app.Menu.Set()`設定的應用程式選單
- 由 @wimaha 新增使用`.icon`檔案（Apple Icon Composer 格式）產生 Liquid Glass 圖示和資產目錄（macOS）的支援（#4934）
- 新增用於無頭／Web 部署的實驗性伺服器模式（`-tags server`）。此模式可讓 Wails 應用程式在不依賴原生 GUI 的情況下作為 HTTP 伺服器執行。請使用`wails3 task build:server`建置。詳情請參閱`examples/server`。
- 新增`internal/libpath`套件，用於在 Linux 上尋找原生程式庫路徑，支援平行搜尋、快取及 Flatpak/Snap/Nix
- 由 @leaanthony 為`MacWindow`新增`CollectionBehavior`選項，用於控制視窗在 macOS Spaces 和全螢幕模式下的行為（#4756）
- 由 @leaanthony 為 pkg/application 新增單元測試
- 由 @leaanthony 為 MSIX 封裝新增自訂通訊協定支援
- 新增 Linux 桌面環境偵測[PR #4797](https://github.com/wailsapp/wails/pull/4797)
- 由 @leaanthony 在 JavaScript 執行階段新增`Window.Print()`方法，以便從前端觸發列印對話方塊（#4290）
- 由 @leaanthony 在 Linux 的`wails3 doctor`輸出中新增`XDG_SESSION_TYPE`
- 由 @leaanthony 為 Linux 新增其他 WebKit2 載入狀態變更事件：`WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished`（#3896）
- 由 @leaanthony 在 Linux 的`wails3 doctor`輸出中新增`XDG_SESSION_TYPE`
- 在 Linux 建置期間產生`.desktop`檔案，而不僅限於封裝期間（#4575）
- 由 @leaanthony 新增 Linux 執行階段相依性文件，其中包含各發行版專屬的套件名稱和 nfpm 封裝範例（#4339）
- 由 @leaanthony 在 Linux 的`wails3 doctor`輸出中新增 NVIDIA 驅動程式版本資訊
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中為原始訊息處理常式新增來源資訊
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中新增 macOS 通用連結支援
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重構繫結傳輸層
- 由 @chinenual 在[PR](https://github.com/wailsapp/wails/pull/4760)中為 helloworld 範本新增 aria-label 識別碼，讓 Appium 測試用戶端可以輕鬆測試範例應用程式
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中為原始訊息處理常式新增來源資訊
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中新增 macOS 通用連結支援
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重構繫結傳輸層
- 由 @fbbdev 和 @ianvs 在[#4633](https://github.com/wailsapp/wails/pull/4633)中新增型別化事件
- 新增`systray-clock`範例，展示會即時更新工具提示的無頭系統匣（#4653）。
- 由 @Tolfx 在 #4510中新增適用於 Windows 的 NSIS Protocol 範本
- 由 @Tolfx 在 #4510中新增 build-assets 測試
- macOS：由 @nidib 在[#4588](https://github.com/wailsapp/wails/pull/4588)中讓選單列顯示原生視窗控制項
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中新增 macOS Dock 服務，用於在 Dock 中隱藏／顯示應用程式圖示
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中新增 macOS Dock 服務，用於在 Dock 中隱藏／顯示應用程式圖示
- 由 @leaanthony 在[#4534](https://github.com/wailsapp/wails/pull/4534)中新增 macOS 原生 Liquid Glass 效果支援，使用 NSGlassEffectView（macOS 15.0 以上版本），並以 NSVisualEffectView 作為後備方案，且提供完整的材質自訂選項
- 由 @leaanthony 在[#4500](https://github.dev/wailsapp/wails/pull/4500)中新增瀏覽器 URL 清理功能。以 @APShenkin 的[#4484](https://github.com/wailsapp/wails/pull/4484)為基礎。
- 由[@leaanthony](https://github.com/leaanthony)新增 Windows/Mac 內容保護功能，以[@Taiterbase](https://github.com/Taiterbase)在此[PR](https://github.com/wailsapp/wails/pull/4241)中的原始工作為基礎
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4488)中新增支援，可透過`wails3 build`和`wails3 package`別名將 CLI 變數傳遞給 Task 命令（#4422）
- 由[@atterpac](https://github.com/atterpac)在[#4318](https://github.com/wailsapp/wails/pull/4318)中新增拖放區支援，事件來源會提供被拖放元素的資料
- 在[PR](https://github.com/wailsapp/wails/pull/4467)中為`WindowsWindow`選項新增`AdditionalLaunchArgs`，以便將額外的命令列引數傳遞給 WebView2 瀏覽器。
- 由[@triadmoko](https://github.com/triadmoko)在[PR](https://github.com/wailsapp/wails/pull/4286)中新增在 wails init 後自動執行 go mod tidy 的功能
- 由 @leaanthony 在[PR](https://github.dev/wailsapp/wails/pull/4463)中新增 Windows Snap Assist 功能
- 在[PR](https://github.com/wailsapp/wails/pull/4467)中為`WindowsWindow`選項新增`AdditionalLaunchArgs`，以便將額外的命令列引數傳遞給 WebView2 瀏覽器。
- 由[@triadmoko](https://github.com/triadmoko)在[PR](https://github.com/wailsapp/wails/pull/4286)中新增在 wails init 後自動執行 go mod tidy 的功能
- 由 @leaanthony 在[PR](https://github.dev/wailsapp/wails/pull/4463)中新增 Windows Snap Assist 功能
- 由[@almas-x](https://github.com/almas-x)在[PR](https://github.com/wailsapp/wails/pull/4427)中新增 Windows `getAccentColor`實作
- 由[@almas-x](https://github.com/almas-x)在[PR](https://github.com/wailsapp/wails/pull/4427)中新增 Windows `getAccentColor`實作
- Windows 深色主題選單和選單列。由 @leaanthony 在[a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)中實作
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4405)中重新命名內建服務，使 JS/TS 繫結更清楚
- 由[@etesam913](https://github.com/etesam913)新增`app.Env.GetAccentColor`，用於取得使用者系統的強調色。支援 macOS。
- 由[@atterpac](https://github.com/atterpac)在[#4137](https://github.com/wailsapp/wails/pull/4137)中新增`window.ToggleFrameless()` API
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4345)中新增 Linux 各發行版專屬的建置相依套件
- @atterpac 在[PR](https://github.com/wailsapp/wails/pull/4404)中新增繫結指南
- **整理測試基礎架構**：將 Docker 測試檔案移至專用的`test/docker/`目錄，並採用最佳化映像檔及提升建置可靠性；由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中完成
- **改進資源管理模式**：在範例中加入適當的事件處理常式清理，以及能感知 context 的 goroutine 管理；由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中完成
- [@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中新增對 aarch64 AppImage 建置的支援
- [@leaanthony](https://github.com/leaanthony)為`wails doctor`新增診斷章節
- [@leaanthony](https://github.com/leaanthony)在呼叫服務方法時將視窗加入 context
- [@leaanthony](https://github.com/leaanthony)新增`window-call`範例，示範如何得知是哪個視窗正在呼叫服務
- [@leaanthony](https://github.com/leaanthony)新增選單指南
- [@leaanthony](https://github.com/leaanthony)改進 panic 處理
- [@leaanthony](https://github.com/leaanthony)新增選單指南
- [@fbbdev](https://github.com/fbbdev)在[#4024](https://github.com/wailsapp/wails/pull/4024)中為 Service API 新增文件註解
- [@leaanthony](https://github.com/leaanthony)在[#4024](https://github.com/wailsapp/wails/pull/4024)中新增`application.NewServiceWithOptions`函式，以使用額外設定初始化服務
- [@FalcoG](https://github.com/FalcoG)與[@leaanthony](https://github.com/leaanthony)在[#4031](https://github.com/wailsapp/wails/pull/4031)中改進選單控制
- [@leaanthony](https://github.com/leaanthony)新增更多文件
- [@leaanthony](https://github.com/leaanthony)新增對在標準事件監聽器中取消事件的支援
- [@leaanthony](https://github.com/leaanthony)新增系統匣`Hide`、`Show`及`Destroy`支援
- [@leaanthony](https://github.com/leaanthony)新增系統匣`SetTooltip`支援。原始構想來自[@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- [@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中讓繫結產生器在不支援類型的警告中回報套件路徑
- [@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中新增繫結產生器對泛型別名的支援
- [@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中新增繫結產生器對`omitzero` JSON 旗標的支援
- [@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中新增`//wails:ignore`指示詞，以防止為所選服務方法產生繫結
- [@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中為服務與模型新增`//wails:internal`指示詞，允許類型在 Go 中匯出，但不在 JS/TS 中匯出
- [@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中新增繫結產生器對別名類型常數的支援，以允許弱型別列舉
- [@fbbdev](https://github.com/fbbdev)在[#4068](https://github.com/wailsapp/wails/pull/4068)中新增針對 Go 1.24功能的繫結產生器測試
- 在`OSInfo.Branding`中新增對 macOS 15「Sequoia」的支援，以改善作業系統版本偵測；見[#4065](https://github.com/wailsapp/wails/pull/4065)
- [@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中新增`PostShutdown`掛鉤，以便在關閉流程完成後執行自訂程式碼
- [@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中新增`FatalError`結構，以支援在自訂錯誤處理常式中偵測嚴重錯誤
- [@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中將服務的啟動與關閉順序標準化並加以記錄
- [@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中新增應用程式啟動／關閉順序的測試框架，以及服務啟動／關閉測試
- [@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中新增`RegisterService`方法，以便在應用程式建立後註冊服務
- [@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中，於應用程式與服務選項新增`MarshalError`欄位，以自訂繫結呼叫的錯誤處理
- [@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中新增可取消的 Promise 包裝器，可沿 Promise 鏈傳遞取消要求
- [@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中新增將繫結呼叫的取消作業繫結至`AbortSignal`的能力
- [@leaanthony](https://github.com/leaanthony)讓 WML 除了支援一般的`wml-*`屬性外，也支援`data-wml-*`屬性
- [@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為所有服務新增`Configure`方法，以支援延後設定與動態重新設定
- `fileserver`服務在未設定時會傳送503「服務無法使用」回應；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- `kvstore`服務在未設定時預設提供記憶體內鍵值儲存區；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- [@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`kvstore`服務新增`Load`方法，以便在設定變更後從檔案重新載入資料
- [@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`kvstore`服務新增`Clear`方法，以刪除所有鍵
- [@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中，於`log`服務新增`Level`類型，以提供 JS 端的記錄層級常數
- [@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`log`服務新增`Log`方法，以動態指定記錄層級
- `sqlite`服務在未設定時預設提供記憶體內資料庫；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- [@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`sqlite`服務新增`Close`方法，以手動關閉資料庫
- [@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`sqlite`服務的查詢方法新增取消支援
- [@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中為`sqlite`服務新增具備 JS 繫結的預備陳述式支援
- 由[Lea Anthony](https://github.com/leaanthony)在[PR](https://github.com/wailsapp/wails/pull/3537)中加入 Gin 支援，此項工作以[@AnalogJ](https://github.com/AnalogJ)在此[PR](https://github.com/wailsapp/wails/pull/3537)中的原始工作為基礎
- 由[@oSethoum](https://github.com/osethoum)在[#4134](https://github.com/wailsapp/wails/pull/4134)中修正自動儲存與密碼自動儲存一律啟用的問題
- 由[@leaanthony](https://github.com/leaanthony)為視窗新增`SetMenu()`，以便為視窗設定選單
- 由[@popaprozac](https://github.com/popaprozac)在[#4098](https://github.com/wailsapp/wails/pull/4098)中加入通知支援
-  由[@wimaha](https://github.com/wimaha)在[#4177](https://github.com/wailsapp/wails/pull/4177)中加入 macOS 檔案關聯支援
- 由[@leaanthony](https://github.com/leaanthony)加入用於遞增語意版本的`wails3 tool version`
- 由[@popaprozac](https://github.com/popaprozac)在[#](https://github.com/wailsapp/wails/pull/4234)中加入 macOS 與 Windows 的徽章支援
- 由[@fbbdev](https://github.com/fbbdev)與[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中加入已註冊／強型別事件的支援
- 由[@fbbdev](https://github.com/fbbdev)與[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中加入為自訂事件註冊掛鉤的功能
- 由[@Krzysztofz01](https://github.com/Krzysztofz01)與[@rcalixte](https://github.com/rcalixte)加入`app.OpenFileManager(path string, selectFile bool)`，以在路徑`path`開啟系統檔案管理員，並可選擇透過`selectFile`醒目提示
- 由[@leaanthony](https://github.com/leaanthony)為`wails3 init`命令新增`-git`旗標
- 由[@leaanthony](https://github.com/leaanthony)新增`wails3 generate webview2bootstrapper`命令
- 由[@leaanthony](https://github.com/leaanthony)在 runtime 中新增`init()`方法，以便手動初始化 runtime
- 由[@leaanthony](https://github.com/leaanthony)在視窗的 WindowOptions 中新增`WindowDidMoveDebounceMS`選項
- 由[@leaanthony](https://github.com/leaanthony)新增單一執行個體功能。此功能以 @APshenkin 的[v2 PR](https://github.com/wailsapp/wails/pull/2951)為基礎。
- 由[@leaanthony](https://github.com/leaanthony)新增`wails3 generate template`命令
- 由[@leaanthony](https://github.com/leaanthony)新增`wails3 releasenotes`命令
- 由[@leaanthony](https://github.com/leaanthony)新增`wails3 update cli`命令
- 由[@leaanthony](https://github.com/leaanthony)為`wails3 generate bindings`命令新增`-clean`選項
- 由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中允許建置 aarch64 (arm64) AppImage Linux 套件
- 由 @ansxuman 在[#3958](https://github.com/wailsapp/wails/pull/3958)中為贊助者新增超連結
- 由以下人員加入 Linux 的 deb、rpm 與 Arch Linux 套件建置支援
- 由以下人員加入 Darwin 通用建置與套件支援
- 由以下人員將事件文件加入網站
- 為非 SSR 開發設定的 sveltekit 與 sveltekit-ts 範本
- 由以下人員使用新的`wails3 update build-assets`命令更新建置資產
- 由以下人員加入測試 HTML Drag and Drop API 的範例
- 由[leaanthony](https://github.com/leaanthony)在以下項目中加入檔案關聯支援
- 由以下人員新增`wails3 generate runtime`命令
- 新增`InitialPosition`選項，用於指定視窗是否應置中，或
- 由以下人員在`application`套件中新增`Path`與`Paths`方法
- 新增`GeneralAutofillEnabled`與`PasswordAutosaveEnabled` Windows 選項
- 由以下人員新增取得呼叫服務方法之視窗的功能
- 由以下人員新增 WebView2 的`EnabledFeatures`與`DisabledFeatures`選項
- ⊞ 由以下人員新增 DIP 系統，以增強高 DPI 顯示器支援
- ⊞ 由[windom](https://github.com/windom/)在以下項目中加入視窗類別名稱選項
- 服務已擴充，可提供外掛功能。由
- 🐧 在以下項目中加入 WindowDidMove / WindowDidResize 事件
- ⊞ 在以下項目中加入 WindowDidResize 事件
-  新增 Event ApplicationShouldHandleReopen，以便處理 Dock
-  由 @tmclane 在以下項目中將 getPrimaryScreen/getScreens 加入實作
-  由以下人員新增在 macOS 全螢幕模式下顯示工具列的選項
- 🐧 新增 onKeyPress 邏輯，將 Linux 按鍵轉換為快速鍵
- 🐧 由以下人員新增工作`run:linux`
- 由[@almas-x](https://github.com/almas-x)在以下項目中匯出`SetIcon`方法
- 由[@almas-x](https://github.com/almas-x)在以下項目中改善`OnShutdown`
- 由以下人員在`Window`介面中還原`ToggleMaximise`方法
- 由 @leaanthony 在以下項目中為`Environment()`新增更多資訊
- 由以下人員在`Window`介面上公開`WebviewWindow.IsFocused`方法
- 由以下人員在 WML 系統中支援以空格分隔的多個觸發事件
- 由以下人員從隨附的 JS runtime 指令碼匯出 ESM
- 新增繫結產生器旗標，以使用隨附的 JS runtime 指令碼，而非
- 由[@abichinger](https://github.com/abichinger)在 Linux 上實作`setIcon`
- 在 dev 命令中新增`-port`旗標，並支援環境變數
- 由以下人員新增繫結方法呼叫的測試
- ⊞ 由以下人員為已建立的視窗新增`SetIgnoreMouseEvents`
-  由以下人員新增設定視窗堆疊層級（順序）的功能

### 已修正

- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5168)中修正 Retina Mac 上`Screen.Bounds`、`WorkArea`與`Size`減半的問題：在`Physical*`欄位中將 NSScreen 點值轉換為裝置像素，並填入頂層`Screen.X`/`Y`，使多顯示器接觸偵測與工作區域定位正確
- 修正 ScreenManager 中的資料競爭；該問題會在顯示器設定變更時（例如睡眠／喚醒期間熱插拔外接顯示器）造成 WebKit DisplayLink 死結
- 當 Assets.car 存在時，直接將 CFBundleIconName 設為 appicon；由 @symball 在[PR](https://github.com/wailsapp/wails/pull/5154)中實作
- 修正`wails3 doctor`在 Fedora、openSUSE、Arch 與 NixOS 上回報錯誤 WebKitGTK 套件的問題——由於v3在編譯時需要4.1 API，因此已移除4.0後援項目（#5071）
- 修正 openSUSE webkit2gtk doctor 套件名稱（`webkit2gtk4_1-devel` → `webkit2gtk3-devel`，後者才是正確的 openSUSE 套件名稱）（#5071）
- 修正在桌面開發模式下缺少`/wails/custom.js`時發生的`Unexpected token '<'`錯誤。在`loadOptionalScript`中新增針對`/wails/custom.js`的明確404處理常式，以及不區分大小寫的`Content-Type`驗證，以防止將 HTML SPA 備援內容注入為 JavaScript。（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- 修正 macOS 上系統匣選單的醒目提示狀態——選單開啟時，圖示現在會顯示為已選取狀態（#4910）
- 修正 macOS 上附加至系統匣的視窗出現在其他視窗後方的問題——現在會使用正確的快顯視窗層級（#4910）
- 修正整份文件中錯誤的`@wailsio/runtime`匯入範例（#4989）
- 修正 darwin 上無框線視窗無法最小化的問題（#4294）
- 將`node_modules/`排除在 go-task 的最新狀態檢查之外，修正執行`wails3 build`和`wails3 dev`時會停滯20-30分鐘的問題。先前，`sources: "**/*"`萬用字元會導致 go-task 列舉`node_modules/`中的每個檔案並計算其總和檢查碼（50000-100000+ 個檔案，包含 MUI 等大型相依套件），在 Windows/NTFS 上尤其緩慢（#4939）
- 修正 C 的`Screen` typedef 與 X11 Xlib.h 衝突所造成的 GTK4 建置失敗（#4957）
- 修正 macOS 上 Dock 徽章方法的一致性
- 修正`InvisibleTitleBarHeight`套用至所有 macOS 視窗，而非僅套用至無框線或標題列透明視窗的問題（#4960）
- 啟用`InvisibleTitleBarHeight`時，略過視窗邊緣附近的拖曳起始操作，修正從頂端角落調整視窗大小時晃動／抖動的問題（#4960）
- 修正 JS/TS 繫結中以列舉作為鍵的對應型別產生問題（#4437），由 @fbbdev 貢獻
- 修正 Windows 上檔案拖放無法在非100% 顯示縮放比例下運作的問題
- 修正 Windows 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 修正 Windows 上檔案放置座標使用錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正 Linux 上檔案拖放搭配懸停效果時無法可靠運作的問題
- 修正 Linux 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 透過使用`gtk_window_present()`，修正 Linux/GTK4 上顯示／隱藏視窗有時會還原為最小化狀態的問題（#4957）
- 透過`XTranslateCoordinates`／`XMoveWindow`新增以 X11 為條件的支援，修正 Linux/GTK4 上取得／設定視窗位置總是傳回0,0的問題（#4957）
- 新增以訊號為基礎的大小限制機制，以取代已移除的`gtk_window_set_geometry_hints`，修正 Linux/GTK4 上未強制執行視窗大小上限的問題（#4957）
- 透過`gdk_monitor_get_scale`實作正確的 PhysicalBounds 計算與小數縮放支援，修正 Linux/GTK4 上的 DPI 縮放（GTK 4.14+）
- 修正 Linux/GTK4 上建立新視窗時選單項目重複的問題
- 修正 JS/TS 繫結中以列舉作為鍵的對應型別產生問題（#4437），由 @fbbdev 貢獻
- 修正 Windows 上檔案拖放無法在非100% 顯示縮放比例下運作的問題
- 修正 Windows 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 修正 Windows 上檔案放置座標使用錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正 Linux 上檔案拖放搭配懸停效果時無法可靠運作的問題
- 修正 Linux 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 透過`gdk_monitor_get_scale`實作正確的 PhysicalBounds 計算與小數縮放支援，修正 Linux/GTK4 上的 DPI 縮放（GTK 4.14+）
- 修正 Linux/GTK4 上建立新視窗時選單項目重複的問題
- 修正 JS/TS 繫結中以列舉作為鍵的對應型別產生問題（#4437），由 @fbbdev 貢獻
- 修正 App.Window.Current() 未從主執行緒存取 AppKit API 而在 macOS 上造成「幽靈視窗」的問題（#4947），由 @wimaha 貢獻
- 透過實作 WKUIDelegate runOpenPanelWithParameters，修正 HTML `<input type="file">`無法在 macOS 上運作的問題（#4862）
- 修正在 macOS/Linux 上使用`@wailsio/runtime` npm 模組時，原生檔案拖放無法運作的問題（#4953），由 @leaanthony 貢獻
- 修正跨套件型別別名的繫結產生問題（#4578），由 @fbbdev 貢獻
- 修正 GTK 執行緒安全違規造成 OpenFileDialog 在 Linux 上當機的問題（#3683），由 @ddmoney420 貢獻
- 修正在隱藏或已銷毀的視窗上呼叫`Focus()`時發生 SIGSEGV 當機的問題（#4890），由 @ddmoney420 貢獻
- 修正在 Linux 上設定空白圖示或點陣圖時可能發生 panic 的問題（#4923），由 @ddmoney420 貢獻
- 修正從服務繫結呼叫 ErrorDialog 時在 macOS 上當機的問題（#3631），由 @leaanthony 貢獻
- 讓選單在 Windows 作業系統的`v3\examples\dialogs`中顯示，由 @ndianabasi 貢獻
- 修正頁面重新載入期間引發 TypeError 的競爭條件（#4872），由 @ddmoney420 貢獻
- 移除`Collector.IsVoidAlias()`方法中的全域狀態，修正繫結產生器測試輸出不正確的問題（#4941），由 @fbbdev 貢獻
- 修正`<input type="file">`檔案選擇器無法在 macOS 上運作的問題（#4862），由 @leaanthony 貢獻
- 修正`Position()`與`SetPosition()`在 macOS 上使用不一致的座標系統，導致儲存／還原狀態時視窗位置偏移的問題（#4816），由 @leaanthony 貢獻
- 修正已透過應用程式資訊清單設定 DPI 感知時發生的 SetProcessDpiAwarenessContext「存取遭拒」錯誤（#4803）
- 更新鍵盤快速鍵文件頁面，並修正`KeyBinding.Add`回呼參數的型別，由 @ndianabasi 貢獻
- 修正自訂繫結產生方式的文件，必須使用`-d String`，而非`-o String`
- 修正在`menu.Update()`時選單未清除子項目的問題
- 修正文件中過時的 Manager API 參照（已更新31檔案，改用`app.Window.New()`、`app.Event.Emit()`等新模式），由 @leaanthony 貢獻
- 修正 WebKit 覆寫訊號處理常式，導致繫結至 JS 的 Go 方法發生 panic 時 Linux 當機的問題（#3965），由 @leaanthony 貢獻
- 修正 SaveFileDialog.SetFilename() 在 Linux 上無效的問題（#4841），由 @samstanier 貢獻
- 修正拖放範例中的放置座標顯示為 undefined 的問題
- 修正 APP_NAME 包含空格時無法建立 macOS 應用程式套件的問題（大括號展開問題）
- 修正在 Windows 上呼叫服務方法時發生索引超出範圍 panic 的問題（還原 goccy/go-json）
- 修正 Windows 上檔案拖放無法在非100% 顯示縮放比例下運作的問題
- 修正 Windows 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 修正 Windows 上檔案放置座標使用錯誤像素空間的問題（實體像素與 CSS 像素）
- 修正 Linux 上檔案拖放搭配懸停效果時無法可靠運作的問題
- 修正 Linux 上啟用檔案放置後，HTML5 內部拖放功能失效的問題
- 更新所有作業系統的 Taskfile.yml 檔案中的全部命令，以支援`APP_NAME`等變數包含空格的情況，由 @ndianabasi 貢獻
- 修正在 Linux 上執行 'build:universal:lipo:go' 工作時發生的命令引數錯誤，由 @wux1an 貢獻
- 修正在 Linux 上執行 'wails3 build GOOS=darwin GOARCH=arm64' 時出現的 Docker 錯誤「undefined symbol: **<em>ubsan</em>handle_xxxxxxx」，由 @wux1an 完成
- 整合自訂通訊協定文件並新增 Universal Links 章節，由 @leaanthony 完成
- 新增防護以避免並行呼叫 TrackPopupMenuEx，修正在 Windows 上重複點擊系統匣圖示時選單當機的問題（#4151），由 @leaanthony 完成
- 防止在 app.Run() 之前呼叫 systray.Run() 時應用程式當機，由 @leaanthony 完成
- 修正在啟用 ApplicationShouldTerminateAfterLastWindowClosed 時，透過 Hide()/Show() 切換視窗可見性會導致 macOS 應用程式當機的問題（#4389），由 @leaanthony 完成
- 修正在 macOS 和 Windows 上重複開啟內容選單時發生的記憶體洩漏（#4012），由 @leaanthony 完成
- 修正 macOS 上未重複使用內容選單原生資源，導致每次顯示時都建立新選單的問題（#4012），由 @leaanthony 完成
- 修正在使用`Hidden: true`啟動應用程式後，點擊 macOS Dock 圖示不會顯示已隱藏視窗的問題（#4583），由 @leaanthony 完成
- 修正 CGO 呼叫使用了錯誤的視窗指標類型，導致 macOS 列印對話方塊無法開啟的問題（#4290），由 @leaanthony 完成
- 修正 appmenu-gtk-module 存取尚未實體化的視窗，導致 Wayland 上的視窗選單當機（#4769），由 @leaanthony 完成
- 修正應用程式名稱含有無效字元（空格、括號等）時 GTK 應用程式當機的問題，由 @leaanthony 完成
- 修正在 Windows 上初始化拖放功能時出現「not enough memory」錯誤的問題（#4701），由 @overlordtm 完成
- 修正 URI 跳脫處理不正確，導致 Linux 上的檔案總管開啟錯誤目錄的問題（#4397），由 @leaanthony 完成
- 透過自動偵測`.relr.dyn` ELF 區段並停用剝除，修正在現代 Linux 發行版（Arch、Fedora 39+、Ubuntu 24.04+）上建置 AppImage 失敗的問題（#4642），由 @leaanthony 完成
- 修正`wails doctor`在 Fedora／DNF 系統上誤報已安裝 webkit 套件的問題（#4457），由 @leaanthony 完成
- 修正預設`config.yml`會以正式環境組建執行`wails3 dev`的問題，由 @mbaklor 完成
- 修正 iOS 服務虛設常式匯入不存在的套件而導致建置失敗的問題，由 @leaanthony 完成
- 修正 debug/info 方法中的結構化記錄導致「no formatting directives」錯誤的問題，由 @leaanthony 完成
- 移除從行動平台合併中意外納入的暫時偵錯列印陳述式，由 @leaanthony 完成
- 透過自動停用 DMA-BUF 算繪器，修正在搭載 NVIDIA GPU 的 Wayland 上 WebKitGTK 當機（Error 71 Protocol error）的問題，由 @leaanthony 完成
- 解決 Linux 上`application.WebviewWindowOptions.BackgroundColour`中的 Alpha 值遭忽略的問題（[#4722](https://github.com/wailsapp/wails/pull/4722)，@BradHacker）
- 修正在未提供自訂圖示時，Windows 系統匣圖示未預設使用應用程式圖示的問題（#4704）
- 追蹤`HICON`的擁有權，僅銷毀使用者建立的控制代碼，以防止 Explorer 回收時當機（#4653）。
- 在銷毀時釋放 Windows 系統佈景主題接聽器和保留的系統匣圖示，以停止洩漏 goroutine 和裝置內容（#4653）。
- 將系統匣工具提示截斷至127個 UTF-16單位，以免損毀代理對和多位元組字符（#4653）。
- 修正 Windows 套件工作失敗的問題（#4667）
- 修正 Linux taskfile 中的 Linux AppImage appicon 變數[PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 修正 go-webview2 v1.0.22簽章變更所導致的 Windows 建置錯誤（#4513、#4645）
- 修正 Linux taskfile 中的 Linux AppImage appicon 變數[PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 透過移除`<.Info.Protocol>`至`<.Protocol>`，修正 Linux desktop.tmpl 的通訊協定範圍，由 @Tolfx 在 #4510 中完成
- 修正 liquid glass 示範中的重複定義錯誤[#4542](https://github.com/wailsapp/wails/pull/4542)，由 @Etesam913 完成
- 修正 Linux 上的系統匣選單更新[#4604](https://github.com/wailsapp/wails/issues/4604)，由[@JackDoan](https://github.com/JackDoan)完成
- 修正在 Windows 上建立隱藏視窗時出現白色視窗的問題，由 @leaanthony 在[#4612](https://github.com/wailsapp/wails/pull/4612)中完成
- 修正文件中的 notifications 套件匯入路徑，由 @rxliuli 在[#4617](https://github.com/wailsapp/wails/pull/4617)中完成
- 修正使用 npm 套件 @wailsio/runtime 時拖放功能無法運作的問題（#4489），由 @leaanthony 在 #4616 中完成
- Windows：修正啟動時視窗閃爍，以及隱藏視窗錯誤顯示的問題；由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4600)中完成。
- 修正 Wayland 視窗最大化時的大小問題（https://github.com/wailsapp/wails/issues/4429），由[@samstanier](https://github.com/samstanier)完成
- 修正 Wayland 視窗最大化時的大小問題（https://github.com/wailsapp/wails/issues/4429），由[@samstanier](https://github.com/samstanier)完成
- 修正 liquid glass 示範中的重複定義錯誤[#4542](https://github.com/wailsapp/wails/pull/4542)，由 @Etesam913 完成
- 修正 AssetServer 在 macOS 上可能當機的問題[#4576](https://github.com/wailsapp/wails/pull/4576)，由 @jghiloni 完成
- 修正使用 NextJs 建置時的編譯問題。由 @rev42 在[#4585](https://github.com/wailsapp/wails/pull/4585)中修正
- 修正夜間版本的管線[#4597](https://github.com/wailsapp/wails/pull/4597)，由 @riadafridishibly 完成
- 修正 liquid glass 示範中的重複定義錯誤[#4542](https://github.com/wailsapp/wails/pull/4542)，由 @Etesam913 完成
- 修正 AssetServer 在 macOS 上可能當機的問題[#4576](https://github.com/wailsapp/wails/pull/4576)，由 @jghiloni 完成
- 修正使用 NextJs 建置時的編譯問題。由 @rev42 在[#4585](https://github.com/wailsapp/wails/pull/4585)中修正
- 修正夜間版本的管線[#4597](https://github.com/wailsapp/wails/pull/4597)，由 @riadafridishibly 完成
- 修正 liquid glass 示範中的重複定義錯誤[#4542](https://github.com/wailsapp/wails/pull/4542)，由 @Etesam913 完成
- 修正 Windows 上的 SetBackgroundColour，由 @PPTGamer 在[PR](https://github.com/wailsapp/wails/pull/4492)中完成
- 更新文件以反映 Manager API 重構所帶來的變更，由 @yulesxoxo 在[PR #4476](https://github.com/wailsapp/wails/pull/4476)中完成
- 修正 Linux taskfile 中 Linux .desktop 檔案的 appicon 變數[PR #4477](https://github.com/wailsapp/wails/pull/4477)
- 更新文件以反映 Manager API 重構所帶來的變更，由 @yulesxoxo 在[PR #4476](https://github.com/wailsapp/wails/pull/4476)中完成
- 修正[#4456](https://github.com/wailsapp/wails/issues/4456)中回報的 Windows nil 指標解除參照錯誤，由 @leaanthony 在[#4460](https://github.com/wailsapp/wails/pull/4460)中完成
- 在 macOS WKWebView 中新增對`allowsBackForwardNavigationGestures`的支援，以啟用雙指滑動導覽手勢（#1857）
- 修正初始設為停用的選單項目無法使用 onClick 的問題，由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中完成。感謝 @IanVS 的初步調查。
- 修正建置失敗時未清理 Vite 伺服器的問題（#4403）
- 修正在 Windows 上關閉或取消`SaveFileDialog`時發生 panic 的問題。由 @hkhere 在[PR](https://github.com/wailsapp/wails/pull/4284)中修正
- 由[@mbaklor](https://github.com/mbaklor)在[#4259](https://github.com/wailsapp/wails/pull/4259)中修正 Windows 上 HTML 層級的拖放功能
- 在 macOS WKWebView 中新增對`allowsBackForwardNavigationGestures`的支援，以啟用雙指滑動導覽手勢（#1857）
- 由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中修正初始設為停用的選單項目無法使用 onClick 的問題。感謝 @IanVS 進行初步調查。
- 修正建置失敗時未清理 Vite 伺服器的問題（#4403）
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4450)中修正 Windows 上的通知剖析
- 由[@kodumulo](https://github.com/kodumulo)在[#4390](https://github.com/wailsapp/wails/issues/4390)中修正 doctor 命令，使其檢查 Windows SDK 相依性
- 由[@etesam913](https://github.com/etesam913)在[#4366](https://github.com/wailsapp/wails/pull/4366)中修正 Mac 上 processURLRequest 的 nil 指標解除參照問題
- 由[@bh90210](https://github.com/bh90210)在[#4287](https://github.com/wailsapp/wails/pull/4287)中修正造成 Linux 無法使用篩選對話方塊的錯誤
- 由[@leaanthony](https://github.com/leaanthony)在[#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)中修正 Windows 與 Linux 的「編輯」選單問題
- 由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中將 macOS .plist 檔案的最低系統版本從10.13.0更新為10.15.0
- 由[@leaanthony](https://github.com/leaanthony)修正視窗 ID 跳號問題
- 由[@leaanthony](https://github.com/leaanthony)修正呼叫 RegisterContextMenu 時的 nil 選單問題
- 由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中修正繫結產生器輸出中的相依性循環
- 由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中修正繫結產生器輸出中先使用後定義的錯誤
- 由[@fbbdev](https://github.com/fbbdev)在[#4023](https://github.com/wailsapp/wails/pull/4023)中將建置旗標傳遞給繫結產生器
- 由[@leaanthony](https://github.com/leaanthony)將 Windows Taskfile 中的路徑改為使用正斜線，以確保其可在非 Windows 平台上運作
- 由[@leaanthony](https://github.com/leaanthony)修正 Mac + Mac JS 事件
- 由[@leaanthony](https://github.com/leaanthony)修正 macOS 上的事件死結
- 由[@leaanthony](https://github.com/leaanthony)修正 Windows 上初始化 Window 時，已提供 HTML 但未提供 JS 所造成的`Parameter incorrect`錯誤
- 由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中修正資產伺服器進行內容類型探測時所用的回應前綴大小
- 由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中修正資產伺服器對根索引路徑上非404回應的處理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器測試泛型型別屬性時的未定義行為
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正底層型別與具名包裝型別屬性不同時，繫結產生器對模型的輸出
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器對映射鍵型別的輸出及預處理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器對實作 marshaler 介面的結構體所產生的輸出
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器對涉及泛型型別之型別循環的偵測
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器輸出中對未匯出模型的無效參照
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中將注入的程式碼移至服務檔案末尾
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修正繫結產生器對檔案關閉作業所傳回錯誤的處理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中抑制對僅定義生命週期或 http 方法、未定義其他繫結方法之服務的警告
- 由[@marcus-crane](https://github.com/marcus-crane)在[#4056](https://github.com/wailsapp/wails/pull/4056)中修正非 React 範本在使用淺色系統配色時無法顯示 Hello World 頁尾的問題
- 由[@leaanthony](https://github.com/leaanthony)修正 macOS 上隱藏的選單項目
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中修正訊息處理器對錯誤的處理與格式化
-  由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中修正結束應用程式時略過服務關閉程序的問題
-  由[@leaanthony](https://github.com/leaanthony)確保選單更新在主執行緒上進行
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中改善拖曳與調整大小機制，使其更穩健，並更貼近平台的預期行為
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中修正[#4097](https://github.com/wailsapp/wails/issues/4097)：Webpack/angular 會捨棄執行階段初始化程式碼
- 由[@IanVS](https://github.com/IanVS)在[#4116](https://github.com/wailsapp/wails/pull/4116)中修正初始為隱藏狀態的選單項目
- 修正當無副檔名的請求中`[request]`不存在但`[request].html`存在時，assetFileServer 不會提供`.html`檔案的問題
- 由[@robin-samuel](https://github.com/robin-samuel)在[#4125](https://github.com/wailsapp/wails/pull/4125)中修正圖示產生路徑
- 由[@oSethoum](https://github.com/osethoum)在[#4130](https://github.com/wailsapp/wails/pull/4130)中修正未發出`fullscreen`、`unfullscreen`、`unminimise`及`unmaximise`事件的問題
- 由[@robin-samuel](https://github.com/robin-samuel)在[#4126](https://github.com/wailsapp/wails/pull/4126)中修正設定內預設版本的前綴不正確所導致的 NSIS 錯誤
- 由[TheGB0077](https://github.com/TheGB0077)在[#4188](https://github.com/wailsapp/wails/pull/4188)中修正 Dialogs 執行階段函式在 Windows 上傳回逸出路徑的問題
- 由[@leaanthony](https://github.com/leaanthony)修正 HKCU 中的 Webview2 偵測路徑。
- 由[@leaanthony](https://github.com/leaanthony)修正 macOS 上的輸入問題。
- 由[@yulesxoxo](https://github.com/yulesxoxo)在[#4219](https://github.com/wailsapp/wails/pull/4219)中修正 Windows 圖示產生工作檔案的名稱。
- 由[@leaanthony](https://github.com/leaanthony)根據 @kron 的工作成果，修正無框線視窗的透明度問題。
- 由[@leaanthony](https://github.com/leaanthony)根據 @kron 的工作成果，修正視窗停用或最小化時的焦點呼叫。
- 由[@leaanthony](https://github.com/leaanthony)根據 @kron 的工作成果，修正工作列重新啟動後系統匣未顯示的問題。
- 已修正fallbackResponseWriter未實作[#4245](https://github.com/wailsapp/wails/pull/4245)中的Flush()
- 已由[@superDingda]修正fallbackResponseWriter未實作Flush()的問題，位於[#4236](https://github.com/wailsapp/wails/issues/4236)
- 已由[@joshhardy](https://github.com/joshhardy)修正macOS視窗在非同步Go繫結函式呼叫仍待處理時關閉所造成的當機，位於[#4354](https://github.com/wailsapp/wails/pull/4354)
- 已由[@leaanthony](https://github.com/leaanthony)修正Windows效率模式啟動時的競爭條件
- 已由[@leaanthony](https://github.com/leaanthony)修正Windows圖示控制代碼的清理問題。
- 已由[@PPTGamer](https://github.com/PPTGamer)修正Windows上的`OpenFileManager`，位於[#4375](https://github.com/wailsapp/wails/pull/4375)。
- 已由@atterpac修正Linux的最小／最大寬度選項，位於[#3979](https://github.com/wailsapp/wails/pull/3979)
- 已由@atterpac透過提升npm版本更新TypeScript範本的型別定義，位於[#3966](https://github.com/wailsapp/wails/pull/3966)
- 已由@atterpac修正SvelteKit範本的CSS參照，位於[#3945](https://github.com/wailsapp/wails/pull/3945)
- 已由[@leaanthony](https://github.com/leaanthony)確保window run()中的關鍵回呼會在主執行緒上呼叫
- 已由[@leaanthony](https://github.com/leaanthony)修正對話方塊目錄選擇器範例
- 已由[@leaanthony](https://github.com/leaanthony)建立index.html遺失時顯示的新中文錯誤頁面
-  已由[@leaanthony](https://github.com/leaanthony)確保`windowDidBecomeKey`回呼在主執行緒上執行
-  已由[@leaanthony](https://github.com/leaanthony)支援無邊框視窗全螢幕顯示
-  已由[@leaanthony](https://github.com/leaanthony)改善視窗銷毀邏輯
-  已由[@leaanthony](https://github.com/leaanthony)修正視窗附加至系統匣時的位置邏輯
-  已由[@leaanthony](https://github.com/leaanthony)支援無邊框視窗全螢幕顯示
- 已由[@leaanthony](https://github.com/leaanthony)修正事件處理
- 已由[@leaanthony](https://github.com/leaanthony)修正視窗關閉邏輯
- 通用Taskfile現在預設會為TypeScript範本產生TypeScript繫結，由[@leaanthony](https://github.com/leaanthony)完成
- 已由[@mmalcek](https://github.com/mmalcek)修正沒有開啟任何視窗／只有系統匣時，收到WM_CLOSE訊息即關閉應用程式的問題，位於[#3990](https://github.com/wailsapp/wails/pull/3990)
- 已由@5aaee9修正garble建置，位於[#3192](https://github.com/wailsapp/wails/pull/3192)
- 已由[@leaanthony](https://github.com/leaanthony)修正Windows NSIS建置
- 已修正Linux多選對話方塊因未關閉而造成的死結
- 已修正Windows建置期間對.syso檔案的跨平台清理問題，由
- 已由@atterpac修正amd64 AppImage編譯，位於
- 已由@ansxuman修正建置資產更新，位於
- 已由@atterpac修正Linux系統匣的`OnClick`與`OnRightClick`實作
- 已修正`AlwaysOnTop`無法在Mac上運作的問題，由
-  已修正`application.NewEditMenu`包含重複項目的問題
- 🐧 已修正aarch64編譯問題
- ⊞ 已修正選項按鈕群組選單項目，由
- 修正在MacOS上建置可執行.app時，若「name」與「outputfilename」
- 已修正拖放範例使用customEventProcessor時的錯誤，由
- 🐧 已修正新增IgnoreMouseEvents所造成的Linux編譯錯誤，由
- ⊞ 已修正syso圖示檔案產生錯誤，由
- 🐧 已整合可在Wayland上原生執行的修正，來源為
- 不要繫結以下項目中的內部服務方法：
- ⊞ 已修正以下項目中的系統匣啟動panic：
- 不要繫結以下項目中的內部服務方法：
- ⊞ 已修正以下項目中的系統匣啟動panic：
- 大幅重構選單項目與事件處理。目前主要改善macOS。作者：
- 修正外掛程式與事件重構後的測試，位於
- ⊞ 已修正`Failed to unregister class Chrome_WidgetWin_0`警告。作者：
- 模組問題
- 已由[atterpac](https://github.com/atterpac)修正調整大小事件的訊息傳遞，位於
- 🐧已修正NixOS上的佈景主題處理錯誤，由
- 已修正Windows跨磁碟區安裝專案的問題，由
- 已修正React範本的CSS以顯示頁尾，由
- 已更新至最新版refresh，修正開發模式中的殭屍程序
- 已由[Atterpac](https://github.com/atterpac)修正AppImage的WebKit檔案來源問題
- 已由[Atterpac](https://github.com/Atterpac)修正Doctor的apt套件驗證，位於
- 已由@5aaee9修正應用程式結束時凍結的問題（Darwin），位於
- 已修正Windows上範例的背景色彩，由
- 已由[mmghv](https://github.com/mmghv)修正預設快顯功能表，位於
- 已修正Darwin上方向鍵的十六進位值，由
- 已讓Windows的拖放功能可正常運作。新增者：
- 已修正使用者沒有適當驅動程式時Doctor在Linux上的錯誤
- 修正啟動時的DPI縮放（Windows）。變更者：[@almas-x](https://github.com/almas-x)，位於
- 修正`go.mod`中的取代行以使用相對路徑。修正包含以下內容的Windows路徑：
- 修正MacOS系統匣未附加視窗時的點擊處理，由
- 修正未知選項造成Windows建置失敗的問題，由
- 修正Windows上沒有以下項目時，以滑鼠左鍵按一下系統匣圖示所造成的當機：
- 修正重複開啟視窗時錯誤的baseURL，由@5aaee9於PR中完成
- 修正`WebviewWindow.Restore`方法中if分支的順序，由
- 在以下情況下，跨多次`GetStartURL`呼叫正確計算`startURL`：
- 修正`Screen`結構的JS型別，使其與對應的Go型別一致，由
- 修正`WML.Reload`方法，確保已註冊事件能正確清理
- 修正Linux上自訂快顯功能表立即關閉的問題，由
- 修正繫結所產生模型檔案的輸出路徑與副檔名
- 修正繫結所產生JS程式碼中模型檔案的匯入路徑
- 修正部分Linux發行版上的拖放功能，由
- 修正使用`wails3 task dev`時缺少macOS工作項目的問題，由
- 修正註冊事件造成指派至nil map的問題，由
- 修正繫結方法參數的反序列化問題，由
- 修正繫結方法傳回多個值時的處理，由
- 修正 doctor 無法偵測未透過系統套件管理員安裝的 npm 的問題
- 修正缺少 MicrosoftEdgeWebview2Setup.exe 的問題。感謝
- 修正 Linux 上因視窗 ID 處理而隨機當機的問題，由 @leaanthony 完成。基於
- 修正 systemTray.setIcon 在 Linux 上造成當機的問題，由
- 修正首次呼叫 `setFrameless` 函式時未套用視窗框架的問題，平台為

### 變更

- **破壞性變更**：產生的 JS/TS 繫結現在會將對映鍵標示為可選，以準確反映 Go map 的語意。現在於 Typescript 中存取對映值會傳回 `T | undefined`，而非 `T`，因此必須執行 null 檢查或使用斷言（#4943），由 `@fbbdev` 完成
- 根據 `@wailsio/runtime` 的變更，將 `Event` 的用法改為 `Events`，並在 `Features/Events/Event System` 的文件中使用適當的函式呼叫，由 @AbdelhadiSeddar 完成
- 將 `EnabledFeatures`、`DisabledFeatures` 和 `AdditionalBrowserArgs` 從各視窗選項移至應用程式層級的 `Options.Windows`（#4559），由 @leaanthony 完成
- 更新 `Drag N Drop` 範例的 README，並強調此範例示範了 `Internal Drag and Drop`，由 @ndianabasi 完成
- 將多項偵錯日誌的層級從 Info 改為 Debug（由 @mbaklor 完成）
- <strong>破壞性變更：</strong>將視窗選項中的 `EnableDragAndDrop` 重新命名為 `EnableFileDrop`
- <strong>破壞性變更：</strong>將事件內容中的 `DropZoneDetails` 重新命名為 `DropTargetDetails`
- <strong>破壞性變更：</strong>將 `WindowEventContext` 上的 `DropZoneDetails()` 方法重新命名為 `DropTargetDetails()`
- <strong>破壞性變更：</strong>移除 `WindowDropZoneFilesDropped` 事件，改用 `WindowFilesDropped`
- <strong>破壞性變更：</strong>將 HTML 屬性從 `data-wails-dropzone` 改為 `data-file-drop-target`
- <strong>破壞性變更：</strong>將 CSS 懸停類別從 `wails-dropzone-hover` 改為 `file-drop-target-active`
- <strong>破壞性變更：</strong>從 Windows 移除 `DragEffect`、`OnEnterEffect`、`OnOverEffect` 選項（原為已移除之 IDropTarget 的一部分）
- 所有執行階段 JSON 處理（方法繫結、事件、webview 請求、通知、kvstore）皆改用 goccy/go-json，效能提升 21-63%，記憶體配置減少 40-60%
- 最佳化 BoundMethod 結構配置並快取 isVariadic 旗標，以降低每次呼叫的額外負擔
- 對具有 `<=8` 個引數的方法使用配置於堆疊的引數緩衝區，以避免堆積配置
- 最佳化方法呼叫中的結果收集，避免為單一傳回值配置切片
- MIME 類型快取改用 sync.Map，以改善並行效能
- 讀取 HTTP 傳輸請求本文時使用緩衝區集區
- 在內容類型探測器中延遲配置 CloseNotify 通道，以減少每個請求的配置
- 移除資產伺服器的 CSS 偵錯日誌
- 擴充 MIME 類型副檔名對映，以涵蓋 50 種以上的常見 Web 格式（字型、音訊、視訊等）
- 更新視窗 `X/Y` 選項的文件，由 @ruhuang2001 完成
- 更新 `Frontend Runtime` 文件，新增更多產生前端繫結的選項，由 @ndianabasi 完成
- 更新 Wails v3 資產伺服器的文件頁面，由 @ndianabasi 完成
- **破壞性變更**：移除套件層級的對話方塊函式（`application.InfoDialog()`、`application.QuestionDialog()` 等）。請改用 `app.Dialog` 管理員：`app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()`
- 更新對話方塊文件以符合實際 API：使用帶有回呼的 `app.Dialog.*`、`AddButton()`（而非 `SetButtons()`）、`SetDefaultButton(*Button)`（而非字串）、`AddFilter()`（而非 `SetFilters()`）、`SetFilename()`（而非 `SetDefaultFilename()`），並使用 `app.Dialog.OpenFile().CanChooseDirectories(true)` 選取資料夾
- **破壞性變更**：現在預設為正式環境組建。若要建立開發組建，請在 Taskfile 中設定 `DEV=true`。由 @leaanthony 產生新專案作為範例
- 發出包含零個或一個資料引數的自訂事件時，資料值將直接指派給 Data 欄位，不再包裝於切片中，由 [@fbbdev](https://github.com/fbbdev) 於 [#4633](https://github.com/wailsapp/wails/pull/4633) 完成
- Windows 系統匣現在會透過切換 `NIS_HIDDEN` 來遵循 `SystemTray.Show()`/`Hide()`，讓應用程式能真正消失並重新顯示（#4653）。
- 系統匣註冊會重複使用已解析的圖示、僅設定一次 `NOTIFYICON_VERSION_4`，並啟用 `NIF_SHOWTIP`，讓工具提示能在 Explorer 重新啟動後恢復（#4653）。
- macOS：視窗置中改用 `visibleFrame` 而非 `frame`，以排除選單列與 Dock 區域
- macOS：視窗置中改用 `visibleFrame` 而非 `frame`，以排除選單列與 Dock 區域
- 使用 `-config` 參數執行 `wails3 update build-assets` 時，透過 `-product*` 參數設定的值會
- `window.NativeWindowHandle()` -> `window.NativeWindow()`，由 @leaanthony 於 [#4471](https://github.com/wailsapp/wails/pull/4471) 完成
- 重構內部視窗處理，由 @leaanthony 於 [#4471](https://github.com/wailsapp/wails/pull/4471) 完成
- 移除 `application.WindowIDKey` 和 `application.WindowNameKey`（由 `application.WindowKey` 取代），由 [@leaanthony](https://github.com/leaanthony) 完成
- ContextMenuData 現在傳回字串而非 any，由 [@leaanthony](https://github.com/leaanthony) 完成
- 在 JS/TS 繫結中，固定長度陣列類型的類別欄位現在會依預期長度初始化，而不再是空陣列，由 [@fbbdev](https://github.com/fbbdev) 於 [#4001](https://github.com/wailsapp/wails/pull/4001) 完成
- ContextMenuData 現在傳回字串而非 any，由 [@leaanthony](https://github.com/leaanthony) 完成
- `application.NewService` 不再接受 options 作為可選參數（請改用 `application.NewServiceWithOptions`），由 [@leaanthony](https://github.com/leaanthony) 於 [#4024](https://github.com/wailsapp/wails/pull/4024) 完成
- 移除 `nanoid` 相依性，由 [@leaanthony](https://github.com/leaanthony) 完成
- 更新 Window 範例，以展示 mica/acrylic/tabbed 視窗樣式，由 [@leaanthony](https://github.com/leaanthony) 完成
- 在 JS/TS 繫結中，已移除 `internal.js/ts` 模型檔案；現在所有模型都可在 `models.js/ts` 中找到，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 完成
- 在 JS/TS 繫結中，具名類型絕不再呈現為其他具名類型的別名；舊行為現在僅限於別名，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 完成
- 在 JS/TS 繫結的類別模式中，類型為型別參數的結構欄位會標示為可選，且絕不會自動初始化，由 [@fbbdev](https://github.com/fbbdev) 於 [#4045](https://github.com/wailsapp/wails/pull/4045) 完成
- 從範本中移除 ESLint，由 [@IanVS](https://github.com/IanVS) 於 [#4059](https://github.com/wailsapp/wails/pull/4059) 完成
- 將著作權年份更新為 2025，由 [@IanVS](https://github.com/IanVS) 於 [#4037](https://github.com/wailsapp/wails/pull/4037) 完成
- 由 [@IanVS](https://github.com/IanVS) 在 [#4075](https://github.com/wailsapp/wails/pull/4075) 中新增 event.Sender 的文件
- 由 [@leaanthony](https://github.com/leaanthony) 新增 Go 1.24 支援
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中變更：`ServiceStartup`掛鉤現在會在呼叫`App.Run`時觸發，而非在`application.New`中觸發
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中變更：現在會從`App.Run`傳回`ServiceStartup`錯誤，而非終止程序
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中變更：來自 JS 的繫結和對話方塊呼叫現在會以錯誤物件而非字串拒絕
- 由 [@leaanthony](https://github.com/leaanthony) 改善 Windows 上的系統匣選單定位
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中將 JS 執行階段移植至 TypeScript
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中變更：執行階段會在匯入後立即初始化，無須等待視窗載入
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中變更：執行階段不再匯出 init 方法；可使用具有副作用的匯入來初始化
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中變更：繫結方法現在會傳回`CancellablePromise`；若取消，該物件會以`CancelError`拒絕，並捨棄實際呼叫結果
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中統一將內建服務型別命名為`Service`
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中統一將含選項的內建服務建立函式命名為`NewWithConfig`
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中變更：`sqlite`服務的`Select`方法現在命名為`Query`，以與 Go API 保持一致
- 範本：由 [@IanVS](https://github.com/IanVS) 在 [#4133](https://github.com/wailsapp/wails/pull/4133) 中將執行階段移至 "dependencies"，並整理 package.json 檔案
- 由 [@popaprozac](https://github.com/popaprozac) 在 [#4171](https://github.com/wailsapp/wails/pull/4171) 中變更：在開發期間建立應用程式套件並進行臨時簽署，以啟用特定 macOS API
- 由 [@leaanthony](https://github.com/leaanthony) 將建置資產移至各平台專用目錄
- 由 [@leaanthony](https://github.com/leaanthony) 將 Taskfile 移至各平台專用目錄並重新命名
- 由 [@leaanthony](https://github.com/leaanthony) 大幅改善缺少`index.html`時的使用體驗
- [Windows] 由 [@leaanthony](https://github.com/leaanthony) 改善最小化和還原的效能。以 [562589540](https://github.com/562589540) 原始的[PR](https://github.com/wailsapp/wails/pull/3955)為基礎
- 由 [@leaanthony](https://github.com/leaanthony) 移除`ShouldClose`選項（請改為替 events.Common.WindowClosing 事件註冊掛鉤）
- [Windows] 由 [@leaanthony](https://github.com/leaanthony) 減少開啟視窗時的閃爍
- 由 [@leaanthony](https://github.com/leaanthony) 移除`Window.Destroy`，因為這原本應是內部函式
- 由 [@leaanthony](https://github.com/leaanthony) 將`WindowClose`事件重新命名為`WindowClosing`
- 由 [@leaanthony](https://github.com/leaanthony) 變更：前端建置現在會依建置類型使用 vite 環境 "development" 或 "production"
- 由 [@leaanthony](https://github.com/leaanthony) 更新至 go-webview2 v1.19
- 確保使用 taskfile 的分支版本，由 @leaanthony 完成
- 更新 Taskfile 的分支版本，以修正使用以下方式安裝時的版本問題
- 使用 Taskfile 的分支版本，以修正使用以下方式安裝時的版本問題
- `service.OnStartup`現在會在發生錯誤時關閉應用程式，並執行
- 重構系統匣點擊訊息傳遞，使其更符合使用者互動方式，由
- 資產嵌入現在包含`all:frontend/dist`，以支援會產生以下內容的框架
- 由 [leaanthony](https://github.com/leaanthony) 在以下項目中重構 Taskfile
- 由以下貢獻者升級至`go-webview2` v1.0.16
- 修正`Screen`型別，使其包含`ID`而非`Id`，由
- 更新`go.mod.tmpl`的 Wails 版本，以支援`application.ServiceOptions`，由
- 由 [windom](https://github.com/windom/) 在以下項目中修正服務名稱判定
- mkdocs serve 現在會使用 Docker，由 [leaanthony](https://github.com/leaanthony) 完成
- 將開發設定整合至`config.yml`，由
- 系統匣對話方塊現在預設使用應用程式圖示（若有）（Windows），由
- 改善 macOS 的 GPU 與記憶體資訊回報，由
- 移除`WebviewGpuIsDisabled`與`EnableFraudulentWebsiteWarnings`
- 事件 API 變更：`On`/`Emit` -> 使用者事件，`OnApplicationEvent` ->
- 由 [TheGB0077](https://github.com/TheGB0077) 在以下項目中修正 Linux 上的事件 API
- [CI] 改善 Actions，並允許在分支複刻中也能執行 Actions，以及
- 由以下貢獻者將`AbsolutePosition()`重新命名為`Position()`
- 將 Linux WebKit 相依套件從 webkitgtk2-4.0 更新為 webkit2gtk-4.1，以
- 隨附的 JS 執行階段指令碼現在是 ESM 模組：匯入它的 script 標籤
- `@wailsio/runtime`套件不會在`window.wails`上公開其 API
- 視窗 API 模組`@wailsio/runtime/src/window`現在會公開其所屬的
- JS 視窗 API 已更新，以符合目前的 Go `WebviewWindow`
- 繫結產生器現在預設使用依 ID 呼叫。`-id` CLI 選項
- 新的繫結程式碼配置：輸出檔案先前依資料夾組織
- 結構欄位`application.Options.Bind`已重新命名為
- 繫結服務的新語法：服務執行個體現在必須包裝在
- 在非終端機或 CI 環境中停用旋轉指示器，由

### 移除

- **重大變更**：從各視窗的`WindowsWindow`選項中移除`EnabledFeatures`、`DisabledFeatures`和`AdditionalLaunchArgs`。請改用應用程式層級的`Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures`和`Options.Windows.AdditionalBrowserArgs`。這些旗標會全域套用至共用的 WebView2 環境（#4559），由 @leaanthony 完成
- 移除 Windows 上原生的`IDropTarget`實作，改採以 JavaScript 為基礎的方法（與v2的行為一致）
- 移除 github.com/wailsapp/mimetype 相依套件，改用擴充的副檔名對照表與標準函式庫 http.DetectContentType，將二進位檔案大小減少約1.2MB
- 實作最小化的 .desktop 檔案剖析器供 Linux 檔案總管使用，移除 gopkg.in/ini.v1 相依套件，節省約 45KB
- 改用 Go 1.21+ 標準函式庫的 slices 套件和最小化的內部輔助函式，從執行階段程式碼中移除 samber/lo，節省約 310KB
- 移除 Darwin URL 配置處理常式中的偵錯 printf 陳述式 (#4834)
- **重大變更**：移除 `linux:WindowLoadChanged` 事件；請改用 `linux:WindowLoadFinished` 偵測 WebView 何時完成載入 (#3896)，由 @leaanthony 貢獻

### 重大變更

- **Manager API 重構**：將應用程式 API 從扁平結構重新組織為分類明確的管理器，以改善程式碼組織及可探索性；由 [@leaanthony](https://github.com/leaanthony) 於 [#4359](https://github.com/wailsapp/wails/pull/4359) 中貢獻
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
- 重新命名 Service 方法：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`；由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 將 `Path` 和 `Paths` 方法移至 `application` 套件；由 [@leaanthony](https://github.com/leaanthony) 貢獻
- 應用程式選單現在僅適用於 macOS；由 [@leaanthony](https://github.com/leaanthony) 貢獻

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

## 已移除

## v3.0.0-alpha.54 - 2025-12-29

## 已新增

## 已修正

## 已移除

## v3.0.0-alpha.53 - 2025-12-27

## 已新增

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

## 已新增

## 已變更

## 已修正

## v3.0.0-alpha.47 - 2025-12-15

## 已新增

## 已修正

## v3.0.0-alpha.46 - 2025-12-14

## 已新增

## 已移除

## v3.0.0-alpha.45 - 2025-12-13

## 已新增

## 已修正

## v3.0.0-alpha.44 - 2025-12-12

## 已新增

## 已變更

## 已修正

## v3.0.0-alpha.43 - 2025-12-11

## 已新增

## v3.0.0-alpha.42 - 2025-12-10

## 已新增

## v3.0.0-alpha.41 - 2025-11-23

## 已修正

## v3.0.0-alpha.40 - 2025-11-13

## 已修正

## v3.0.0-alpha.39 - 2025-11-12

## 已新增

## 已變更

## v3.0.0-alpha.38 - 2025-11-04

## 已新增

## 已變更

## 已修正

## v3.0.0-alpha.37 - 2025-11-02

## 已修正

## v3.0.0-alpha.36 - 2025-10-15

## 已修正

## v3.0.0-alpha.35 - 2025-10-14

## 已修正

## v3.0.0-alpha.34 - 2025-10-06

## 已新增

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

不再被忽略，並會覆寫設定值。

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

`go install`，作者為@leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### 已變更

`go install`，作者為@leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### 已新增

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### 已變更

`service.OnShutdown`適用於先前已啟動的所有服務，作者為@atterpac，收錄於[#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac，收錄於[#3907](https://github.com/wailsapp/wails/pull/3907)   子資料夾，作者為@atterpac，收錄於[#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913)，收錄於[#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes)，收錄於[#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   （已由`EnabledFeatures`和`DisabledFeatures`選項取代），作者為[leaanthony](https://github.com/leaanthony)

### 已修正

通道變數，作者為@michael-freling，收錄於[#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman)，收錄於[#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   收錄於[#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony)，收錄於[#3841](https://github.com/wailsapp/wails/pull/3841)   Darwin 編輯選單中的`PasteAndMatchStyle`角色，作者為[johnmccabe](https://github.com/johnmccabe)，收錄於[#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840)，收錄於[#3854](https://github.com/wailsapp/wails/pull/3854)，作者為[kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   並不相同。作者為@nickisworking，收錄於[#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### 已新增

[mmghv](https://github.com/mmghv)，收錄於[#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac)和[leaanthony](https://github.com/leaanthony)，收錄於[#3570](https://github.com/wailsapp/wails/pull/3570)

### 已變更

應用程式事件`OnWindowEvent` -> 視窗事件，作者為[leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   以`v3/`或`v3-`為前綴的分支，作者為[stendler](https://github.com/stendler)，收錄於[#3747](https://github.com/wailsapp/wails/pull/3747)

### 已修正

[etesam913](https://github.com/etesam913)，收錄於[#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac)，收錄於[#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac)，收錄於[#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811)，收錄於[#3614](https://github.com/wailsapp/wails/pull/3614)，作者為[@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720)，作者為[leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693)，作者為[@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720)，作者為[leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693)，作者為[@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746)，作者為[@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 已修正

## v3.0.0-alpha.5 - 2024-07-30

### 已新增

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   圖示點擊，作者為@5aaee9，收錄於[#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev)，收錄於[#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)，收錄於[#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane)，收錄於[#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev)，收錄於[#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   以@Mai-Lapyst 的[PR](https://github.com/wailsapp/wails/pull/2044)為基礎   [@fbbdev](https://github.com/fbbdev)，收錄於[#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev)，收錄於[#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev)，收錄於[#3295](https://github.com/wailsapp/wails/pull/3295)   npm 套件，作者為[@fbbdev](https://github.com/fbbdev)，收錄於[#3334](https://github.com/wailsapp/wails/pull/3334)   收錄於[#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT`，作者為[@abichinger](https://github.com/abichinger)，收錄於[#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger)，收錄於[#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev)，收錄於[#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy)，收錄於[#3674](https://github.com/wailsapp/wails/pull/3674)

### 已修正

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) 於   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) 於   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) 於   [#3477](https://github.com/wailsapp/wails/pull/3477)   由 [Atterpac](https://github.com/atterpac) 於   [#3320](https://github.com/wailsapp/wails/pull/3320)。   於 [#3306](https://github.com/wailsapp/wails/pull/3306)。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) 於   [#2750](https://github.com/wailsapp/wails/pull/2750)。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   [jaybeecave](https://github.com/jaybeecave) 於   [#3052](https://github.com/wailsapp/wails/pull/3052)。   [@pylotlight](https://github.com/pylotlight) 於   [PR](https://github.com/wailsapp/wails/pull/3039)   安裝。由 [@pylotlight](https://github.com/pylotlight) 於   [PR](https://github.com/wailsapp/wails/pull/3032) 新增   [PR](https://github.com/wailsapp/wails/pull/3145)   空格 — @leaanthony。   [thomas-senechal](https://github.com/thomas-senechal) 於 PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) 於 PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   附加視窗 [tw1nk](https://github.com/tw1nk) 於 PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) 於   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` 存在。   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   監聽器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) 於   [#3330](https://github.com/wailsapp/wails/pull/3330)   產生器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3334](https://github.com/wailsapp/wails/pull/3334)   產生器，由 [@fbbdev](https://github.com/fbbdev) 於   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) 於   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) 於   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) 於   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) 於   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) 於   [#3431](https://github.com/wailsapp/wails/pull/3431)   由 [@pekim](https://github.com/pekim) 於   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel)。   PR [#3466](https://github.com/wailsapp/wails/pull/3622)，由   [@5aaee9](https://github.com/5aaee9) 提交。   [@windom](https://github.com/windom/) 於   [#3636](https://github.com/wailsapp/wails/pull/3636)。   Windows，由 [@bruxaodev](https://github.com/bruxaodev/) 於   [#3691](https://github.com/wailsapp/wails/pull/3691)。

### 已變更

[mmghv](https://github.com/mmghv) 於   [#3611](https://github.com/wailsapp/wails/pull/3611)   支援 Ubuntu 24.04 LTS，由 [atterpac](https://github.com/atterpac) 於   [#3461](https://github.com/wailsapp/wails/pull/3461)   必須具有 `type="module"` 屬性。由   [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   物件，且不會啟動 WML 系統。這項變更旨在改善封裝。若需要，可呼叫新的 `WML.Enable` 方法手動啟動 WML 系統。隨附的 JS 執行階段指令碼仍會自動執行這兩項操作。由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   視窗物件作為預設匯出。現在已無法透過 ESM 具名匯入或命名空間匯入語法匯入個別方法。   API。部分方法的名稱或原型已變更，具體如下：`Screen`   變為 `GetScreen`；`GetZoomLevel`/`SetZoomLevel` 變為 `GetZoom`/`SetZoom`；   `GetZoom`、`Width` 和 `Height` 現在會直接傳回值，而不是將值包裝在物件中。由 [@fbbdev](https://github.com/fbbdev) 於   [#3295](https://github.com/wailsapp/wails/pull/3295)   已移除。使用 `-names` CLI 選項可切換回依名稱呼叫。   由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   先前以其所在套件命名；現在則使用完整的 Go 匯入路徑，包括模組路徑。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   對 `application.NewService` 的呼叫。由 [@fbbdev](https://github.com/fbbdev) 於   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) 於   [#3574](https://github.com/wailsapp/wails/pull/3574)
